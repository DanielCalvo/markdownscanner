package cmd

import (
	"log"
	"markdownscanner/internal/config"
	"markdownscanner/internal/mdscanner"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var scanAllCmd = &cobra.Command{
	Use:   "scan-all",
	Short: "scans all repositories as defined in the yaml config file",
	Long: `scans all repositories as defined in the yaml config file
Will look for mdscanner.yaml in the current directory
The path to this file can also be specified, ex: mdscanner --config=/home/user/mdscanner.yaml`,
	Run: RunScanAll,
}

func init() {
	rootCmd.AddCommand(scanAllCmd)

	pwd, err := os.Getwd()
	if err != nil {
		log.Panicln("Unable to get working dir:", err)
	}

	scanAllCmd.Flags().String("config", pwd+string(os.PathSeparator)+"mdscanner.yaml", "Path to mdscanner config file")
}

func RunScanAll(cmd *cobra.Command, args []string) {

	//This needs redoing/refining
	conf, err := config.New(cmd.Flag("config").Value.String())
	if err != nil {
		log.Panicln("Could not load initialize config:", err)
	}

	//These 3 seem to be a repeating pattern -- can you put this in a function perhaps?
	repoUrlsFromProjects := mdscanner.GetRepoUrlsFromProjects(conf.GithubProjects)
	repoUrlsToScan := append(conf.Repositories, repoUrlsFromProjects...)
	repositories := mdscanner.NewRepositories(&conf, repoUrlsToScan)

	//Try storing repository scan state somewhere persistent (S3?) but not on the local filesystem, as that's ephemeral
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

		log.Println("Deleting cloned repository files for", repo.Name)
		err = mdscanner.DeleteRepository(repo)
		if err != nil {
			log.Println("Unable to delete repository:", err)
			continue
		}

		//Figure out later why this doesn't work in place
		log.Println("Checking Markdown links")
		repo.MarkdownLinks = mdscanner.CheckMarkdownLinksWithSleep(repo.MarkdownLinks, time.Second)
		repo.MarkdownLinks = mdscanner.SortLinksBy404(repo.MarkdownLinks)

		//The ordering of the operations with S3 is very strange, you should adjust it!
		log.Println("Uploading scan result to S3")
		err = mdscanner.UploadResultsToS3(conf, repo)

		if err != nil {
			log.Println("Could not upload to S3:", err)
			continue
		}

		//Can you put this into a function to make it prettier?
		repo.LastScanned = time.Now().Format(time.RFC3339)
		repo.LinksScanned = len(repo.MarkdownLinks)
		repo.Links404 = mdscanner.Count404MarkdownLinks(repo.MarkdownLinks)

		log.Println("Saving scan metadata")
		err = mdscanner.SaveScanMetadata(repo)
		if err != nil {
			log.Println("Could not SaveScanMetadata:", err)
		}

		log.Println("Generating and uploading index html")
		err = mdscanner.GenerateAndUploadIndexHtml(conf)
		if err != nil {
			log.Println("Could not GenerateIndexHtml:", err)
		}
	}
}
