package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"markdownscanner/internal/mdscanner"
	"os"
	"time"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan will check all the markdown links for a given github repository on a remote url",
	Long: `scan will download a repository and check all it's links
ex: markdownscanner https://github.com/kubernetes/kubectl
In the case above, the results will be saved to a file named kubernetes_kubectl.json
The scan results can also be sent to stdout with the '--stdout=true' setting`,
	Run: RunScan,
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().Bool("stdout", false, "If the output is to be printed to stdout instead of saved to a file. Default is false")
}

func RunScan(cmd *cobra.Command, args []string) {
	for _, repo := range args {
		repoUrlPath, err := mdscanner.GetUrlPath(repo)
		if err != nil {
			log.Println(repo, "is not a valid URL")
			continue
		}

		//This struct was created with the scan-all.go use case in mind. It's painful to instantiate it like this
		repo := mdscanner.Repository{
			URL:                repo, //url is repo and name is repourl? That's confusing!
			Name:               repoUrlPath,
			FilesystemPath:     "/tmp/" + repoUrlPath,
			LastScanned:        "",
			LinksScanned:       0,
			Links404:           0,
			JSONReportPath:     "",
			HTMLReportPath:     "",
			MetadataReportPath: "",
			MarkdownFiles:      nil,
			MarkdownLinks:      nil,
		}

		log.Println("Clonning:", repo)
		err = mdscanner.CloneRepository(repo.URL, repo.FilesystemPath)
		if err != nil {
			log.Println("Could not clone" + repo.URL + " to " + repo.FilesystemPath)
			continue
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

		log.Println("Checking Markdown links")
		repo.MarkdownLinks = mdscanner.CheckMarkdownLinksWithSleep(repo.MarkdownLinks, time.Second)
		repo.MarkdownLinks = mdscanner.SortLinksBy404(repo.MarkdownLinks)

		currentDir, err := os.Getwd()
		if err != nil {
			log.Panicln("Unable to get current directory:", err)
		}

		if cmd.Flag("stdout").Value.String() == "false" {
			err = mdscanner.SaveStructToJson(repo, currentDir+string(os.PathSeparator)+mdscanner.GetURLWithUnderscores(repo.Name)+".json") //to be fixed
			if err != nil {
				log.Panicln(err)
			}
		} else {
			repoJSON, err := json.MarshalIndent(repo, "", "  ")
			if err != nil {
				log.Fatalf(err.Error())
			}
			fmt.Println(string(repoJSON))
		}
	}
}
