//Hang on let me just try testing connecting to S3 to learn a thing or two

package cmd

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"

	"github.com/spf13/cobra"
)

// s3testCmd represents the s3test command
var s3testCmd = &cobra.Command{
	Use:   "s3test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: S3test,
}

func init() {
	rootCmd.AddCommand(s3testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// s3testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// s3testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func S3test(cmd *cobra.Command, args []string) {
	fmt.Println("Hello world from s3 test!")
	//var S3session *session.Session
	S3session, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-1")})
	if err != nil {
		log.Fatalln(err)
	}
	//input := &s3.ListBucketsInput{}

	output, err := s3.New(S3session).PutObject(&s3.PutObjectInput{
		Bucket:      aws.String("mdscanner-bucket"),
		Key:         aws.String("mdscanner-testfile"),
		ACL:         aws.String("public-read"),
		ContentType: aws.String("text/html"),
		Body:        bytes.NewReader([]byte("Hello world!")),
	})
	if err != nil {
		log.Fatalln("uh it error out putting object:", err)
	}
	fmt.Println(output)
}
