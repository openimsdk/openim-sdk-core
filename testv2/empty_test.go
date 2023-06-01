package testv2

import (
	"testing"
	"time"
)

func Test_Empty(t *testing.T) {
	for {
		time.Sleep(time.Second * 10)
	}
}

func Test_RunWait(t *testing.T) {
	time.Sleep(time.Second * 10)
}
