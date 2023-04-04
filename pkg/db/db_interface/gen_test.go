package db_interface

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	data, err := os.ReadFile("databse.go")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n")
	var arr []string
	for _, line := range lines {
		if strings.Index(line, "(") > 0 && strings.Index(line, ")") > 0 {
			line = strings.Replace(line, "(", "(ctx context.Context, ", 1)
		}
		arr = append(arr, line)
	}
	fmt.Println(strings.Join(arr, "\n"))

}
