package install

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/sund3RRR/cure/pkg/types/file"
)

type PathBuilder struct {
	ProfilePath string
	ExcludeDirs []string
}

func NewPathBuilder(profilePath string) *PathBuilder {
	return &PathBuilder{
		ProfilePath: profilePath,
		ExcludeDirs: []string{
			"bin",
			"lib",
			"lib/system/systemd",
			"lib64",
			"libexec",
			"include",
			"etc",
			"share",
			"share/applications",
			"doc",
			"man",
		},
	}
}

func (s *PathBuilder) Build(packagePath string, files []*file.File) error {
	for _, file := range files {
		if file.IsDir() {
			if slices.Contains(s.ExcludeDirs, file.Name) {
				os.Mkdir(file.GetPath(), os.ModePerm)
			} else {
				err := os.Symlink(packagePath+file.GetPath(), s.ProfilePath+file.GetPath())
				if err != nil {
					return err
				}
			}
			continue
		}

		err := os.WriteFile(filepath.Join(s.ProfilePath, file.GetPath()), file.Content, file.Mode)
		if err != nil {
			return err
		}
	}

	err := filepath.WalkDir(packagePath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(packagePath, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(s.ProfilePath, relPath)

		if d.IsDir() {
			if slices.Contains(s.ExcludeDirs, relPath) {
				if _, err := os.Lstat(dstPath); os.IsNotExist(err) {
					if err := os.Mkdir(dstPath, 0755); err != nil {
						return fmt.Errorf("ошибка при создании директории %s: %v", path, err)
					}
				}
			} else {
				if _, err := os.Lstat(dstPath); os.IsNotExist(err) {
					if err := os.Symlink(path, dstPath); err != nil {
						return fmt.Errorf("ошибка при создании симлинка для директории %s: %v", path, err)
					}
				}
			}
		} else {
			if _, err := os.Lstat(dstPath); os.IsNotExist(err) {
				err := os.Symlink(path, dstPath)
				if err != nil {
					return fmt.Errorf("ошибка при создании симлинка для файла %s: %v", path, err)
				}
			}
		}

		return nil
	})

	return err
}
