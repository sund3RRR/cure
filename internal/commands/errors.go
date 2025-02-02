package commands

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrPackageNotProvided error appears when package was not provided by the user
	ErrPackageNotProvided = errors.New("you must provide package name in format <package> " +
		"for default registry (nixpkgs) or <registry#package> for another package registry")
)

// UnrecognizedFlagValueError is an error type for unrecognized flag value
type UnrecognizedFlagValueError[S ~[]E, E comparable] struct {
	FlagValue   E
	ValidValues S
}

// Error returns string representation of the error
func (e *UnrecognizedFlagValueError[S, E]) Error() string {
	validValuesStr := make([]string, len(e.ValidValues))
	for i, v := range e.ValidValues {
		validValuesStr[i] = fmt.Sprintf("%v", v)
	}
	return fmt.Sprintf("Unrecognized flag value: %v. Valid values are: %s", e.FlagValue, strings.Join(validValuesStr, ", "))
}
