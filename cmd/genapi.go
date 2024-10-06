package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func NewGenControllerCmd() *cobra.Command {
	var controllerCmd = &cobra.Command{
		Use:   "controller",
		Short: "Generate a new controller",
		Run:   generateController,
	}

	controllerCmd.Flags().StringP(
		"name",
		"n",
		"",
		"Name of the controller",
	)

	controllerCmd.Flags().StringP(
		"path",
		"p",
		"app/core/controllers",
		"Path to controllers",
	)

	return controllerCmd
}

func generateController(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		fmt.Println("Controller name is required.")
		return
	}

	controllerDir := filepath.Join("app/core/controllers", name)
	os.MkdirAll(controllerDir, os.ModePerm)

	// Generate proto files and OpenAPI annotations
	protoFile := filepath.Join(controllerDir, fmt.Sprintf("%s.proto", name))
	generateProtoFile(name, protoFile)

	// Other files generation logic...

	fmt.Printf("Generated controller: %s\n", name)
}

func generateProtoFile(name, protoFile string) {
	// Template and file creation logic...
	fmt.Printf("Generated proto file: %s\n", protoFile)
}
