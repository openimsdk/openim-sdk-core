package testv3new

import (
	"fmt"
	"github.com/OpenIMSDK/tools/log"
	"open_im_sdk/testv3new/testcore"
	"testing"
	"time"
)

func init() {
	if err := log.InitFromConfig("sdk.log", "sdk", 3,
		true, false, "", 2, 24); err != nil {
		panic(err)
	}
}

func TestPressureTester_PressureSendMsgs(t *testing.T) {
	sendUserID := []string{"register_test_493"}
	recvUserID := []string{"5338610321"}

	p := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	for i := 0; i < 10; i++ {
		p.WithTimer(p.PressureSendMsgs2)(sendUserID, recvUserID, 1000, 100*time.Millisecond)
		time.Sleep(time.Second)
	}
	// time.Sleep(1000 * time.Second)
}

func TestPressureTester_PressureSendGroupMsgs(t *testing.T) {
	sendUserID := "register_test_4334"
	groupID := "3411007805"

	pressureTester := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	pressureSendGroupMsgsWithTime := pressureTester.WithTimer(pressureTester.PressureSendGroupMsgs)
	pressureSendGroupMsgsWithTime([]string{sendUserID}, groupID, 100, time.Duration(100))
}

func TestPressureTester_PressureSendGroupMsgs2(t *testing.T) {
	start := 850
	count := 900
	step := 10
	for j := start; j <= count; j += step {
		var sendUserIDs []string
		startTime := time.Now().UnixNano()
		for i := j; i < j+step; i++ {
			sendUserIDs = append(sendUserIDs, fmt.Sprintf("register_test_%v", i))
		}
		// groupID := "3411007805"
		// groupID := "2347514573"
		groupID := "3167736657"
		pressureTester := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
		pressureTester.PressureSendGroupMsgs(sendUserIDs, groupID, 1, 0)
		endTime := time.Now().UnixNano()
		nanoSeconds := float64(endTime - startTime)
		t.Log("", nanoSeconds)
		fmt.Println()
	}
}

func TestPressureTester_Conversation(t *testing.T) {
	sendUserID := "5338610321"
	var recvUserIDs []string
	for i := 1; i <= 1000; i++ {
		recvUserIDs = append(recvUserIDs, fmt.Sprintf("register_test_%v", i))
	}
	p := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	p.WithTimer(p.PressureSendMsgs)(sendUserID, recvUserIDs, 1, 100*time.Millisecond)
}

func TestPressureTester_PressureSendMsgs2(t *testing.T) {
	recvUserID := "5338610321"
	var sendUserIDs []string
	for i := 1; i <= 1000; i++ {
		sendUserIDs = append(sendUserIDs, fmt.Sprintf("register_test_%v", i))
	}
	p := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	p.WithTimer(p.PressureSendMsgs2)(sendUserIDs, []string{recvUserID}, 1, 100*time.Millisecond)
}
