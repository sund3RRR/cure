package nix

import "errors"

var (
	// ErrInvalidBuildData error appears when build data is invalid
	ErrInvalidBuildData = errors.New("invalid build data")
	// ErrPackageNotFound error appears when package is not found
	ErrPackageNotFound = errors.New("package not found in the specified registry")
	// ErrUnknownNixError error appears when nix returns unknown error
	ErrUnknownNixError = errors.New("unknown nix error")
	// ErrCannotDownloadFlake error appears when flake repository is not found
	ErrCannotDownloadFlake = errors.New("flake repository not found")
	// ErrCannotFindFlake error appears when flake is not found in the registry
	ErrCannotFindFlake = errors.New("cannot find flake in the registry")
)
