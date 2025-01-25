// cure
package main

import (
	"context"
	"log"
	"os"
	"os/user"

	"github.com/sund3RRR/cure/internal/app"
	"github.com/sund3RRR/cure/internal/config"
	"github.com/urfave/cli/v3"
)

var home = getUserHome()
var configPaths = []string{home + "/.config/cure/cure.yaml", "/etc/cure/cure.yaml"}

func main() {
	ctx := context.Background()

	cfg := config.NewConfig(configPaths...)

	logger, err := cfg.Logger.Build()
	if err != nil {
		log.Fatal("failed to create logger: ", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatal(err)
		}
	}()

	cmd := &cli.Command{
		Name:     "cure",
		Usage:    "a package manager",
		Commands: app.NewApp(cfg, logger).GetCommands(),
	}

	if err := cmd.Run(ctx, os.Args); err != nil {
		log.Fatal("failed to start cmd: ", err)
	}
}

func getUserHome() string {
	u, err := user.Current()
	if err != nil {
		log.Fatal("failed to get user information: ", err)
	}

	return u.HomeDir
}
