package testv3new

import (
	"open_im_sdk/testv3new/testcore"
	"testing"
	"time"
)

func init() {

}

func TestPressureTester_PressureSendMsgs(t *testing.T) {
	sendUserID := "bantanger"
	recvUserID := []string{"9927048690"}

	pressureTester := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	pressureTester.PressureSendMsgs(sendUserID, recvUserID, 3, time.Duration(1))
}
