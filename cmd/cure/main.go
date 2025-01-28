// cure
package main

import (
	"context"
	"log"
	"os"
	"os/user"

	"github.com/sund3RRR/cure/internal/app"
	"github.com/sund3RRR/cure/internal/config"
	"github.com/sund3RRR/cure/internal/modules/install"
	"github.com/sund3RRR/cure/pkg/adapters/nix"
	"github.com/urfave/cli/v3"
)

var home = getUserHome()
var configPaths = []string{home + "/.config/cure/cure.yaml", "/etc/cure/cure.yaml"}

func main() {
	// Create main application context
	ctx := context.Background()

	// Create main config
	cfg := config.NewConfig(configPaths...)

	// Init logger
	logger, err := cfg.Logger.Build()
	if err != nil {
		log.Fatal("failed to create logger: ", err)
	}
	defer logger.Sync() //nolint

	// Create adapters
	nixAdapter := nix.NewNix(logger)

	// Create commands
	installCmd := install.NewInstall(cfg, logger, nixAdapter)

	cmd := &cli.Command{
		Name:     "cure",
		Usage:    "a package manager",
		Commands: app.NewApp(installCmd).GetCommands(),
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
