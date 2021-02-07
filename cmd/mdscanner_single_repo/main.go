package main

import (
	"fmt"
	"github.com/DanielCalvo/markdownscanner/pkg/config"
	"github.com/DanielCalvo/markdownscanner/pkg/mdscanner"
	"log"
	"os"
	"time"
)

func main() {
	if os.Getenv("repo") == "" {
		log.Fatal("Repository not passed as env variable \"repo\" environment variable")
	}

	c := config.Config{} //Set sane values for config if you find no file? (aka running from container in cmdline mode?)
	repo, err := mdscanner.NewRepository(&c, os.Getenv("repo"))

	if err != nil {
		log.Fatal("Error instantiating new repository:", err)
	}

	fmt.Print(repo)

	err = mdscanner.CloneRepository(repo)
	if err != nil {
		log.Fatal("Error clonning repository:", err)
	}

	log.Println("Getting Markdown files")
	mdscanner.GetMarkdownFiles(&repo)
	log.Println("Getting Markdown links from files")
	mdscanner.GetMarkdownLinksFromFiles(&repo)

	log.Println("Deleting repository", repo.Name)
	err = mdscanner.DeleteRepository(repo)
	if err != nil {
		log.Println("Unable to delete repository:", err)
	}

	//Figure out later why this doesn't work in place
	log.Println("Checking Markdown links")
	repo.MarkdownLinks = mdscanner.CheckMarkdownLinksWithSleep(&repo, time.Second)
	repo.MarkdownLinks = mdscanner.SortLinksBy404(repo.MarkdownLinks)

	fmt.Println(repo.MarkdownLinks)

}
