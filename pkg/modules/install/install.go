package install

import (
	"strings"

	"github.com/sund3RRR/cure/pkg/adapters/gpu"
	"github.com/sund3RRR/cure/pkg/types"
)

type Nix interface {
	GetPackage(registry, pkg string) (types.PackageInfo, error)
	PathInfo(registry, pkg string) (types.PackageInfo, error)
	AddRegistry(alias, registry string) error
}

type Gpu interface {
	GetManufacturer() (types.GPUManufacturer, error)
}

type File interface {
	Write() (string, error)
	GetPath() string
	IsDir() bool
}

type InstallerModule interface {
	CheckAndPrepare(pkgPath types.Path, params Params) error
	Apply(pkgPath types.Path, files []types.File) []types.File
}

type Params struct {
	NixGL        types.NixGL
	NixGLPackage types.NixGLPackage
}

type Installer struct {
	nix     Nix
	modules []InstallerModule
	builder *PathBuilder
}

func NewInstaller(nix Nix) *Installer {
	return &Installer{
		nix: nix,
		modules: []InstallerModule{
			NewNixGLWrapper(nix, gpu.NewGPU()),
		},
		builder: NewPathBuilder("/home/sunder/dev/cure/profile"),
	}
}

func (i *Installer) InstallPackage(name string, params Params) error {
	// Substitute empty registry with 'nixpkgs'
	splitted := strings.Split(name, "#")
	var registry, pkg string
	if len(splitted) == 1 {
		registry, pkg = "nixpkgs", splitted[0]
	} else {
		registry, pkg = splitted[0], splitted[1]
	}

	// Download package to /nix/store
	pi, err := i.nix.GetPackage(registry, pkg)
	if err != nil {
		return err
	}

	// Prepare modules for building profile
	for _, m := range i.modules {
		if err := m.CheckAndPrepare(pi.Out, params); err != nil {
			return err
		}
	}

	// Apply modules, e.g. modifications to packages
	files := make([]types.File, 0)
	for _, m := range i.modules {
		files = m.Apply(pi.Out, files)
	}

	// Build profile path
	return i.builder.Build(pi.Out, types.NewPath("/opt/cure"), files)
}
