package install

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sund3RRR/cure/pkg/adapters/nix"
	"github.com/sund3RRR/cure/pkg/types"
	"github.com/sund3RRR/cure/pkg/types/file"
)

type NixGLWrapper struct {
	nix      Nix
	gpu      Gpu
	enable   types.NixGL
	pkg      types.NixGLPackage
	execPath string
}

func NewNixGLWrapper(nix Nix, gpu Gpu) *NixGLWrapper {
	return &NixGLWrapper{
		nix:    nix,
		gpu:    gpu,
		enable: types.NixGLAuto,
		pkg:    types.NixGLPackageAuto,
	}
}

func (n *NixGLWrapper) SetParams(params Params) {
	n.enable = params.NixGL
	n.pkg = params.NixGLPackage
}

func (n *NixGLWrapper) CheckAndPrepare() error {
	// Identify nixGL package will determine what gpu is using
	// and which nixGL package to install
	n.pkg = n.identifyPackage(n.pkg)

	// Convert enum representation of package to real one
	pkg := convertToNixGLPackage(n.pkg)

	// Deal with 'nixGL' registry. If there is no such, add the 'nixGL' registry
	pi, err := n.nix.PathInfo("nixgl", pkg)
	if err != nil {
		if errors.Is(err, nix.ErrCannotFindFlake) {
			err := n.nix.AddRegistry("nixgl", "github:nix-community/nixGL")
			if err != nil {
				return err
			}
		}
		return err
	}

	// Check that selected package exists in the store.
	// If it doesn't, download package from 'nixGL' registry
	if !pi.Out.Exists() {
		pi, err = n.nix.GetPackage("nixgl", pkg)
		if err != nil {
			return err
		}
	}

	// Set exec path with selected package
	n.execPath = filepath.Join(pi.Out.String(), "bin", pkg)

	return nil
}

func (n *NixGLWrapper) Apply(packagePath string, files []*file.File) []*file.File {
	binDir := filepath.Join(packagePath, "bin")
	dirFiles, err := os.ReadDir(binDir)
	if err != nil {
		return files
	}

	for _, f := range dirFiles {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		content := strings.Join(
			[]string{
				"#!/usr/bin/env bash",
				fmt.Sprintf("%s/bin/nixGL %s \"$@\"", n.execPath, filepath.Join(packagePath, "bin", f.Name())),
			},
			"\n",
		)
		wrappedBin := file.NewFile(filepath.Join("bin", f.Name()), types.RootOwner, 0755, []byte(content))
		files = append(files, wrappedBin)
	}

	return files
}

func (n *NixGLWrapper) identifyPackage(pkg types.NixGLPackage) types.NixGLPackage {
	if n.pkg == types.NixGLPackageAuto {
		gpu, _ := n.gpu.GetManufacturer()
		if gpu == types.NVIDIA {
			return types.NixGLNvidia
		}
		return types.NixGLMesa
	}

	return pkg
}

func convertToNixGLPackage(pkg types.NixGLPackage) string {
	result := ""
	switch pkg {
	case types.NixGLPackageAuto:
		result = "nixGLIntel"
	case types.NixGLMesa:
		result = "nixGLIntel"
	case types.NixVulkanMesa:
		result = "nixVulkanIntel"
	case types.NixGLNvidia:
		result = "nixGLNvidia"
	case types.NixVulkanNvidia:
		result = "nixVulkanNvidia"
	}
	return result
}

package install

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sund3RRR/cure/pkg/adapters/nix"
	"github.com/sund3RRR/cure/pkg/types"
)

type NixGLWrapper struct {
	nix      Nix
	gpu      Gpu
	enable   types.NixGL
	pkg      types.NixGLPackage
	execPath string
}

func NewNixGLWrapper(nix Nix, gpu Gpu) *NixGLWrapper {
	return &NixGLWrapper{
		nix:    nix,
		gpu:    gpu,
		enable: types.NixGLAuto,
		pkg:    types.NixGLPackageAuto,
	}
}

func (n *NixGLWrapper) SetParams(params Params) {
	n.enable = params.NixGL
	n.pkg = params.NixGLPackage
}

func (n *NixGLWrapper) CheckAndPrepare() error {
	// Identify nixGL package will determine what gpu is using
	// and which nixGL package to install
	n.pkg = n.identifyPackage(n.pkg)

	// Convert enum representation of package to real one
	pkg := convertToNixGLPackage(n.pkg)

	// Deal with 'nixGL' registry. If there is no such, add the 'nixGL' registry
	pi, err := n.nix.PathInfo("nixgl", pkg)
	if err != nil {
		if errors.Is(err, nix.ErrCannotFindFlake) {
			err := n.nix.AddRegistry("nixgl", "github:nix-community/nixGL")
			if err != nil {
				return err
			}
		}
		return err
	}

	// Check that selected package exists in the store.
	// If it doesn't, download package from 'nixGL' registry
	if !pi.Out.Exists() {
		pi, err = n.nix.GetPackage("nixgl", pkg)
		if err != nil {
			return err
		}
	}

	// Set exec path with selected package
	n.execPath = filepath.Join(pi.Out.String(), "bin", pkg)

	return nil
}

func (n *NixGLWrapper) Apply(packagePath string, files []types.File) []types.File {
	binDir := filepath.Join(packagePath, "bin")
	dirFiles, err := os.ReadDir(binDir)
	if err != nil {
		return files
	}

	binFiles := make([]types.File, 0)
	for _, f := range dirFiles {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		content := strings.Join(
			[]string{
				"#!/usr/bin/env bash",
				fmt.Sprintf("%s %s \"$@\"", n.execPath, filepath.Join(packagePath, "bin", f.Name())),
			},
			"\n",
		)
		wrappedBin := types.NewRegularFile(types.NewPath("bin", f.Name()), []byte(content))
		binFiles = append(binFiles, wrappedBin)
	}

	dir := types.NewDirectory(types.NewPath("bin"), binFiles...)
	return append(files, dir)
}

func (n *NixGLWrapper) identifyPackage(pkg types.NixGLPackage) types.NixGLPackage {
	if n.pkg == types.NixGLPackageAuto {
		gpu, _ := n.gpu.GetManufacturer()
		if gpu == types.NVIDIA {
			return types.NixGLNvidia
		}
		return types.NixGLMesa
	}

	return pkg
}

func convertToNixGLPackage(pkg types.NixGLPackage) string {
	result := ""
	switch pkg {
	case types.NixGLPackageAuto:
		result = "nixGLIntel"
	case types.NixGLMesa:
		result = "nixGLIntel"
	case types.NixVulkanMesa:
		result = "nixVulkanIntel"
	case types.NixGLNvidia:
		result = "nixGLNvidia"
	case types.NixVulkanNvidia:
		result = "nixVulkanNvidia"
	}
	return result
}
