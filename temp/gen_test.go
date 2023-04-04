package temp

import (
	"fmt"
	"os"
	"regexp"
	"testing"
)

func TestName(t *testing.T) {
	data, err := os.ReadFile("./../open_im_sdk/group.go")
	if err != nil {
		panic(err)
	}
	reg := regexp.MustCompile(`func \w*\(callback open_im_sdk_callback.Base, operationID string.*\) {`)
	//reg.MatchString(string(data))

	fmt.Println(reg.FindAllString(string(data), -1))

}
