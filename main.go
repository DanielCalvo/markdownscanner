package main

import (
	"encoding/json"
	"flag"
	. "github.com/DanielCalvo/markdownscanner/markdownscanner"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type MarkdownFile struct {
	FileName string
	FilePath string
	HTTPAddr string
}

type ScanMetadata struct {
	ProjectName            string
	RepoURL                string
	LastScanned            string
	LinksScanned           int
	Links404               int
	HTMLReportFilename     string
	JSONReportFilename     string
	MetadataReportFileName string
}

type GlobalConfig struct {
	RepositoriesFolder   string
	StaticFolder         string
	TemplateFolder       string
	RepositoriesYamlFile string
	WorkerNum            int
}

//If ran on a Unix system:
//If this function receives /tmp it will return /tmp/
//If it receives /tmp/ it will return /tmp/
func CheckAndAddPathSeparatorSuffix(fsPath string) string {
	if !strings.HasSuffix(fsPath, string(os.PathSeparator)) {
		fsPath = fsPath + string(os.PathSeparator)
		return fsPath
	} else {
		return fsPath
	}
}

//You need to handle the errors here
func GetMarkdownFiles(repositoryFilesystemPath, gitRepositoryUrl string) []MarkdownFile {
	var markdownFiles []MarkdownFile

	//Returning nil but no err? Is this ok?
	url, err := url.ParseRequestURI(gitRepositoryUrl)
	if err != nil {
		return nil
	}

	err = filepath.Walk(repositoryFilesystemPath, func(path string, file os.FileInfo, err error) error {
		if strings.HasSuffix(file.Name(), ".md") && DoesExist(file.Name()) {
			s := strings.Split(path, url.Path)
			mdFile := MarkdownFile{
				FilePath: path,
				FileName: file.Name(),
				HTTPAddr: gitRepositoryUrl + "/tree/master" + s[1],
			}
			markdownFiles = append(markdownFiles, mdFile)
		}
		return err
	})

	if err != nil {
		log.Printf("Error walking the path %q: %v\n", err)
	}
	return markdownFiles
}

func GetMarkdownLinksFromFiles(mdFiles []MarkdownFile) ([]MarkdownLink, error) {

	var markdownLinks []MarkdownLink

	for _, mdFile := range mdFiles {
		fileContents, err := ioutil.ReadFile(mdFile.FilePath)
		if err != nil {
			return nil, err
		}

		//regex for footnote style MarkdownLinks
		re := regexp.MustCompile(`(\[.+\])\s*:\s*(.+)`)
		for _, matchedMarkdownLink := range re.FindAllStringSubmatch(string(fileContents), -1) {
			mdLink := MarkdownLink{
				FileName:      mdFile.FileName,
				LocalFilePath: mdFile.FilePath,
				HTTPFilePath:  mdFile.HTTPAddr,
				Name:          matchedMarkdownLink[1],
				Destination:   matchedMarkdownLink[2],
			}
			markdownLinks = append(markdownLinks, mdLink)
		}

		//regex for inline style links
		re = regexp.MustCompile(`(\[.+?\])((\()(.+?)(\)))`)
		for _, matchedMarkdownLink := range re.FindAllStringSubmatch(string(fileContents), -1) {
			mdLink := MarkdownLink{
				FileName:      mdFile.FileName,
				LocalFilePath: mdFile.FilePath,
				HTTPFilePath:  mdFile.HTTPAddr,
				Name:          matchedMarkdownLink[1],
				Destination:   matchedMarkdownLink[4],
			}
			markdownLinks = append(markdownLinks, mdLink)
		}
	}
	return markdownLinks, nil
}

func CheckMarkdownLinksWithSleep(mdLinks []MarkdownLink, sleepTime time.Duration) []MarkdownLink {
	var scannedLinks []MarkdownLink
	for _, mdLink := range mdLinks {
		mdLink.CheckLink()
		if mdLink.Type == "HTTP" {
			time.Sleep(sleepTime)
		}
		scannedLinks = append(scannedLinks, mdLink)
		log.Println("Checked:", mdLink)
	}
	return scannedLinks
}

func Count404MarkdownLinks(mdLinks []MarkdownLink) int {
	var c int
	for _, link := range mdLinks {
		if link.Status == "404" {
			c++
		}
	}
	return c
}

func SortLinksByStatus(mdLinks []MarkdownLink, status string) []MarkdownLink {
	var tmpLinks []MarkdownLink

	//returns a slice bounds out of error if 404 link is on the last element. Redo do this!
	for i, v := range mdLinks {
		if strings.HasPrefix(v.Status, "4") {
			tmpLinks = append(tmpLinks, v)
			mdLinks = append(mdLinks[:i], mdLinks[i+1:]...)
		}
	}
	mdLinks = append(tmpLinks, mdLinks...)
	return mdLinks
}

func SaveStructToJson(s interface{}, jsonFilesystemPath string) error {
	resultJson, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(jsonFilesystemPath, resultJson, 0644)
	if err != nil {
		return err
	}
	return nil
}

func SaveCheckedLinksToJsonAndHtml(mdLinks []MarkdownLink, gitRepositoryUrl string, gc GlobalConfig) error {
	repoJsonFileName, err := GetRepoFilenameWithExtension(gitRepositoryUrl, "json")
	if err != nil {
		return err
	}

	err = SaveStructToJson(mdLinks, gc.StaticFolder+repoJsonFileName)
	if err != nil {
		return err
	}

	HtmlFileName, err := GetRepoFilenameWithExtension(gitRepositoryUrl, "html")
	if err != nil {
		return err
	}

	tpl := template.Must(template.ParseFiles(gc.TemplateFolder + string(os.PathSeparator) + "results_table.gohtml"))

	htmlFile, err := os.Create(gc.StaticFolder + string(os.PathSeparator) + HtmlFileName)
	if err != nil {
		return err
	}
	defer htmlFile.Close()

	err = tpl.Execute(htmlFile, mdLinks)
	if err != nil {
		return err
	}
	return nil
}

func GenerateIndexHtml(gc GlobalConfig) error {
	var scans []ScanMetadata

	err := filepath.Walk(gc.StaticFolder, func(path string, file os.FileInfo, err error) error {
		var scan ScanMetadata

		//What you actually want is "starts with". Improve this!
		if strings.Contains(file.Name(), "metadata_") {
			jsonFile, err := os.Open(path)
			if err != nil {
				log.Println(err)
			}
			defer jsonFile.Close()
			//Ignoring error? Questionable :o
			scanMetadataJson, _ := ioutil.ReadAll(jsonFile)
			err = json.Unmarshal(scanMetadataJson, &scan)
			if err != nil {
				//panicking doesn't seem very productive
				panic(err)
			}
			scans = append(scans, scan)
		}
		return nil
	})
	if err != nil {
		return err
	}

	indexTpl := template.Must(template.ParseFiles(gc.TemplateFolder + "/index.gohtml"))

	indexFile, err := os.Create(gc.StaticFolder + string(os.PathSeparator) + "index.html")
	if err != nil {
		return err
	}
	defer indexFile.Close()

	err = indexTpl.Execute(indexFile, scans)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	log.SetOutput(os.Stdout)

	webMode := flag.Bool("webmode", false, "Does this application start in webserver mode?")
	flag.Parse()

	pwd, err := os.Getwd()
	if err != nil {
		log.Println("Could not get current directory at startup")
		panic(err)
	}

	globalConfig := GlobalConfig{
		RepositoriesFolder:   pwd + string(os.PathSeparator) + "repositories" + string(os.PathSeparator),
		StaticFolder:         pwd + string(os.PathSeparator) + "html_static" + string(os.PathSeparator),
		TemplateFolder:       pwd + string(os.PathSeparator) + "html_templates" + string(os.PathSeparator),
		RepositoriesYamlFile: pwd + string(os.PathSeparator) + "repositories.yaml",
		WorkerNum:            6,
	}

	log.Println("Global Config set:", globalConfig)

	_, err = os.Stat(globalConfig.StaticFolder)
	if os.IsNotExist(err) {
		err = os.Mkdir(globalConfig.StaticFolder, 0755)
		if err != nil {
			log.Println("Unable to create results directory: "+globalConfig.StaticFolder+", error:", err)
			panic(err)
		}
	}

	if *webMode {
		log.Println("webmode!")
		http.Handle("/", http.FileServer(http.Dir(globalConfig.StaticFolder)))
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

	gitRepositoryUrls, err := GetRepositoryUrlsFromYaml(globalConfig.RepositoriesYamlFile)
	if err != nil {
		log.Println("Error getting repositories from yaml", err)
		panic(err)
	}

	gitRepositoryUrls = SortRepositoriesByUnscannedFirst(gitRepositoryUrls, globalConfig.StaticFolder)

	log.Println("Beginning main loop: Iterating over Git Repositories from repositories.yaml")

	for _, gitRepositoryUrl := range gitRepositoryUrls {
		log.Println("Clonning", gitRepositoryUrl)
		gitRepoFilesystemPath, err := CloneGitRepository(gitRepositoryUrl, globalConfig.RepositoriesFolder)
		if err != nil {
			log.Println("Error cloning or updating: "+gitRepositoryUrl, err)
			continue
		}
		log.Println("Cloning complete for:", gitRepositoryUrl)

		log.Println("Getting markdownfiles for", gitRepositoryUrl)
		markdownFiles := GetMarkdownFiles(gitRepoFilesystemPath, gitRepositoryUrl)
		log.Println("Found", len(markdownFiles), "markdown FILES on", gitRepositoryUrl)

		markdownLinks, err := GetMarkdownLinksFromFiles(markdownFiles)
		if err != nil {
			log.Println("Unable to GetMarkdownLinksFromFiles:", err)
			continue
		}
		log.Println("Found", len(markdownLinks), "markdown LINKS on", gitRepositoryUrl)

		var checkedLinks []MarkdownLink
		checkedLinks = CheckMarkdownLinksWithSleep(markdownLinks, time.Second)

		log.Println("Markdown link check complete")

		checkedLinks = SortLinksByStatus(checkedLinks, "404")

		err = SaveCheckedLinksToJsonAndHtml(checkedLinks, gitRepositoryUrl, globalConfig)
		if err != nil {
			log.Println("Could not SaveCheckedLinksToJsonAndHtml:", err)
		}
		log.Println("Links for", gitRepositoryUrl, "saved to json")

		//See what to do about these possible errors at some point
		projectUrl, _ := url.ParseRequestURI(gitRepositoryUrl)
		repoJsonFileName, _ := GetRepoFilenameWithExtension(gitRepositoryUrl, "json")
		repoHtmlFileName, _ := GetRepoFilenameWithExtension(gitRepositoryUrl, "html")

		sm := ScanMetadata{
			ProjectName:            projectUrl.Path,
			RepoURL:                gitRepositoryUrl,
			LastScanned:            time.Now().Format(time.RFC3339),
			LinksScanned:           len(checkedLinks),
			Links404:               Count404MarkdownLinks(checkedLinks),
			HTMLReportFilename:     repoHtmlFileName,
			JSONReportFilename:     repoJsonFileName,
			MetadataReportFileName: globalConfig.StaticFolder + string(os.PathSeparator) + "metadata_" + repoJsonFileName,
		}

		err = SaveStructToJson(sm, sm.MetadataReportFileName)
		if err != nil {
			log.Println("Could not SaveStructToJson:", err)
		}
		log.Println("Scan metadata saved to json")

		err = GenerateIndexHtml(globalConfig)
		if err != nil {
			log.Println("Could not GenerateIndexHtml:", err)
		}
		log.Println("Index page updated")

		log.Println("Deleting:", gitRepoFilesystemPath)
		err = os.RemoveAll(gitRepoFilesystemPath)
		if err != nil {
			log.Println("Could not remove gitRepoFilesystemPath:", err)
		}
	}
}
