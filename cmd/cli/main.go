package main

import (
	"log"
	"os"

	"github.com/go-batteries/bananas/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bananas",
	Short: "Bananas framework for Go applications",
}

func main() {
	log.SetFlags(0)

	rootCmd.AddCommand(cmd.NewInitAppCmd())
	rootCmd.AddCommand(cmd.NewGenControllerCmd())
	rootCmd.AddCommand(cmd.NewGenDocsCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

