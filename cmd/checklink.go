/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// checklinkCmd represents the checklink command
var checklinkCmd = &cobra.Command{
	Use:   "checklink",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: Checklink,
}

func init() {
	rootCmd.AddCommand(checklinkCmd)
}

func Checklink(cmd *cobra.Command, args []string) {
	fmt.Println("checklink called!")
	fmt.Println("Args:", args)
}
