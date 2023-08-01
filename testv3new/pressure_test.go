package testv3new

import (
	"fmt"
	"testing"
	"time"

	"github.com/OpenIMSDK/tools/log"
)

var (
	pressureTestAttribute PressureTestAttribute
)

func init() {
	pressureTestAttribute.InitWithFlag()
	//pressureTestAttribute.recvUserIDs = []string{"6680650275"}
	//pressureTestAttribute.sendUserIDs = []string{"8430973211", "4098531159", "3171794400"}
	//pressureTestAttribute.messageNumber = 10
	//pressureTestAttribute.timeInterval = 100

	if err := log.InitFromConfig("sdk.log", "sdk", 4,
		true, false, "./chat_log", 2, 24); err != nil {
		panic(err)
	}
}

func TestPressureTester_PressureSendMsgs(t *testing.T) {
	ParseFlag()
	p := NewPressureTester(APIADDR, WSADDR)
	for i := 0; i < 10; i++ {
		p.WithTimer(p.PressureSendMsgs2)(pressureTestAttribute.sendUserIDs, pressureTestAttribute.recvUserIDs, pressureTestAttribute.messageNumber, time.Duration(pressureTestAttribute.timeInterval)*time.Millisecond)
		time.Sleep(time.Second)
	}
	// time.Sleep(1000 * time.Second)
}

func TestPressureTester_PressureSendGroupMsgs(t *testing.T) {
	ParseFlag()
	p := NewPressureTester(APIADDR, WSADDR)
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
	p := NewPressureTester(APIADDR, WSADDR)
	p.WithTimer(p.PressureSendMsgs)(sendUserID, recvUserIDs, 1, 100*time.Millisecond)
}

func TestPressureTester_PressureSendMsgs2(t *testing.T) {
	recvUserID := "6680650275"
	var sendUserIDs []string
	for i := 1; i <= 100; i++ {
		sendUserIDs = append(sendUserIDs, fmt.Sprintf("register_%v", i))
	}
	p := NewPressureTester(APIADDR, WSADDR)
	p.WithTimer(p.PressureSendMsgs2)(sendUserIDs, []string{recvUserID}, 1000, 100*time.Millisecond)
}
