package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type genApiRunner struct{}

func NewGenApiCmd() *cobra.Command {
	r := genApiRunner{}

	var controllerCmd = &cobra.Command{
		Use:   "gen:structs",
		Short: "Generate all go structs from proto definitions",
		Run:   r.generateApiCode,
	}

	controllerCmd.Flags().StringP(
		"path",
		"p",
		"./protos/web",
		"Immutable!! Path to proto definitions directory",
	)

	controllerCmd.Flags().StringP(
		"out_path",
		"o",
		"./app/core",
		"Output path for compiled proto definitions",
	)

	controllerCmd.Flags().BoolP(
		"grpc",
		"",
		false,
		"Turn on the flag to enable grpc gateway",
	)
	return controllerCmd
}

func (r genApiRunner) generateApiCode(cmd *cobra.Command, args []string) {
	// path, err := cmd.Flags().GetString("path")
	// if err != nil {
	// 	cmd.Usage()
	// 	log.Fatalln("failed to get path from options. reason:", err)
	// }
	//
	// if path == "" {
	// 	cmd.Usage()
	// 	log.Fatalln("path is required.")
	// }

	outPath, err := cmd.Flags().GetString("out_path")
	if err != nil {
		cmd.Usage()
		log.Fatalln("failed to get path from options. reason:", err)
	}

	if outPath == "" {
		cmd.Usage()
		log.Fatalln("out_path is required.")
	}

	path := "./protos/web"

	isgRPCMode, err := cmd.Flags().GetBool("grpc")
	if err != nil {
		isgRPCMode = false
	}

	err = r.genWebProto(strings.TrimSuffix(path, "/"), outPath, isgRPCMode)
	if err != nil {
		log.Fatalln("failed to generate pb go files. reason:", err)
	}

	log.Println("generated pb files")
}

func (r genApiRunner) buildArgs(protosPath string, overrideOutPath string, isgRPCMode bool) []string {
	protoFiles, err := filepath.Glob(fmt.Sprintf("%s/**/*.proto", protosPath))
	if err != nil {
		log.Println("failed to find .proto files. reason:", err)
		return []string{}
	}
	if len(protoFiles) == 0 {
		log.Println("no .proto files found in path", protosPath)
		return []string{}
	}

	outPath, err := filepath.Abs(overrideOutPath)
	if err != nil {
		log.Fatal("failed to find directory, init first")
	}
	// Build the protoc command with all .proto files
	args := []string{
		"-I", protosPath,
		"-I", "protos/includes/googleapis",
		"-I", "protos/includes/grpc_ecosystem",
		fmt.Sprintf("--go_out=%s", outPath), "--go_opt=paths=import",
		fmt.Sprintf("--go-grpc_out=%s", outPath), "--go-grpc_opt=paths=import",
	}

	if isgRPCMode {
		args = append(args,
			fmt.Sprintf("--grpc-gateway_out=%s", outPath),
			"--grpc-gateway_opt=paths=import",
		)
	}

	args = append(args, protoFiles...)
	return args
}

func (r genApiRunner) genWebProto(protosPath string, outPath string, isgRPCMode bool) error {
	cliArgs := r.buildArgs(protosPath, outPath, isgRPCMode)
	log.Println("protoc", strings.Join(cliArgs, " "))

	return Execute("protoc", cliArgs...)
}
