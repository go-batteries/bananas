package cmd

import (
	"embed"
	"log"
	"os"
	"strings"

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
	"protos",
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
			data:         bananas.TemplData{projectName: projectName},
			fs:           r.baseFS,
		},
		"cli": {
			tmplFilePath: "templates/cmd/cli.go.tmpl",
			data: bananas.TemplData{
				projectName: projectName,
				binaryName:  binaryName,
			},
			fs: r.baseFS,
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
