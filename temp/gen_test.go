package temp

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	type Temp struct {
		Name string
		Type string
	}
	argsTypeName := func(s string) []Temp {
		var ts []Temp
		for _, v := range strings.Split(s, ",") {
			args := strings.Split(strings.ReplaceAll(strings.TrimSpace(v), "  ", " "), " ")
			if len(args) == 1 {
				ts = append(ts, Temp{Name: args[0]})
			} else {
				ts = append(ts, Temp{Name: args[0], Type: args[1]})
			}
		}
		lastType := ts[len(ts)-1].Type
		for i := len(ts) - 2; i >= 0; i-- {
			if ts[i].Type == "" {
				ts[i].Type = lastType
			} else {
				lastType = ts[i].Type
			}
		}
		return ts
	}
	formatArgs := func(ts []Temp) string {
		var arr []string
		for _, temp := range ts {
			arr = append(arr, fmt.Sprintf("%s %s", temp.Name, temp.Type))
		}
		return strings.Join(arr, ", ")
	}
	argsNames := func(ts []Temp) []string {
		var arr []string
		for _, temp := range ts {
			arr = append(arr, temp.Name)
		}
		return arr
	}

	filemap := map[string]string{
		"conversation_msg.go": "Conversation",
		"friend.go":           "Friend",
		"group.go":            "Group",
		"user.go":             "User",
		"signaling.go":        "Signaling",
		"workmoments.go":      "WorkMoments",
	}

	for name, val := range filemap {
		data, err := os.ReadFile(filepath.Join("./../open_im_sdk", name))
		if err != nil {
			panic(err)
		}
		reg := regexp.MustCompile(`func \w*\(callback open_im_sdk_callback.Base, operationID string.*\) {`)
		var arr []string
		arr = append(arr, "package open_im_sdk\n\n")
		for _, line := range reg.FindAllString(string(data), -1) {
			funcname := line[len("func "):strings.Index(line, "(")]
			args := argsTypeName(line[strings.Index(line, "(")+1 : strings.LastIndex(line, ")")])
			arr = append(arr, fmt.Sprintf("func %s(%s) {", funcname, formatArgs(args)))
			if len(args) > 2 {
				arr = append(arr, fmt.Sprintf("\tcall(%s, %s, userForSDK.%s().%s, %s)", args[0].Name, args[1].Name, val, funcname, strings.Join(argsNames(args[2:]), ",")))
			} else {
				arr = append(arr, fmt.Sprintf("\tcall(%s, %s, userForSDK.%s().%s)", args[0].Name, args[1].Name, val, funcname))
			}
			arr = append(arr, "}")
			arr = append(arr, "\n")
		}
		if err := os.WriteFile(name, []byte(strings.Join(arr, "\n")), 0666); err != nil {
			panic(err)
		}
	}
}
