// Package install implements 'cure install [...]' command
package install

import (
	"context"
	"fmt"
	"strings"

	"github.com/sund3RRR/cure/internal/config"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
)

// Nix interface for interacting with nix
type Nix interface {
	GetPackage(registry, pkg string) string
}

// Install command
type Install struct {
	logger *zap.Logger
	cmd    *cli.Command
	nix    Nix
}

// NewInstall creates new install command
func NewInstall(_ config.Config, logger *zap.Logger, nix Nix) *Install {
	install := &Install{
		logger: logger,
		nix:    nix,
	}

	install.cmd = &cli.Command{
		Name:   "install",
		Usage:  "cure install <package>, if you want install from default registry (nixpkgs) or cure install <registry#package>",
		Action: install.handleCommand,
	}

	return install
}

// Command is the method returning install command
func (install *Install) Command() *cli.Command {
	return install.cmd
}

func (install *Install) handleCommand(_ context.Context, cmd *cli.Command) error {
	install.logger.Debug("Install command call")

	if cmd.NArg() == 0 {
		return ErrPackageNotProvided
	}

	for _, rawPkg := range cmd.Args().Slice() {
		splitted := strings.Split(rawPkg, "#")

		var registry, pkg string
		if len(splitted) == 1 {
			registry, pkg = "nixpkgs", splitted[0]
		} else {
			registry, pkg = splitted[0], splitted[1]
		}

		path := install.nix.GetPackage(registry, pkg)

		fmt.Println("Installed package path: ", path)
	}

	return nil
}
