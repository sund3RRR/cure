package types

import "fmt"

// PackageInfo contains information about package
type PackageInfo struct {
	Name    string
	Pname   string
	Version string
	Out     Path            // main package output
	Outputs map[string]Path // other outputs WITHOUT 'out'
	System  string
}

func (pi PackageInfo) String() string {
	return fmt.Sprintf(
		"Name: %s\nPname: %s\nVersion: %s\nOut: %s\nOutputs: %v\nSystem: %s\n",
		pi.Name, pi.Pname, pi.Version, pi.Out, pi.Outputs, pi.System,
	)
}
