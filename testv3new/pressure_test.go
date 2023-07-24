package testv3new

import (
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"open_im_sdk/testv3new/testcore"
	"testing"
	"time"
)

func init() {
	if err := log.InitFromConfig("sdk.log", "sdk", 3,
		true, false, "", 2); err != nil {
		panic(err)
	}
}

func TestPressureTester_PressureSendMsgs(t *testing.T) {
	sendUserID := "bantanger"
	recvUserID := []string{"9927048690"}

	pressureTester := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	pressureTester.PressureSendMsgs(sendUserID, recvUserID, 3, 100)
}

func TestPressureTester_PressureSendGroupMsgs(t *testing.T) {
	sendUserID := "register_test_1"
	groupID := "3411007805"

	pressureTester := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	pressureTester.PressureSendGroupMsgs([]string{sendUserID}, groupID, 100, 100)
}

func TestPressureTester_PressureSendGroupMsgs2(t *testing.T) {
	start := 1
	count := 10000
	step := 100
	for j := start; j <= count; j += step {
		var sendUserIDs []string
		startTime := time.Now().UnixNano()
		for i := j; i <= j+step; i++ {
			sendUserIDs = append(sendUserIDs, fmt.Sprintf("register_test_%v", i))
		}
		// groupID := "3411007805"
		// groupID := "2347514573"
		groupID := "3813706739"
		pressureTester := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
		pressureTester.PressureSendGroupMsgs(sendUserIDs, groupID, 1, 0)
		endTime := time.Now().UnixNano()
		nanoSeconds := float64(endTime - startTime)
		t.Log("", nanoSeconds)
		fmt.Println()
	}
}
