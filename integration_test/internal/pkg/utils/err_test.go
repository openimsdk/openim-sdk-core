package utils

import (
	"fmt"
	"github.com/openimsdk/tools/errs"
	"testing"
)

func TestErr(t *testing.T) {
	err := A5()
	fmt.Println(FormatErrorStack(err))
}

func A1() error {
	err := errs.New("err1").Wrap()
	return err
}

func A2() error {
	return A1()
}

func A3() error {
	return A2()
}

func A4() error {
	return A3()
}

func A5() error {
	return A4()
}
