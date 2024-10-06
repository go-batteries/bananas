package main

import (
	dbmgr "testgrpcsetup/pkg/databases/dbsqlite"
	"testgrpcsetup/pkg/config"
	"context"
	"errors"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

const (
	EnvDirCmd = "envdir"
	Direction = "dir"
)

type MigrateCmd struct{}

func (mc MigrateCmd) Run(c *cli.Context) error {
	dir := c.String(Direction)
	envDir := c.String(EnvDirCmd)

	if dir != "up" && dir != "down" {
		return errors.New("invalid command direction")
	}

	cfg := config.Load(envDir)

	conn := dbmgr.ConnectSqlite(cfg.DbName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn.Connect(ctx)

	ctx = context.Background()

	if err := conn.Setup(ctx, dir == "up"); err != nil {
		log.Error().Err(err).Msg("failed to setup db schema")
		return err
	}

	log.Info().Msg("database schema migration success")
	return nil
}

func main() {
	migrateCmd := MigrateCmd{}

	app := &cli.App{
		Name:  "arbok",
		Usage: "cli app for testgrpcsetup",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    EnvDirCmd,
				Aliases: []string{"e"},
				Value:   "./config",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "migrate",
				Usage: "testgrpcsetup migrate -dir [up/down]",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: Direction, Required: true},
				},
				Action: migrateCmd.Run,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("command failed")
	}
}
