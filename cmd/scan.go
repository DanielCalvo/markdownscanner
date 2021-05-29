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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: RunScan, //I have no idea if this is a good name but I'm just gonna roll with it
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	scanCmd.Flags().String("stdout", "false", "If the output is to be printed to stdout instead of saved to a file. Default is false")
}

func RunScan(cmd *cobra.Command, args []string) {
	fmt.Println("scan caaaaaaaaaalled", args)

	//Usage: mdscanner scan <repo-url>
	//Intended functionality: Save a json file to the local directory with the links checked
	//--stdout: print result to stdout (default)
	//--output-file: save json to given file

	//You need to handle arguments

	for _, repo := range args {
		fmt.Println(repo)
		repoUrlPath, err := mdscanner.GetUrlPath(repo)
		if err != nil {
			log.Println(repo, "is not a valid URL")
			continue
		}

		fmt.Println("repoUrlPath:", repoUrlPath)

		repo := mdscanner.Repository{
			URL:                repo, //Atrocious naming. url is repo and name is repourl? what were you thinking man?
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

		err = mdscanner.CloneRepository(repo.URL, repo.FilesystemPath)
		if err != nil {
			fmt.Println("Could not clone" + repo.URL + " to " + repo.FilesystemPath)
		}

		log.Println("Getting Markdown files")
		mdscanner.GetMarkdownFiles(&repo)
		log.Println("Getting Markdown links from files")
		mdscanner.GetMarkdownLinksFromFiles(&repo)

		fmt.Println(repo)

		log.Println("Deleting repository", repo.Name)
		err = mdscanner.DeleteRepository(repo)
		if err != nil {
			log.Println("Unable to delete repository:", err)
		}

		log.Println("Checking Markdown links")
		repo.MarkdownLinks = mdscanner.CheckMarkdownLinksWithSleep(repo.MarkdownLinks, time.Second)
		repo.MarkdownLinks = mdscanner.SortLinksBy404(repo.MarkdownLinks)

		//save result to current dir
		//or print it to stdout

		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalf(err.Error())
		}

		//default
		err = mdscanner.SaveStructToJson(repo, currentDir+string(os.PathSeparator)+mdscanner.GetURLWithUnderscores(repo.Name)+".json") //to be fixed
		if err != nil {
			log.Fatalf(err.Error())
		}

		//Uuuuh this should be a boolean
		if cmd.Flag("stdout").Value.String() == "true" {
			repoJSON, err := json.MarshalIndent(repo, "", "  ")
			if err != nil {
				log.Fatalf(err.Error())
			}
			fmt.Println(string(repoJSON))
		}

		//implement print to stdout here!

		//repo: https://github.com/kubernetes/kubectl

		//for _, repo := range args {
		//	fmt.Println(repo)
		//}
		//
		//err := mdscanner.CloneRepository("https://github.com/DanielCalvo/github-actions-shenanigans")
		//if err != nil {
		//	fmt.Println("Could not clone") //Args go here!
		//}

		//Clone all repos
		//???? (Scan links?)
		//Delete all repos
	}
}
