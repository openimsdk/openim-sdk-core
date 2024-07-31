package utils

import (
	"errors"
	"fmt"
)

// PrintErrorStack print err stack
func PrintErrorStack(err error) string {
	var stack []string

	for err != nil {
		stack = append(stack, err.Error())
		err = errors.Unwrap(err)
	}

	return fmt.Sprintf("Error stack:\n%s", formatStack(stack))
}

// formatStack
func formatStack(stack []string) string {
	var result string
	for i := len(stack) - 1; i >= 0; i-- {
		result += fmt.Sprintf("%d: %s\n", len(stack)-i, stack[i])
	}
	return result
}
