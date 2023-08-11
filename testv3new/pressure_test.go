package testv3new

import (
	"fmt"
	"github.com/OpenIMSDK/tools/mcontext"
	"open_im_sdk/pkg/utils"
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
func Test_CreateConversationsAndSendMessages(t *testing.T) {
	if err := log.InitFromConfig("sdk.log", "sdk", 6, true, false, "", 2, 24); err != nil {
		panic(err)
	}
	recvID := "6959062403"
	conversationNum := 3
	onePeopleMessageNum := 100
	pressureTester := NewPressureTester(APIADDR, WSADDR, SECRET, Admin)
	ctx := pressureTester.NewAdminCtx()
	ctx = mcontext.SetOperationID(ctx, utils.OperationIDGenerator())
	fixedUserIDs := []string{"register_test_1", "register_test_2", "register_test_3"}
	pressureTester.CreateConversationsAndBatchSendMsg(ctx, conversationNum, onePeopleMessageNum, recvID, fixedUserIDs)
	time.Sleep(time.Minute * 10)
}
func Test_CreateConversationsAndSendGroupMessages(t *testing.T) {
	if err := log.InitFromConfig("sdk.log", "sdk", 6, true, false, "./", 2, 24); err != nil {
		panic(err)
	}
	groupID := "227809258"
	conversationNum := 3
	onePeopleMessageNum := 1000
	pressureTester := NewPressureTester(APIADDR, WSADDR, SECRET, Admin)
	ctx := pressureTester.NewAdminCtx()
	ctx = mcontext.SetOperationID(ctx, utils.OperationIDGenerator())
	fixedUserIDs := []string{"register_test_1", "register_test_2", "register_test_3"}
	pressureTester.CreateConversationsAndBatchSendGroupMsg(ctx, conversationNum, onePeopleMessageNum, groupID, fixedUserIDs)
	time.Sleep(time.Minute * 10)
}
func Test_CreateGroup(t *testing.T) {
	count := 9
	ownerUserID := "6959062403"
	p := NewPressureTester(APIADDR, WSADDR, SECRET, Admin)
	ctx := p.NewAdminCtx()
	ctx = mcontext.SetOperationID(ctx, utils.OperationIDGenerator())
	token, err := p.testUserMananger.GetToken(ctx, ownerUserID, p.platformID)
	if err != nil {
		t.Fatal(err)
	}
	var userIDs []string

	for i := 1; i <= count; i++ {
		userIDs = append(userIDs, fmt.Sprintf("register_test_%v", i))
	}
	//if err := p.testUserMananger.RegisterUsers(ctx, userIDs...); err != nil {
	//	t.Fatal(err)
	//}
	ctx = p.NewCtx(ownerUserID, token)
	ctx = mcontext.SetOperationID(ctx, utils.OperationIDGenerator())
	err = p.testUserMananger.CreateGroup(ctx, "", ownerUserID, userIDs, fmt.Sprintf("group_test_%v", "Gordon"))
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("create group success")
	}
}
