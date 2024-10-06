package cmd

import (
	"embed"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-batteries/bananas"
	"github.com/spf13/cobra"
)

var dirs = []string{
	"app/web/middlewares",
	"app/web/routers",
	"app/core",
	"cmd/server",
	"cmd/cli",
	"cmd/workers",
	"pkg/config",
	"pkg/httputils",
	"pkg/workerpool",
	"pkg/databases/dbsqlite",
	"migrations/sqlite",
	"migrations/mysql",
	"migrations/postgres",
	"config",
	"protos/web",
	"openapiv2",
}

type appInitRunner struct {
	dirs       []string
	baseFS     embed.FS
	databaseFS embed.FS
	pkgFS      embed.FS
}

func NewInitAppCmd() *cobra.Command {
	r := appInitRunner{
		dirs:       dirs,
		baseFS:     bananas.BaseFS,
		databaseFS: bananas.DbFS,
		pkgFS:      bananas.PkgFS,
	}

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Bananas application",
		Run:   r.initApp,
	}
	initCmd.Flags().StringP(
		"mode",
		"m",
		"http",
		"Specify the mode (http or grpc)",
	)

	initCmd.Flags().StringP(
		"name",
		"n",
		"",
		"Specify the project name as per go.mod",
	)

	return initCmd
}

func getBinaryName(projectName string) string {
	splits := strings.Split(projectName, "/")
	idx := len(splits) - 1

	return splits[idx]
}

func (r appInitRunner) initApp(cmd *cobra.Command, args []string) {
	mode, err := cmd.Flags().GetString("mode")
	if err != nil {
		log.Fatalf("mode is empty. %v\n", err)
	}

	appName, err := cmd.Flags().GetString("name")
	if err != nil || appName == "" {
		cmd.Usage()
		log.Fatal("\napp name not provided.")
	}

	{
		// Create the base directory structure
		log.Println("creating directories..")

		for _, dir := range dirs {
			log.Println(dir)
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				log.Fatalf("Error creating directory: %v\n", err)
			}
		}

		log.Println("directory setup done..")
	}

	{
		// Copy Templates and shit for base project
		log.Println("setting up base project..")

		r.copyTemplates(appName, mode)

		log.Println("setting up base done..")
	}

	{
		// Setup Required Protos
		err := r.setupRequiredProtos()
		if err != nil {
			log.Fatal("project setup failed for", appName, "reason:", err)
		}
	}

	log.Printf("\nInitialized a new Bananas app %s in '%s' mode.\n", appName, mode)
}

type render struct {
	tmplFilePath string
	data         bananas.TemplData
	fs           embed.FS
}

func (r appInitRunner) copyTemplates(projectName, mode string) {
	binaryName := getBinaryName(projectName)

	renderers := map[string]render{
		"server": {
			tmplFilePath: "templates/cmd/server.go.tmpl",
			data:         bananas.TemplData{"projectName": projectName},
			fs:           r.baseFS,
		},
		"cli": {
			tmplFilePath: "templates/cmd/cli.go.tmpl",
			data: bananas.TemplData{
				"projectName": projectName,
				"binaryName":  binaryName,
			},
			fs: r.baseFS,
		},
		"tools": {
			tmplFilePath: "templates/cmd/tools.go.tmpl",
			fs:           r.baseFS,
		},
		"app.env": {
			tmplFilePath: "templates/cmd/app.env.tmpl",
			fs:           r.baseFS,
		},

		"config": {
			tmplFilePath: "templates/pkg/config.go.tmpl",
			fs:           r.pkgFS,
		},
		"httpclient": {
			tmplFilePath: "templates/pkg/httpclient.go.tmpl",
			fs:           r.pkgFS,
		},
		"workerpool": {
			tmplFilePath: "templates/pkg/workerpool.go.tmpl",
			fs:           r.pkgFS,
		},

		"sqlxdbconn": {
			tmplFilePath: "templates/databases/dbconnect.sqlx.go.tmpl",
			fs:           r.databaseFS,
		},
		"redisconn": {
			tmplFilePath: "templates/databases/redisconn.go.tmpl",
			fs:           r.databaseFS,
		},
		"crypter": {
			tmplFilePath: "templates/databases/crypter.go.tmpl",
			fs:           r.databaseFS,
		},
		"idgen": {
			tmplFilePath: "templates/databases/idgen.go.tmpl",
			fs:           r.databaseFS,
		},
		"timestamper": {
			tmplFilePath: "templates/databases/timestamper.go.tmpl",
			fs:           r.databaseFS,
		},
	}

	for _, r := range renderers {
		data := bananas.TemplData{}

		if len(r.data) != 0 {
			data = r.data
		}

		outPath, content, ok := bananas.MustRenderTmpl(
			r.fs,
			r.tmplFilePath,
			data,
		)
		if ok {
			log.Println(outPath)
			bananas.WriteFile(outPath, content)
		}
	}

}

