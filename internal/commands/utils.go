package commands

import (
	"fmt"
	"slices"
	"strings"

	"github.com/sund3RRR/cure/pkg/types"
)

func validateFlagValue[S ~[]E, E comparable](flagValue E, validValues S) error {
	if !slices.Contains(validValues, flagValue) {
		return &UnrecognizedFlagValueError[S, E]{
			FlagValue:   flagValue,
			ValidValues: validValues,
		}
	}
	return nil
}

func possibleValuesString[S ~[]E, E comparable](values S, defaultValue E, highlighter types.Sequence) string {
	var formattedValues []string
	for _, value := range values {
		strValue := fmt.Sprintf("%v", value)
		if value == defaultValue {
			formattedValues = append(formattedValues, applyFormat(highlighter, strValue))
		} else {
			formattedValues = append(formattedValues, strValue)
		}
	}

	return fmt.Sprintf("Possible values: [%s]", strings.Join(formattedValues, ", "))
}

func applyFormat(format types.Sequence, s string) string {
	return string(format) + s + string(types.ResetSeq)
}
