package main

import (
	"flag"
	"fmt"
	"github.com/DanielCalvo/markdownscanner/pkg/config"
	"github.com/DanielCalvo/markdownscanner/pkg/mdscanner"
	"log"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)

	configFile := flag.String("config.file", "", "Filesystem path for markdown scanner configuration file")
	flag.Parse()

	log.Println("Loading config file")
	conf, err := config.LoadFile(*configFile)
	if err != nil {
		fmt.Println("Could not load config file:", err)
		os.Exit(1)
	}

	log.Println("Initializing config file")
	err = config.Initialize(conf)
	if err != nil {
		fmt.Println("Error initiating config:", err)
		os.Exit(1)
	}

	//Please improve this and put it in a function, it's awful here
	//Also merge conf.Repositories and the repositories from conf.GithubProjects properly into another slice
	var projectRepos []string
	for _, p := range conf.GithubProjects {
		prepos, err := mdscanner.GetRepoUrlsFromProject(p)
		if err != nil {
			continue
		}
		projectRepos = append(projectRepos, prepos...)
	}

	//conf.Repositories is a poor name. conf.RepoUrls would be better I think
	conf.Repositories = append(conf.Repositories, projectRepos...)

	//Later at some point you can read from disk when repos were last scanned too and sort by oldest scan first
	repositories := mdscanner.NewRepositories(conf, conf.Repositories)
	repositories = mdscanner.SortRepositoriesByUnscannedFirst(repositories)

	for _, repo := range repositories {
		fmt.Println(repo.Name)

		log.Println("Cloning repostory:", repo.Name)
		err = mdscanner.CloneRepository(repo)
		if err != nil {
			fmt.Println("Error clonning repository:", err)
			continue
		}

		log.Println("Getting Markdown files")
		mdscanner.GetMarkdownFiles(&repo)
		log.Println("Getting Markdown links from files")
		mdscanner.GetMarkdownLinksFromFiles(&repo)

		log.Println("Deleting repository", repo.Name)
		err = mdscanner.DeleteRepository(repo)
		if err != nil {
			fmt.Println("Unable to delete repository")
			continue
		}

		//Figure out later why this doesn't work in place
		log.Println("Checking Markdown links")
		repo.MarkdownLinks = mdscanner.CheckMarkdownLinksWithSleep(&repo, time.Second)
		repo.MarkdownLinks = mdscanner.SortLinksBy404(repo.MarkdownLinks)

		log.Println("Uploading scan result to S3")
		err = mdscanner.UploadResultsToS3(*conf, repo)

		if err != nil {
			fmt.Println("Could not upload to S3:", err)
			continue
		}

		//put this into a function to make it prettier!
		repo.LastScanned = time.Now().Format(time.RFC3339)
		repo.LinksScanned = len(repo.MarkdownLinks)
		repo.Links404 = mdscanner.Count404MarkdownLinks(repo.MarkdownLinks)

		log.Println("Saving scan metadata")
		err = mdscanner.SaveScanMetadata(repo)
		if err != nil {
			log.Println("Could not SaveScanMetadata:", err)
		}

		log.Println("Generating and uploading index html")
		err = mdscanner.GenerateAndUploadIndexHtml(*conf)
		if err != nil {
			log.Println("Could not GenerateIndexHtml:", err)
		}
	}
}
