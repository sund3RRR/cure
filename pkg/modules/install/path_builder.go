package install

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/sund3RRR/cure/pkg/types"
)

type PathBuilder struct {
	ProfilePath string
}

func NewPathBuilder(profilePath string) *PathBuilder {
	return &PathBuilder{
		ProfilePath: profilePath,
	}
}

func (b *PathBuilder) Build(packagePath types.Path, prefix types.Path, files []types.File) error {
	params := types.WriteParams{Prefix: prefix, Mode: 0755}
	errhandler := func(fileType types.FileType, file types.File, err error) error {
		switch fileType {
		case types.RegularFile, types.SymlinkFile:
			return nil
		case types.DirectoryFile:
			if errors.Is(err, os.ErrExist) {
				link, err := os.Readlink(file.GetPath().String())
				if err != nil {
					return nil
				}

				if err := os.Remove(file.GetPath().String()); err != nil {
					return err
				}
				if err := os.Mkdir(file.GetPath().String(), 0755); err != nil {
					return err
				}
				dir, err := os.ReadDir(link)
				if err != nil {
					return err
				}

				for _, f := range dir {
					if err := os.Symlink(filepath.Join(link, f.Name()), filepath.Join(file.GetPath().String(), f.Name())); err != nil {
						return err
					}
				}
			}
			return nil
		}
		return nil
	}

	for _, file := range files {
		_, err := file.Write(params, errhandler)
		if err != nil {
			return err
		}
	}

	return b.symlinkAll(packagePath, prefix)
}

func (b *PathBuilder) symlinkAll(old, prefix types.Path) error {
	dir, err := os.ReadDir(old.String())
	if err != nil {
		return err
	}

	for _, file := range dir {
		err := os.Symlink(filepath.Join(old.String(), file.Name()), filepath.Join(prefix.String(), file.Name()))
		if err != nil {
			if file.IsDir() {
				return b.symlinkAll(types.NewPath(old.String(), file.Name()), prefix+types.Path(file.Name()))
			}
		}
	}

	return nil
}
