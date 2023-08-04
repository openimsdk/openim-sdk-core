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
	// pressureTestAttribute.InitWithFlag()
	//pressureTestAttribute.recvUserIDs = []string{"6680650275"}
	//pressureTestAttribute.sendUserIDs = []string{"8430973211", "4098531159", "3171794400"}
	//pressureTestAttribute.messageNumber = 10
	//pressureTestAttribute.timeInterval = 100
	pressureTestAttribute.InitWithFlag()
	if err := log.InitFromConfig("sdk.log", "sdk", 4,
		true, false, "./chat_log", 2, 24); err != nil {
		panic(err)
	}
}

func TestPressureTester_PressureSendMsgs(t *testing.T) {
	ParseFlag()
	var sendUserIDs []string
	var recvUserIDs []string
	for i := 0; i < pressureTestAttribute.sendNums; i++ {
		sendUserIDs = append(sendUserIDs, fmt.Sprintf("register_test_%v", i))
	}
	for i := 0; i < pressureTestAttribute.recvNums; i++ {
		recvUserIDs = append(recvUserIDs, fmt.Sprintf("register_test_%v", i+pressureTestAttribute.sendNums))
	}
	p := NewPressureTester(APIADDR, WSADDR, SECRET, Admin)
	for i := 0; i < 10; i++ {
		p.WithTimer(p.PressureSendMsgs2)(sendUserIDs, recvUserIDs, pressureTestAttribute.messageNumber, time.Duration(pressureTestAttribute.timeInterval)*time.Millisecond)
		time.Sleep(time.Duration(pressureTestAttribute.timeInterval) * time.Second)
	}
	// time.Sleep(1000 * time.Second)
}

func TestPressureTester_PressureSendGroupMsgs(t *testing.T) {
	ParseFlag()
	var sendUserIDs []string
	var groupIDs []string
	for i := 0; i < pressureTestAttribute.sendNums; i++ {
		sendUserIDs = append(sendUserIDs, fmt.Sprintf("register_test_%v", i))
	}
	for i := 0; i < pressureTestAttribute.groupNums; i++ {
		sendUserIDs = append(sendUserIDs, fmt.Sprintf("group_test_%v", i))
	}
	p := NewPressureTester(APIADDR, WSADDR, SECRET, Admin)
	for i := 0; i < 10; i++ {
		p.WithTimer(p.PressureSendGroupMsgs2)(sendUserIDs, groupIDs, pressureTestAttribute.messageNumber, time.Duration(pressureTestAttribute.timeInterval)*time.Millisecond)
		time.Sleep(time.Duration(pressureTestAttribute.timeInterval) * time.Second)
	}
}

func TestPressureTester_Conversation(t *testing.T) {

}

func TestPressureTester_PressureSendMsgs2(t *testing.T) {
	recvUserID := "register_test_0"
	count := 1001
	messageNum := 1
	LoopNumber := 10
	var sendUserIDs []string
	for i := 1; i <= count; i++ {
		sendUserIDs = append(sendUserIDs, fmt.Sprintf("register_test_%v", i))
	}
	p := NewPressureTester(APIADDR, WSADDR, SECRET, Admin)
	for i := 0; i < LoopNumber; i++ {
		p.WithTimer(p.PressureSendMsgs2)(sendUserIDs, []string{recvUserID}, messageNum, 100*time.Millisecond)
		time.Sleep(time.Second)
	}
}