func (r appInitRunner) setupRequiredProtos() error {
	googleApiDirRoot := "protos/includes/googleapis"
	grpcEcosystemDirRoot := "protos/includes/grpc_ecosystem/protoc-gen-openapiv2"

	// Create directories
	err := os.MkdirAll(filepath.Join(googleApiDirRoot, "google/api"), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(googleApiDirRoot, "google/protobuf"), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(grpcEcosystemDirRoot, "options"), os.ModePerm)
	if err != nil {
		return err
	}

	// Download proto files using HTTP
	var protos = []struct {
		url  string
		path string
	}{
		{
			url:  "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto",
			path: filepath.Join(googleApiDirRoot, "google/api/http.proto"),
		},
		{
			url:  "https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto",
			path: filepath.Join(googleApiDirRoot, "google/api/annotations.proto"),
		},
		{
			url:  "https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/descriptor.proto",
			path: filepath.Join(googleApiDirRoot, "google/protobuf/descriptor.proto"),
		},
		{
			url:  "https://raw.githubusercontent.com/protocolbuffers/protobuf/refs/heads/main/src/google/protobuf/empty.proto",
			path: filepath.Join(googleApiDirRoot, "google/protobuf/empty.proto"),
		},
		{
			url:  "https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/struct.proto",
			path: filepath.Join(googleApiDirRoot, "google/protobuf/struct.proto"),
		},
		{
			url:  "https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/main/protoc-gen-openapiv2/options/annotations.proto",
			path: filepath.Join(grpcEcosystemDirRoot, "options/annotations.proto"),
		},
		{
			url:  "https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/main/protoc-gen-openapiv2/options/openapiv2.proto",
			path: filepath.Join(grpcEcosystemDirRoot, "options/openapiv2.proto"),
		},
	}

	var wg = sync.WaitGroup{}
	wg.Add(len(protos))

	for i := 0; i < len(protos); i++ {
		go func(_wg *sync.WaitGroup, idx int) {
			defer _wg.Done()

			proto := protos[idx]

			err := downloadProtoFile(proto.url, proto.path)
			if err != nil {
				log.Println("failed to download protos from", proto.url)
			}
		}(&wg, i)
	}

	wg.Wait()
	// err = downloadProtoFile("https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto", filepath.Join(googleApiDirRoot, "google/api/http.proto"))
	// if err != nil {
	// 	return err
	// }
	//
	// err = downloadProtoFile("https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto", filepath.Join(googleApiDirRoot, "google/api/annotations.proto"))
	// if err != nil {
	// 	return err
	// }

	// err = downloadProtoFile("https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/descriptor.proto", filepath.Join(googleApiDirRoot, "google/protobuf/descriptor.proto"))
	// if err != nil {
	// 	return err
	// }

	// err = downloadProtoFile("https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/struct.proto", filepath.Join(googleApiDirRoot, "google/protobuf/struct.proto"))
	// if err != nil {
	// 	return err
	// }

	// err = downloadProtoFile("https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/main/protoc-gen-openapiv2/options/annotations.proto", filepath.Join(grpcEcosystemDirRoot, "options/annotations.proto"))
	// if err != nil {
	// 	return err
	// }

	// err = downloadProtoFile("https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/main/protoc-gen-openapiv2/options/openapiv2.proto", filepath.Join(grpcEcosystemDirRoot, "options/openapiv2.proto"))
	// if err != nil {
	// 	return err
	// }

	return nil
}

func downloadProtoFile(url, outputPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("failed to download %s: %v", url, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("failed to download %s: server returned %d", url, resp.StatusCode)
		return err
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("failed to create file %s: %v", outputPath, err)
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		log.Fatalf("failed to save file %s: %v", outputPath, err)
		return err
	}

	log.Printf("downloaded %s to %s\n", url, outputPath)
	return nil
}
