// Package nix is an adapter to nix command
package nix

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"go.uber.org/zap"
)

// Nix implements and adapter to nix command
type Nix struct {
	logger  *zap.Logger
	command string
	flags   []string
}

// NewNix creates new nix adapter
func NewNix(logger *zap.Logger, flags ...string) *Nix {
	flags = append(
		flags,
		"--extra-experimental-features", "nix-command",
		"--extra-experimental-features", "flakes",
	)

	return &Nix{
		logger:  logger,
		command: "nix",
		flags:   flags,
	}
}

// GetPackage takes provided 'registry' and 'pkg', downloads it to /nix/store and
// returns package path
func (nix *Nix) GetPackage(registry, pkg string) string {
	args := append(nix.flags, "build", "--no-link", nix.compilePackage(registry, pkg))

	cmd := exec.Command(nix.command, args...) //nolint

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	nix.logger.Debug(nix.command + " " + strings.Join(args, " "))

	err := cmd.Run()
	if err != nil {
		nix.logger.Warn("failed to execute nix build", zap.Error(err))
		return ""
	}

	return nix.PathInfo(registry, pkg)
}

// PathInfo takes provided 'registry' and 'pkg' and returns package path
// in the /nix/store. Returns empty string if there is no such package.
func (nix *Nix) PathInfo(registry, pkg string) string {
	args := append(nix.flags, "eval", nix.compilePackage(registry, pkg)+".outPath")

	cmd := exec.Command(nix.command, args...) //nolint

	output, err := cmd.CombinedOutput()
	if err != nil {
		nix.logger.Error("failed to execute nix path-info", zap.Error(err))
		return ""
	}

	return string(output)
}

func (nix *Nix) compilePackage(registry, pkg string) string {
	return fmt.Sprintf("%s#%s", registry, pkg)
}
