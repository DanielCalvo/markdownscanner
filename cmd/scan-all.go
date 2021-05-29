package cmd

import (
	"fmt"
	"log"
	"markdownscanner/internal/config"
	"markdownscanner/internal/mdscanner"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// scanAllCmd represents the scanAll command
var scanAllCmd = &cobra.Command{
	Use:   "scan-all",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: RunScanAll,
}

func init() {
	rootCmd.AddCommand(scanAllCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanAllCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	scanAllCmd.Flags().String("config", "", "Path to mdscanner config file")
	err := scanAllCmd.MarkFlagRequired("config")

	//This is pointless as MarkFlagRequired seems to panic for you
	if err != nil {
		fmt.Println("config flag missing!")
	}
}

func RunScanAll(cmd *cobra.Command, args []string) {
	fmt.Println("scan-aaaall called")

	//I want to print the value of the config flag, hmm
	//If path for the config file is not set, look in the current directory!
	fmt.Println(cmd.Flag("config").Value)

	//If something goes wrong, do I panic here? I suppose so /shrug

	//go run main.go scan-all --config=mdscanner.yaml
	conf, err := config.LoadFile(cmd.Flag("config").Value.String())
	if err != nil {
		//Print path of the config file that was tried to load
		log.Println("Could not load config file:", err)
		os.Exit(1)
	}

	log.Println("Initializing config file")
	err = config.Initialize(conf)
	if err != nil {
		log.Println("Error initiating config:", err)
		os.Exit(1)
	}

	fmt.Println(conf)

	//This is questionable -- I want this to be simpler!
	repoUrlsFromProjects := mdscanner.GetRepoUrlsFromProjects(conf.GithubProjects)
	repoUrlsToScan := append(conf.Repositories, repoUrlsFromProjects...)
	repositories := mdscanner.NewRepositories(conf, repoUrlsToScan)
	repositories = mdscanner.SortRepositoriesByUnscannedFirst(repositories)

	for _, repo := range repositories {

		log.Println("Cloning repostory:", repo.Name)
		err = mdscanner.CloneRepository(repo.URL, repo.FilesystemPath)
		if err != nil {
			log.Println("Error clonning repository:", err)
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
			continue
		}

		//Figure out later why this doesn't work in place
		log.Println("Checking Markdown links")
		repo.MarkdownLinks = mdscanner.CheckMarkdownLinksWithSleep(repo.MarkdownLinks, time.Second)
		repo.MarkdownLinks = mdscanner.SortLinksBy404(repo.MarkdownLinks)

		//The ordering of this is very strange, you should adjust it!
		log.Println("Uploading scan result to S3")
		err = mdscanner.UploadResultsToS3(*conf, repo)

		if err != nil {
			log.Println("Could not upload to S3:", err)
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

	//copy and paste everything on main.go here all the way to S3 upload
	//Do you want a flag named --s3-upload which defaults to true? Hmm

}
