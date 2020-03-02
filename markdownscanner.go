package main

import (
	"fmt"
	. "github.com/DanielCalvo/markdownscanner/markdownlink"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

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

func PrintAndPanic(s string, err error) {
	fmt.Print(s)
	panic(err)
}

//CloneGitRepository
//UpdateGitRepository
//Only works with http(s). Will there be SSH support?
func GetGitRepository(repositoryUrl, tmpDir string) (string, error) {
	tmpDir = CheckAndAddPathSeparatorSuffix(tmpDir)
	url, err := url.ParseRequestURI("http://github.com/kubernetes/kubectl")
	if err != nil {
		return "", err
	}

	repoFilesystemPath := tmpDir + url.Path
	_, fsErr := os.Stat(repoFilesystemPath)

	if os.IsNotExist(fsErr) {
		log.Println("Cloning", repositoryUrl)
		_, err := git.PlainClone(repoFilesystemPath, false, &git.CloneOptions{
			URL:      repositoryUrl,
			Progress: os.Stdout,
		})
		if err != nil {
			return "", err
		}
	} else if fsErr != nil {
		return repoFilesystemPath, fsErr
	}

	repository, err := git.PlainOpen(repoFilesystemPath)
	if err != nil {
		return "", err
	}

	workTree, err := repository.Worktree()
	if err != nil {
		return "", err
	}

	log.Println("Pulling", repositoryUrl)
	err = workTree.Pull(&git.PullOptions{RemoteName: "origin"})
	if err == git.NoErrAlreadyUpToDate {
		return repoFilesystemPath, nil
	} else {
		return "", err
	}
}

func GetMarkdownFilepaths(repoFilesystemPath string) []string {
	var MarkdownFilepaths []string

	err := filepath.Walk(repoFilesystemPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".md") {
			MarkdownFilepaths = append(MarkdownFilepaths, path)
		}
		return err
	})

	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", err)
	}
	return MarkdownFilepaths
}

////leave some big comments explaining the regexes in this function
//this can still be improved, the regex could go on their own functions and run in parallel
func GetMarkdownLinksFromFiles(filePaths []string) ([]MarkdownLink, error) {

	var markdownLinks []MarkdownLink

	for _, f := range filePaths {
		fileContents, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}

		//GetInlineLinks
		//regex for footnote style MarkdownLinks
		re := regexp.MustCompile(`(\[.+\])\s*:\s*(.+)`)
		for _, matchedMarkdownLink := range re.FindAllStringSubmatch(string(fileContents), -1) {
			mdLink := MarkdownLink{
				File:        f,
				Name:        matchedMarkdownLink[1],
				Destination: matchedMarkdownLink[2],
			}
			markdownLinks = append(markdownLinks, mdLink)
		}

		//GetFootnoteLinks
		//regex for inline style links
		re = regexp.MustCompile(`(\[.+?\])((\()(.+?)(\)))`)
		for _, matchedMarkdownLink := range re.FindAllStringSubmatch(string(fileContents), -1) {
			mdLink := MarkdownLink{
				File:        f,
				Name:        matchedMarkdownLink[1],
				Destination: matchedMarkdownLink[4],
			}
			markdownLinks = append(markdownLinks, mdLink)
		}
	}

	return markdownLinks, nil
}

//
//func handler(w http.ResponseWriter, req *http.Request) {
//	///tmp/mdscanner/results/kubernetes/kubelet/
//
//	jsonFile, err := os.Open(config.tmpDir + "/results" + req.URL.Path + "/report.json")
//	if err != nil {
//		fmt.Fprintf(w, "Could not open json at: "+config.tmpDir+"/results"+req.URL.Path+"/report.json")
//	}
//	byteValue, _ := ioutil.ReadAll(jsonFile)
//
//	fmt.Fprintf(w, string(byteValue))
//
//}

//Implement logging!

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

//Name things properly in this function
//Leave a comment or two in here
//"myLinks" and "n" are not proper variable names
func CheckMarkdownLinks(mdLinks []MarkdownLink) []MarkdownLink {
	linkChan := make(chan MarkdownLink)
	go func() {
		for _, link := range mdLinks {
			linkChan <- link
		}
		close(linkChan)
	}()

	myLinks := CheckMarkdownLinksWorker(linkChan, 30)

	var mdLinksProcessed []MarkdownLink

	for n := range myLinks {
		mdLinksProcessed = append(mdLinksProcessed, n)
	}
	return mdLinksProcessed
}

//You need to add a way to ignore certain MarkdownLinks (minutes, slack MarkdownLinks, mailto, etc)
//Formatting the JSON to one check per line would be neat for readability
//What arguments are you going to be taking? Call from web part of service? Command line flags?
//Unit tests are missing
//How do you handle an invalid repository?
//Think about your queueing package, do a queue.Start() at the beginning and have other things send things to it maybe!
//Don't forge to find a way to implement a list of things you want to ignore: if strings.Contains(strings.ToLower(l.File), "changelog") || strings.Contains(strings.ToLower(l.File), "minute") || strings.Contains(strings.ToLower(l.File), "meeting") || strings.Contains(strings.ToLower(l.File), "release") {
//Make your impromptu comments actual package documentation (that would be cool!)
func main() {

	log.SetOutput(os.Stdout)
	GitRepository := "https://github.com/kubernetes/kubectl"
	tmpDir := "/tmp"

	gitRepoFilesystemPath, err := GetGitRepository(GitRepository, tmpDir)
	if err != nil {
		PrintAndPanic("Error running GetGitRepository:", err)
	}

	markdownFilepaths := GetMarkdownFilepaths(gitRepoFilesystemPath)

	var markdownLinks []MarkdownLink

	markdownLinks, err = GetMarkdownLinksFromFiles(markdownFilepaths)
	if err != nil {
		fmt.Println(err)
	}

	checkedLinks := CheckMarkdownLinks(markdownLinks)

	for _, v := range checkedLinks {
		fmt.Println(v)
	}

	//Do something with the results!

}

//	jsonSavePath := config.tmpDir + "/results" + u.EscapedPath()
//
//	err = os.MkdirAll(jsonSavePath, os.ModePerm)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	file, _ := json.MarshalIndent(ll, "", "")
//	_ = ioutil.WriteFile(jsonSavePath+"/report.json", file, 0644)
//	log.Println("Report saved at: ", jsonSavePath+"/report.json")
//
//}
//
//http.HandleFunc("/", handler)
//log.Fatal(http.ListenAndServe(":8080", nil))
