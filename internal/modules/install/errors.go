package install

import "errors"

// ErrPackageNotProvided error appears when package was not provided by the user
var ErrPackageNotProvided = errors.New("you must provide package name in format <package> " +
	"for default registry (nixpkgs) or <registry#package> for another package registry")
