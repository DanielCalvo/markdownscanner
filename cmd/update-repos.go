/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"markdownscanner/internal/mdscanner"
	"os"
)

// updateReposCmd represents the updateRepos command
var updateReposCmd = &cobra.Command{
	Use:   "update-repos",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: RunUpdateRepos,
}

func init() {
	rootCmd.AddCommand(updateReposCmd)

	pwd, err := os.Getwd()
	if err != nil {
		log.Panicln("Unable to get working dir:", err)
	}

	updateReposCmd.Flags().String("config", pwd+string(os.PathSeparator)+"config.yaml", "Path to mdscanner config file")
}

func RunUpdateRepos(cmd *cobra.Command, args []string) {
	fmt.Println("yeeeeeah update repos!")
	conf, err := mdscanner.New(cmd.Flag("config").Value.String())
	if err != nil {
		log.Panicln("Could not load initialize config:", err)
	}

	//These 3 seem to be a repeating pattern -- can you put this in a function perhaps?
	repoUrlsFromProjects := mdscanner.GetRepoUrlsFromProjects(conf.GithubProjects) //What projects? Github projects you mean...
	repoUrlsToScan := append(conf.Repositories, repoUrlsFromProjects...)
	repositories := mdscanner.NewRepositories(&conf, repoUrlsToScan)

	//Try storing repository scan state somewhere persistent (S3?) but not on the local filesystem, as that's ephemeral
	repositories = mdscanner.SortRepositoriesByUnscannedFirst(repositories)

	for _, repo := range repositories {
		log.Println("Fetching changes from git for repostory:", repo.Name)
		err = mdscanner.CloneRepository(repo.URL, repo.FilesystemPath)
		if err != nil {
			log.Println("Error clonning repository:", err)
			continue
		}
	}

}
