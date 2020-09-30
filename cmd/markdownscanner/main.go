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

	conf, err := config.LoadFile(*configFile)
	if err != nil {
		fmt.Println("Could not load config file:", err)
		os.Exit(1)
	}

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
	conf.Repositories = append(conf.Repositories, projectRepos...)

	//Sort repos by Unscanned first

	for _, repoURL := range conf.Repositories {
		repo, err := mdscanner.NewRepository(conf, repoURL)
		if err != nil {
			fmt.Println("Error encountered initializing repository:", err)
			continue
		}

		err = mdscanner.CloneRepository(repo)
		if err != nil {
			fmt.Println("Error clonning repository:", err)
			continue
		}

		mdscanner.GetMarkdownFiles(&repo)
		mdscanner.GetMarkdownLinksFromFiles(&repo)

		err = mdscanner.DeleteRepository(repo)
		if err != nil {
			fmt.Println("Unable to delete repository")
			continue
		}

		//Figure out later why this doesn't work in place
		repo.MarkdownLinks = mdscanner.CheckMarkdownLinksWithSleep(&repo, time.Second)
		//checkedLinks = SortLinksByStatus(checkedLinks, "404")

		err = mdscanner.UploadResultsToS3(*conf, repo)

		if err != nil {
			fmt.Println("Could not upload to S3:", err)
			continue
		}

		//put this into a function to make it prettier!
		repo.LastScanned = time.Now().Format(time.RFC3339)
		repo.LinksScanned = len(repo.MarkdownLinks)
		repo.Links404 = mdscanner.Count404MarkdownLinks(repo.MarkdownLinks)

		err = mdscanner.SaveScanMetadata(repo)
		if err != nil {
			log.Println("Could not SaveScanMetadata:", err)
		}
		log.Println("Scan metadata saved to json")

		err = mdscanner.GenerateAndUploadIndexHtml(*conf)
		if err != nil {
			log.Println("Could not GenerateIndexHtml:", err)
		}
	}
}
