// Package commands implements cure cli commands
package commands

import (
	"context"
	"fmt"

	"github.com/sund3RRR/cure/internal/config"
	"github.com/sund3RRR/cure/pkg/modules/install"
	"github.com/sund3RRR/cure/pkg/types"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
)

// Nix interface for interacting with nix
type Nix interface {
	GetPackage(registry, pkg string) (types.PackageInfo, error)
	PathInfo(registry, pkg string) (types.PackageInfo, error)
	AddRegistry(alias, registry string) error
}

// Install command
type Install struct {
	logger    *zap.Logger
	cmd       *cli.Command
	nix       Nix
	installer *install.Installer
}

// NewInstall creates new install command
func NewInstall(_ config.Config, logger *zap.Logger, nix Nix) *Install {
	install := &Install{
		logger:    logger,
		nix:       nix,
		installer: install.NewInstaller(nix),
	}

	install.cmd = &cli.Command{
		Name: "install",
		Usage: "command to install package from nix flake registry.\n" +
			"Type 'cure install <package>', if you want install from default registry (nixpkgs) or 'cure install <registry#package>'",
		Action: install.handleCommand,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "gl",
				Value: types.NixGLAuto.String(),
				Usage: "enable nixGL wrapper for package.\n" +
					possibleValuesString(types.NixGLValues, types.NixGLAuto, types.UnderlinedSeq) + "\n" +
					"Auto will automatically detect GUI binaries and enable nixGL wrapper for them",
				Action: func(_ context.Context, _ *cli.Command, s string) error {
					return validateFlagValue(types.NixGL(s), types.NixGLValues)
				},
			},
			&cli.StringFlag{
				Name:    "gl-type",
				Aliases: []string{"glt"},
				Value:   types.NixGLAuto.String(),
				Usage: "nixGL wrapper type for package.\n" +
					possibleValuesString(types.NixGLPackageValues, types.NixGLPackageAuto, types.UnderlinedSeq) + "\n" +
					"Auto will automatically detect needed nixGL wrapper",
				Action: func(_ context.Context, _ *cli.Command, s string) error {
					return validateFlagValue(types.NixGLPackage(s), types.NixGLPackageValues)
				},
			},
		},
	}

	return install
}

// Command is the method returning install command
func (ins *Install) Command() *cli.Command {
	return ins.cmd
}

func (ins *Install) handleCommand(_ context.Context, cmd *cli.Command) error {
	ins.logger.Debug("Install command call")

	if cmd.NArg() == 0 {
		return ErrPackageNotProvided
	}

	for _, pkg := range cmd.Args().Slice() {
		err := ins.installer.InstallPackage(pkg, install.Params{
			NixGL:        types.NixGL(cmd.String("gl")),
			NixGLPackage: types.NixGLPackage(cmd.String("gl-type")),
		})
		if err != nil {
			ins.logger.Error("failed to install package", zap.Error(err))
			continue
		}
		fmt.Println("Successfuly installed package ", pkg)
	}

	return nil
}
