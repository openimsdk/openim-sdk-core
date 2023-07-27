package testv3new

import (
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"open_im_sdk/testv3new/testcore"
	"testing"
	"time"
)

var (
	pressureTestAttribute PressureTestAttribute
)

func init() {
	pressureTestAttribute.InitWithFlag()
	if err := log.InitFromConfig("sdk.log", "sdk", 3,
		true, false, "", 2); err != nil {
		panic(err)
	}
}

func TestPressureTester_PressureSendMsgs(t *testing.T) {
	ParseFlag()
	p := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	for i := 0; i < 10; i++ {
		p.WithTimer(p.PressureSendMsgs2)(pressureTestAttribute.sendUserIDs, pressureTestAttribute.recvUserIDs, pressureTestAttribute.messageNumber, time.Duration(pressureTestAttribute.timeInterval)*time.Millisecond)
		time.Sleep(time.Second)
	}
	// time.Sleep(1000 * time.Second)
}

func TestPressureTester_PressureSendGroupMsgs(t *testing.T) {
	ParseFlag()
	pressureTester := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	pressureSendGroupMsgsWithTime := pressureTester.WithTimer(pressureTester.PressureSendGroupMsgs)
	pressureSendGroupMsgsWithTime(pressureTestAttribute.sendUserIDs, pressureTestAttribute.groupIDs, pressureTestAttribute.groupIDs, time.Duration(pressureTestAttribute.timeInterval)*time.Millisecond)
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
