package types

import "os"

type Path string

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
