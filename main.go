package main

import (
	"encoding/json"
	"flag"
	. "github.com/DanielCalvo/markdownscanner/markdownlink"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
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

//Possibly split into: CloneGitRepository, UpdateGitRepository
//Only works with http(s). Will there be SSH support?
//There is indeed something wrong with this function as oftentimes the program crashes when trying to clone a repository for the first time
func GetGitRepository(repositoryUrl, tmpDir string) (string, error) {
	//tmpDir = CheckAndAddPathSeparatorSuffix(tmpDir)

	url, err := url.ParseRequestURI(repositoryUrl)
	if err != nil {
		return "", err
	}

	gitRepoFolder := tmpDir + url.Path

	_, fsErr := os.Stat(gitRepoFolder)

	//There's a bug in here? Are you pulling if the repo already exists?
	if os.IsNotExist(fsErr) {
		_, err := git.PlainClone(gitRepoFolder, false, &git.CloneOptions{
			URL:      repositoryUrl,
			Progress: os.Stdout,
		})
		if err != nil {
			return "", err
		}
	} else if fsErr != nil {
		return gitRepoFolder, fsErr
	}

	repository, err := git.PlainOpen(gitRepoFolder)
	if err != nil {
		return "", err
	}

	workTree, err := repository.Worktree()
	if err != nil {
		return "", err
	}

	err = workTree.Pull(&git.PullOptions{RemoteName: "origin"})
	if err == git.NoErrAlreadyUpToDate {
		return gitRepoFolder, nil
	} else {
		return "", err
	}
}

//You need to handle the errors here
//"info" is a poor variable name
func GetMarkdownFiles(repositoryFilesystemPath, gitRepositoryUrl string) []MarkdownFile {
	var markdownFiles []MarkdownFile

	url, err := url.ParseRequestURI(gitRepositoryUrl)
	if err != nil {
		return nil
	}

	//rename this "info" in here maybe, that's not a very descriptive variable name
	err = filepath.Walk(repositoryFilesystemPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".md") {
			s := strings.Split(path, url.Path)
			mdFile := MarkdownFile{
				FilePath: path,
				FileName: info.Name(),
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

//leave some big comments explaining the regexes in this function
//this can still be improved, the regex could go on their own functions and run in parallel
func GetMarkdownLinksFromFiles(mdFiles []MarkdownFile) ([]MarkdownLink, error) {

	var markdownLinks []MarkdownLink

	for _, mdFile := range mdFiles {
		fileContents, err := ioutil.ReadFile(mdFile.FilePath)
		if err != nil {
			return nil, err
		}

		//GetInlineLinks
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

		//GetFootnoteLinks
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

func CheckMarkdownLinksWorker(mdLinkIn <-chan MarkdownLink, workerNum int) <-chan MarkdownLink {
	mdLinkOut := make(chan MarkdownLink)
	var wg sync.WaitGroup

	wg.Add(workerNum)
	go func() {
		for i := 0; i < workerNum; i++ {
			go func() {
				defer wg.Done()
				for mdLink := range mdLinkIn {
					mdLink.CheckLink()
					mdLinkOut <- mdLink
				}
			}()
		}
		wg.Wait()
		close(mdLinkOut)
	}()
	return mdLinkOut
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

//Name things properly in this function
//Leave a comment or two in here

func CheckMarkdownLinks(mdLinks []MarkdownLink, workerNum int) []MarkdownLink {
	linkChan := make(chan MarkdownLink)
	go func() {
		for _, link := range mdLinks {
			linkChan <- link
		}
		close(linkChan)
	}()

	checkedLinks := CheckMarkdownLinksWorker(linkChan, workerNum)

	var mdLinksProcessed []MarkdownLink

	for checkedLink := range checkedLinks {
		mdLinksProcessed = append(mdLinksProcessed, checkedLink)
		log.Println("Checked:", checkedLink)
	}
	return mdLinksProcessed
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

func GetRepoFilenameWithExtension(repoUrl, extension string) (string, error) {
	url, err := url.ParseRequestURI(repoUrl)
	if err != nil {
		return "", err
	}
	urlPathUnderscores := strings.ReplaceAll(url.Path, "/", "_")
	return url.Host + urlPathUnderscores + "." + extension, nil
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

	//Uh, a hardcoded path. This should be changed...
	//SaveStructToHtml in here perhaps?
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

func GetRepositoriesFromYaml(yamlPath string) ([]string, error) {
	var gitRepositories []string

	repositoriesYamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(repositoriesYamlFile, &gitRepositories)
	if err != nil {
		log.Println("Cannot unmarshall repositories.yaml:", err)
		panic(err)
	}
	return gitRepositories, nil
}

func GenerateIndexHtml(gc GlobalConfig) error {
	var scans []ScanMetadata

	err := filepath.Walk(gc.StaticFolder, func(path string, info os.FileInfo, err error) error {
		var scan ScanMetadata

		//This can be improved, what you want is "starts with"
		//"Info" is not a good variable name
		if strings.Contains(info.Name(), "metadata_") {
			jsonFile, err := os.Open(path)
			if err != nil {
				log.Println(err)
			}
			defer jsonFile.Close()
			//Ignoring error? Questionable :o
			scanMetadataJson, _ := ioutil.ReadAll(jsonFile)
			err = json.Unmarshal(scanMetadataJson, &scan)
			if err != nil {
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
	slowScan := flag.Bool("slowscan", false, "Flag to scan things slowly to avoid generating too many requests")
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

	gitRepositoryUrls, err := GetRepositoriesFromYaml(globalConfig.RepositoriesYamlFile)
	if err != nil {
		log.Println("Error getting repositories from yaml", err)
		panic(err)
	}

	log.Println("Beginning main loop: Iterating over Git Repositories from repositories.yaml")
	for _, gitRepositoryUrl := range gitRepositoryUrls {

		gitRepoFilesystemPath, err := GetGitRepository(gitRepositoryUrl, globalConfig.RepositoriesFolder)
		if err != nil {
			log.Println("Error cloning or updating: "+gitRepositoryUrl, err)
			continue
		}

		markdownFiles := GetMarkdownFiles(gitRepoFilesystemPath, gitRepositoryUrl)
		log.Println("Found", len(markdownFiles), "markdown files on", gitRepositoryUrl)

		markdownLinks, err := GetMarkdownLinksFromFiles(markdownFiles)
		if err != nil {
			log.Println("Unable to GetMarkdownLinksFromFiles:", err)
			continue
		}
		log.Println("Found", len(markdownLinks), "markdown links on", gitRepositoryUrl)

		var checkedLinks []MarkdownLink

		if *slowScan {
			checkedLinks = CheckMarkdownLinksWithSleep(markdownLinks, time.Second)
		} else {
			checkedLinks = CheckMarkdownLinks(markdownLinks, globalConfig.WorkerNum)
		}
		log.Println("Markdown link check complete")

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
