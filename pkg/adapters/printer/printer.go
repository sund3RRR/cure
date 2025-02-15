// Package printer implements printing methods
package printer

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/text"
)

// Alignment is the alignment of the output
const Alignment = 5

// Success prints '✅' with formatted string and args
func Success(fd *os.File, format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	success := text.FgGreen.Sprintf("✅ %s\n", formatted)
	aligned := text.AlignCenter.Apply(success, Alignment)

	_, _ = fmt.Fprint(fd, aligned)
}

// Failed prints '❌' with formatted string and args
func Failed(fd *os.File, format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	failed := text.FgRed.Sprintf("❌ %s\n", formatted)
	aligned := text.AlignCenter.Apply(failed, Alignment)

	_, _ = fmt.Fprint(fd, aligned)
}

func Skip(fd *os.File, format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	skipped := text.FgHiGreen.Sprintf("🌿 %s\n", formatted)
	aligned := text.AlignCenter.Apply(skipped, Alignment)

	_, _ = fmt.Fprint(fd, aligned)
}

func Processing(fd *os.File, format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	processing := fmt.Sprintf("🌀 %s\n", formatted)
	aligned := text.AlignCenter.Apply(processing, Alignment)

	_, _ = fmt.Fprint(fd, aligned)
}

func Info(fd *os.File, format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	info := fmt.Sprintf("ℹ️ %s\n", formatted)
	aligned := text.AlignCenter.Apply(info, Alignment)

	_, _ = fmt.Fprint(fd, aligned)
}
