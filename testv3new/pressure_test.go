package testv3new

import (
	"fmt"
	"github.com/OpenIMSDK/tools/log"
	"open_im_sdk/testv3new/testcore"
	"testing"
	"time"
)

var (
	pressureTestAttribute PressureTestAttribute
)

func init() {
	pressureTestAttribute.InitWithFlag()
	if err := log.InitFromConfig("sdk.log", "sdk", 4,
		true, true, "./chat_log", 2, 24); err != nil {
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
	p := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	for i := 0; i < 10; i++ {
		p.WithTimer(p.PressureSendGroupMsgs2)(pressureTestAttribute.sendUserIDs, pressureTestAttribute.groupIDs, pressureTestAttribute.messageNumber, time.Duration(pressureTestAttribute.timeInterval)*time.Millisecond)
		time.Sleep(time.Second)
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
	recvUserID := "register_test_0"
	count := 1001
	var sendUserIDs []string
	for i := 1; i <= count; i++ {
		sendUserIDs = append(sendUserIDs, fmt.Sprintf("register_test_%v", i))
	}
	p := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	p.WithTimer(p.PressureSendMsgs2)(sendUserIDs, []string{recvUserID}, 1, 100*time.Millisecond)
}
