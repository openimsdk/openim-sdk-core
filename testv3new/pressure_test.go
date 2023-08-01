package testv3new

import (
	"fmt"
	"github.com/OpenIMSDK/tools/log"
	"open_im_sdk/testv3new/constant"
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
	var sendUserIDs []string
	var recvUserIDs []string
	for i := 0; i < pressureTestAttribute.sendNums; i++ {
		sendUserIDs = append(sendUserIDs, fmt.Sprintf("register_test_%v", i))
	}
	for i := 0; i < pressureTestAttribute.recvNums; i++ {
		recvUserIDs = append(recvUserIDs, fmt.Sprintf("register_test_%v", i+pressureTestAttribute.sendNums))
	}
	p := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
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
	p := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	for i := 0; i < 10; i++ {
		p.WithTimer(p.PressureSendGroupMsgs2)(sendUserIDs, groupIDs, pressureTestAttribute.messageNumber, time.Duration(pressureTestAttribute.timeInterval)*time.Millisecond)
		time.Sleep(time.Duration(pressureTestAttribute.timeInterval) * time.Second)
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
	messageNum := 1
	LoopNumber := 10
	var sendUserIDs []string
	for i := 1; i <= count; i++ {
		sendUserIDs = append(sendUserIDs, fmt.Sprintf("register_test_%v", i))
	}
	p := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
	for i := 0; i < LoopNumber; i++ {
		p.WithTimer(p.PressureSendMsgs2)(sendUserIDs, []string{recvUserID}, messageNum, 100*time.Millisecond)
		time.Sleep(time.Second)
	}
}

func Test_CreateGroup(t *testing.T) {
	count := 1000
	ownerUserID := "register_test_0"
	for i := 0; i < count; i++ {
		p := NewPressureTester(testcore.APIADDR, testcore.WSADDR)
		err := p.CreateGroup(fmt.Sprintf("group_test_%v", i), ownerUserID, []string{constant.DefaultGroupMember}, fmt.Sprintf("group_test_%v", i))
		if err != nil {
			t.Fatal(err)
		}
	}
}
