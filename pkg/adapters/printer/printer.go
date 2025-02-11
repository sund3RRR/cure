// Package printer implements printing methods
package printer

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/text"
)

// Alignment is the alignment of the output
const Alignment = 5

// Success prints '✅ Success'
func Success(fd *os.File) {
	success := text.FgGreen.Sprint("✅ Success\n")
	str := text.AlignCenter.Apply(success, Alignment)

	_, _ = fmt.Fprint(fd, str)
}

// Successf prints '✅ Successfully' with formatted string and args
func Successf(fd *os.File, format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	success := text.FgGreen.Sprintf("✅ Successfully %s\n", formatted)
	aligned := text.AlignCenter.Apply(success, Alignment)

	_, _ = fmt.Fprint(fd, aligned)
}

// Failed prints '❌ Failed'
func Failed(fd *os.File) {
	failed := text.FgRed.Sprint("❌ Failed\n")
	str := text.AlignCenter.Apply(failed, Alignment)

	_, _ = fmt.Fprint(fd, str)
}

// Failedf prints '❌ Failed' with formatted string and args
func Failedf(fd *os.File, format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	failed := text.FgRed.Sprintf("❌ Failed %s\n", formatted)
	aligned := text.AlignCenter.Apply(failed, Alignment)

	_, _ = fmt.Fprint(fd, aligned)
}
