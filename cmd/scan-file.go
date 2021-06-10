package cmd

import (
	"fmt"
	"io"
	"markdownscanner/internal/mdscanner"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// scanFileCmd represents the scanFile command
var scanFileCmd = &cobra.Command{
	Use:   "scan-file",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: RunScanfile,
}

//Oh this is mega unfinished

func init() {
	rootCmd.AddCommand(scanFileCmd)
}

func RunScanfile(cmd *cobra.Command, args []string) {
	fmt.Println("ayooooooo :D")

	//Download file

	err := DownloadFile("https://raw.githubusercontent.com/etcd-io/cetcd/master/CONTRIBUTING.md")
	if err != nil {
		fmt.Println("uuuh error", err)
	}

	contrib := mdscanner.MarkdownFile{
		FileName: "CONTRIBUTING.md",
		FilePath: "/tmp/CONTRIBUTING.md",
		HTTPAddr: "https://raw.githubusercontent.com/etcd-io/cetcd/master/CONTRIBUTING.md",
	}

	//func CheckMarkdownLinksWithSleep(mdLinks []MarkdownLink, sleepTime time.Duration) []MarkdownLink {
	mdLinks, err := mdscanner.GetMarkdownLinksFromFile(contrib)
	if err != nil {
		fmt.Println("uuuh error", err)
	}

	scannedLinks := mdscanner.CheckMarkdownLinksWithSleep(mdLinks, time.Second)

	fmt.Println(scannedLinks)

}

//Downloads a file to the current dir
//For the purposes of testing/initial development, let's stick to
//https://raw.githubusercontent.com/etcd-io/cetcd/master/CONTRIBUTING.md
//- Read the [README](README.md) for build and test instructions <- Aaah, this is a file link that gets interpreted as an HTTP link!

//Ah, not sure if a simple implementation is possible here as we have to clone the entire repository to check for file links

func DownloadFile(url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//Ah, we need to
	//Get a raw url from a regular URL
	//Get the file name from the URL and save that to disk
	out, err := os.Create("/tmp" + string(os.PathSeparator) + "CONTRIBUTING.md")
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
