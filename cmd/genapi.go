package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type genApiRunner struct{}

func NewGenApiCmd() *cobra.Command {
	r := genApiRunner{}

	var controllerCmd = &cobra.Command{
		Use:   "gen:controllers",
		Short: "Generate all go structs from proto definitions",
		Run:   r.generateApiCode,
	}

	controllerCmd.Flags().StringP(
		"path",
		"p",
		"protos/web",
		"Path to protos directory",
	)

	return controllerCmd
}

func (r genApiRunner) generateApiCode(cmd *cobra.Command, args []string) {
	path, err := cmd.Flags().GetString("path")
	if err != nil {
		cmd.Usage()
		log.Fatalln("failed to get path from options. reason:", err)
	}

	if path == "" {
		cmd.Usage()
		log.Fatalln("path is required.")
	}

	err = r.genWebProto(strings.TrimSuffix(path, "/"))
	if err != nil {
		log.Fatalln("failed to generate pb go files. reason:", err)
	}

	log.Println("generated pb files")
}

func (r genApiRunner) buildArgs(protosPath string) []string {
	protoFiles, err := filepath.Glob(fmt.Sprintf("%s/**/*.proto", protosPath))
	if err != nil {
		log.Println("failed to find .proto files. reason:", err)
		return []string{}
	}
	if len(protoFiles) == 0 {
		log.Println("no .proto files found in path", protosPath)
		return []string{}
	}

	// Build the protoc command with all .proto files
	args := []string{
		"-I", protosPath,
		"-I", "protos/includes/googleapis",
		"-I", "protos/includes/grpc_ecosystem",
		"--go_out=./app/core", "--go_opt=paths=source_relative",
		"--go-grpc_out=./app/core", "--go-grpc_opt=paths=source_relative",
	}

	args = append(args, protoFiles...)
	return args
}

func (r genApiRunner) genWebProto(protosPath string) error {
	cmd := exec.Command("protoc", r.buildArgs(protosPath)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
