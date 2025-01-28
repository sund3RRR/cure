// Package desktop implements a module, that manages desktop files
package desktop

import (
	"fmt"
	"os"
	"strings"
)

// Desktop is the module, that manages desktop files
type Desktop struct {
}

// NewDesktop returns a new desktop module
func NewDesktop() *Desktop {
	return &Desktop{}
}

// GetDesktopItems returns desktop items from provided package
func (d *Desktop) GetDesktopItems(pkg string) ([]Item, error) {
	dir, err := os.ReadDir(strings.TrimSuffix(pkg, "/") + "/share/applications")
	if err != nil {
		return nil, ErrDesktopDirNotFound
	}

	desktopItems := make([]Item, 0, len(dir))
	for _, file := range dir {
		fileName := file.Name()

		if !file.IsDir() && strings.Contains(fileName, ".desktop") {
			item, err := NewItemFromPath(d.buildPath(pkg, fileName))
			if err != nil {
				return nil, err
			}
			desktopItems = append(desktopItems, item)
		}
	}

	return desktopItems, nil
}

func (d *Desktop) buildPath(pkg, file string) string {
	return fmt.Sprintf("%s/share/applications/%s", pkg, file)
}
