package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	_outFile = "openapiv2/apidocs.json"
)

type genDocsRunner struct{}

func NewGenDocsCmd() *cobra.Command {
	r := genDocsRunner{}

	var docsCmd = &cobra.Command{
		Use:   "gen:docs",
		Short: "Generate documentation from protobuf files",
		Run:   r.generateDocs,
	}

	docsCmd.Flags().StringP(
		"path",
		"p",
		"./protos/web",
		"Immutable!! Path to proto definitions directory",
		// "Name of the controller for which to generate docs",
	)

	return docsCmd
}

func (r genDocsRunner) generateDocs(cmd *cobra.Command, args []string) {
	// path, err := cmd.Flags().GetString("path")
	// if err != nil {
	// 	cmd.Usage()
	// 	log.Fatalln("failed to get path from options. reason:", err)
	// }

	// if path == "" {
	// 	cmd.Usage()
	// 	log.Fatalln("path is required.")
	// }
	path := "./protos/web"

	if err := r.genAPIProto(path); err != nil {
		log.Fatalln("failed to generate swagger yamls. reason: ", err)
	}

	if err := r.mergeSwaggerFiles(_outFile); err != nil {
		log.Fatalln("failed to merge swagger yamls. reason: ", err)
	}

	log.Println("combined swagger generated at", _outFile)
}

func (r genDocsRunner) buildArgs(protosPath string) []string {
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
	args := []string{"-I", protosPath}

	args = append(args, DefaultProtocArgs...)
	args = append(args,
		"--openapiv2_out", "./openapiv2",
		"--openapiv2_opt", "logtostderr=true",
	)
	// Append all the .proto files found by the glob
	args = append(args, protoFiles...)
	return args
}

func (r genDocsRunner) genAPIProto(protosPath string) error {
	cliArgs := r.buildArgs(protosPath)
	log.Println("protoc", strings.Join(cliArgs, " "))

	return Execute("protoc", cliArgs...)
}

func (r genDocsRunner) mergeSwaggerFiles(outFile string) error {
	files, err := filepath.Glob("./openapiv2/**/*.swagger.json")
	if err != nil {
		return fmt.Errorf("failed to find swagger files: %v", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no swagger files found")
	}

	// Join the files for swagger mixin
	filesStr := strings.Join(files, " ")

	if len(files) == 1 {
		if err := exec.Command("cp", filesStr, outFile).Run(); err != nil {
			return fmt.Errorf("failed to create output file: %v", err)
		}

		return nil
	}

	// Build the swagger mixin command
	cmd := exec.Command("swagger", "mixin", filesStr)

	// Set the output file
	output, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer output.Close()

	// Redirect command output to the file
	cmd.Stdout = output

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	log.Println("Running swagger mixin command...")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to merge swagger files: %v\n%s", err, stderr.String())
		return err
	}

	return nil
}
