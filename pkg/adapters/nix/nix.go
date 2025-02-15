// Package nix is an adapter to nix command
package nix

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"

	"github.com/sund3RRR/cure/pkg/types"
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
		command: "/nix/store/ar6980is6p3993lgpr07wk2i0ynsm3lf-nix-2.27.0pre19700101_dirty/bin/nix",
		flags:   flags,
	}
}

// GetPackage takes provided 'registry' and 'pkg', downloads it to /nix/store and
// returns package path
func (nix *Nix) GetPackage(registry, pkg string) (types.PackageInfo, error) {
	args := append(
		nix.flags,
		"build",
		"--no-link",
		nix.compilePackage(registry, pkg),
	)

	cmd := exec.Command(nix.command, args...) //nolint

	var outputBuffer bytes.Buffer
	multiWriter := io.MultiWriter(os.Stdout, &outputBuffer)
	cmd.Stdout = multiWriter
	cmd.Stderr = multiWriter

	err := cmd.Run()
	cmdOutput := outputBuffer.String()
	if err != nil {
		err := reflectNixError(cmdOutput)
		if err == ErrUnknownNixError {
			nix.logger.Error(
				"failed to execute nix build",
				zap.Error(err), zap.Strings("args", args), zap.String("cmdOutput", cmdOutput),
			)
		}
		return types.PackageInfo{}, err
	}

	return nix.PathInfo(registry, pkg)
}

// PathInfo takes provided 'registry' and 'pkg' and returns package path
// in the /nix/store. Returns empty string if there is no such package.
func (nix *Nix) PathInfo(registry, pkg string) (types.PackageInfo, error) {
	type RawPackageInfo struct {
		Name string `json:"name"`
		Env  struct {
			Pname   string `json:"pname"`
			Version string `json:"version"`
			Out     string `json:"out"`
			Outputs string `json:"outputs"`
		} `json:"env"`
		Outputs map[string]struct {
			Path string `json:"path"`
		} `json:"outputs"`
		System string `json:"system"`
	}

	args := append(
		nix.flags,
		"derivation",
		"show",
		nix.compilePackage(registry, pkg),
		"--impure",
	)

	cmd := exec.Command(nix.command, args...) //nolint
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		err := reflectNixError(string(cmdOutput))
		if err == ErrUnknownNixError {
			nix.logger.Error(
				"failed to execute nix derivation show",
				zap.Error(err), zap.Strings("args", args), zap.ByteString("cmdOutput", cmdOutput),
			)
		}
		return types.PackageInfo{}, err
	}

	var data map[string]RawPackageInfo
	err = json.Unmarshal(cmdOutput, &data)
	if err != nil {
		nix.logger.Error("failed to unmarshal nix derivation show output", zap.Error(err))
		return types.PackageInfo{}, err
	}

	if len(data) != 1 {
		return types.PackageInfo{}, ErrInvalidBuildData
	}

	var rawPackageInfo RawPackageInfo
	for _, packageValue := range data {
		rawPackageInfo = packageValue
	}

	packageInfo := types.PackageInfo{
		Name:    rawPackageInfo.Name,
		Pname:   rawPackageInfo.Env.Pname,
		Version: rawPackageInfo.Env.Version,
		Out:     types.Path(rawPackageInfo.Env.Out),
		Outputs: make(map[string]types.Path),
		System:  rawPackageInfo.System,
	}

	for outputKey, output := range rawPackageInfo.Outputs {
		if outputKey == "out" {
			continue
		}
		packageInfo.Outputs[outputKey] = types.Path(output.Path)
	}

	return packageInfo, nil
}
