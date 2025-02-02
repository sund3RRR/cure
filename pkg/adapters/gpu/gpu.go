package gpu

import (
	"os/exec"
	"regexp"

	"github.com/sund3RRR/cure/pkg/types"
)

type Gpu struct {
}

func NewGPU() *Gpu {
	return &Gpu{}
}

func (g *Gpu) GetManufacturer() (types.GPUManufacturer, error) {
	cmd := exec.Command("lspci", "|", "grep", "-i", "vga")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	vendorPatterns := map[types.GPUManufacturer]*regexp.Regexp{
		types.Intel:  regexp.MustCompile(`.*Intel.*VGA.*`),
		types.AMD:    regexp.MustCompile(`.*Advanced.*AMD.*VGA.*`),
		types.NVIDIA: regexp.MustCompile(`.*NVIDIA.*VGA.*`),
	}

	for vendor, pattern := range vendorPatterns {
		if pattern.MatchString(string(output)) {
			return vendor, nil
		}
	}

	return "", ErrUnknownGPUManufacturer
}
