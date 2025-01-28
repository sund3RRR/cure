package desktop

import "errors"

var (
	// ErrDesktopDirNotFound error appears when share/applications directory not found
	ErrDesktopDirNotFound = errors.New("share/applications directory not found")
)
