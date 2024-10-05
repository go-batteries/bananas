package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewGenDocsCmd() *cobra.Command {
	var docsCmd = &cobra.Command{
		Use:   "docs",
		Short: "Generate documentation from protobuf files",
		Run:   generateDocs,
	}

	docsCmd.Flags().StringP(
		"name",
		"n",
		"",
		"Name of the controller for which to generate docs",
	)

	return docsCmd
}

func generateDocs(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")

	// Merge swagger files logic...
	if name != "" {
		fmt.Printf("Generating docs for controller: %s\n", name)
	} else {
		fmt.Println("Generating docs for all controllers.")
	}
}
