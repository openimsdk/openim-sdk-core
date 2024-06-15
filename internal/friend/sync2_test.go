package friend

import (
	"fmt"
	"testing"
)

func Test_main(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	fmt.Println(a[:3])
	fmt.Println(a[3:])
	fmt.Println(a[2:4])
}
