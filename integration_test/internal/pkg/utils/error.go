package utils

import (
	"fmt"
	"strings"
)

func FormatErrorStack(err error) string {
	if err == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%+v", err))

	stack := sb.String()
	stack = strings.ReplaceAll(stack, "\n", " => ")
	stack = strings.ReplaceAll(stack, "\t", " ")

	return stack
}
