package main

import (
	"errors"
	"fmt"
	"github.com/DanielCalvo/markdownscanner/pkg/mdscanner"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

//First things first: Create a function to convert a github URL to a github RAW url
/*
There are 2 ways of doing this:
- https://raw.githubusercontent.com/kubernetes-sigs/external-dns/master/docs/tutorials/alb-ingress.md
- https://api.github.com/repos/kubernetes-sigs/external-dns/contents/docs/tutorials/alb-ingress.md (with curl -H 'Accept: application/vnd.github.v3.raw')
- Both require modifying the URL. I'll go with the raw one.
*/

//Pls write a test for this function
func GetGithubRawUrlFromGithubUrl(s string) (string, error) {
	ghUrl, err := url.Parse(s)

	if err != nil {
		return "", err
	}

	ghUrl.Host = "raw.githubusercontent.com"
	urlSplit := strings.Split(ghUrl.Path, "/")

	//You will have a remarkably poor time trying to debug this. Don't try to be this clever, make it simpler
	finalString := strings.Join([]string{urlSplit[1], urlSplit[2], "master", strings.Join(urlSplit[5:], "/")}, "/")

	ghUrl.Path = finalString
	return ghUrl.String(), nil
}

//You copied and pasted that from here: https://golangbyexample.com/download-image-fileurl-golang/
func DownloadFile(URL string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}

	file, err := os.Create(path.Base(URL))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Println("hey")

	if os.Getenv("fileurl") == "" {
		log.Fatal("HTTP file target not passed as \"fileurl\" environment variable")
	}

	//Download file to the current dir (no tmpdir shenanigans!)

	fmt.Println(os.Getwd())

	raw, _ := GetGithubRawUrlFromGithubUrl(os.Getenv("fileurl"))
	fmt.Println(raw)

	err := DownloadFile(raw)
	if err != nil {
		fmt.Println(err)
	}

	//There's a proper way to do this...
	mdFile := mdscanner.MarkdownFile{
		FileName: path.Base(raw),
		FilePath: path.Base(raw),
		HTTPAddr: raw,
	}

	mdLinks, err := mdscanner.GetMarkdownLinksFromFile(mdFile)
	if err != nil {
		fmt.Println(err)
	}

	mdLinks = mdscanner.CheckMarkdownLinksWithSleep(mdLinks, time.Second)
	mdLinks = mdscanner.SortLinksBy404(mdLinks)
	//RemoveFile?

	err = mdscanner.SaveStructToJson(mdLinks, "single-file.json")
	if err != nil {
		fmt.Println(err)
	}

}
