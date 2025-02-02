package types

import (
	"os"
	"path/filepath"
)

type Path string

func NewPath(e ...string) Path {
	return Path(filepath.Join(e...))
}
func (p Path) String() string {
	return string(p)
}

func (p Path) Exists() bool {
	_, err := os.Stat(p.String())
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}
