package nix

import (
	"fmt"
	"regexp"
)

func (nix *Nix) compilePackage(registry, pkg string) string {
	return fmt.Sprintf("%s#%s", registry, pkg)
}

type ErrorPattern struct {
	pattern string
	err     error
}

// reflectNixError returns error based on nix output
func reflectNixError(nixOutput string) error {
	patterns := []ErrorPattern{
		{
			pattern: `error: flake '([^']+)' does not provide attribute '([^']+)'`,
			err:     ErrPackageNotFound,
		},
		{
			pattern: `error: unable to download '([^']+)'(?:.*HTTP error 404`,
			err:     ErrCannotDownloadFlake,
		},
		{
			pattern: `error: cannot find flake '([^']+)' in the flake registries`,
			err:     ErrCannotFindFlake,
		},
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern.pattern)
		matches := re.FindStringSubmatch(nixOutput)

		if len(matches) > 0 {
			return pattern.err
		}
	}

	return ErrUnknownNixError
}
