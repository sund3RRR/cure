package install

import (
	"os"
	"strings"

	"github.com/sund3RRR/cure/pkg/adapters/gpu"
	"github.com/sund3RRR/cure/pkg/adapters/printer"
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
	GetName() string
	CheckAndPrepare(pkgPath types.Path, params Params) (bool, error)
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

func (installer *Installer) InstallPackage(name string, params Params) error {
	// Substitute empty registry with 'nixpkgs'
	splitted := strings.Split(name, "#")
	var registry, pkg string
	if len(splitted) == 1 {
		registry, pkg = "nixpkgs", splitted[0]
	} else {
		registry, pkg = splitted[0], splitted[1]
	}

	// Download package to /nix/store
	printer.Processing(os.Stdout, "Downloading %s#%s...", registry, pkg)
	pi, err := installer.nix.GetPackage(registry, pkg)
	if err != nil {
		return err
	}
	printer.Success(os.Stdout, "Successfully downloaded %s#%s", registry, pkg)

	// Prepare modules for building profile
	modulesBoolMap := make([]bool, len(installer.modules))
	for i, m := range installer.modules {
		enable, err := m.CheckAndPrepare(pi.Out, params)
		if err != nil {
			return err
		}
		modulesBoolMap[i] = enable

		if enable {
			printer.Success(os.Stdout, "using %s module", m.GetName())
		} else {
			printer.Skip(os.Stdout, "skipping %s module", m.GetName())
		}
	}

	// Apply modules, e.g. modifications to packages
	files := make([]types.File, 0)
	for i, m := range installer.modules {
		if modulesBoolMap[i] {
			files = m.Apply(pi.Out, files)
		}
	}

	// Build profile path
	return installer.builder.Build(pi.Out, types.NewPath("/opt/cure"), files)
}
