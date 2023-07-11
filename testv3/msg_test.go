// @Author BanTanger 2023/7/10 15:30:00
package testv3

import (
	"context"
	"fmt"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/testv3/funcation"
	"testing"
	"time"
)

// 定点发送消息
func Test_SendMsg(t *testing.T) {
	// 用户A 向 用户B 发送一条消息
	userA := "2935954421"
	userB := "4587931616"
	msg := fmt.Sprintf("%v send to %v a message", userA, userB)
	t.Log("prefix ", msg)
	funcation.LoginOne(userA)
	// funcation.LoginOne(userB)
	operationID := utils.OperationIDGenerator()
	ctx := ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID: userA,
		Token:  funcation.AllLoginMgr[userA].Token,
	})
	ctx = ccontext.WithOperationID(ctx, operationID)
	ctx = ccontext.WithSendMessageCallback(ctx, funcation.TestSendMsgCallBack{
		OperationID: operationID,
	})

	funcation.SendMsg(ctx, userA, userB, "", msg)

	// conversationID := "si_" + userA + userB
	// count := 20
	//
	// // 获取会话查看消息是否抵达
	// funcation.GetConversation(userA, conversationID, message.ClientMsgID, count)
	time.Sleep(10 * time.Second)
}

func Test_SendMsgBatch(t *testing.T) {
	// 用户A 向 用户B 发送一条消息
	userA := "2935954421"
	userB := "4587931616"

	funcation.LoginOne(userA)
	// funcation.LoginOne(userB)
	operationID := utils.OperationIDGenerator()
	ctx := ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID: userA,
		Token:  funcation.AllLoginMgr[userA].Token,
	})
	ctx = ccontext.WithOperationID(ctx, operationID)
	ctx = ccontext.WithSendMessageCallback(ctx, funcation.TestSendMsgCallBack{
		OperationID: operationID,
	})
	time.Sleep(10 * time.Millisecond)
	for i := 1; i <= 100000; i++ {
		// time.Sleep(time.Duration(100) * time.Millisecond)
		msg := fmt.Sprintf("%v send to %v message by %d ", userA, userB, i)
		funcation.SendMsg(ctx, userA, userB, "", msg)
		t.Log("prefix " + msg)
	}
	time.Sleep(1000 * time.Second)
}

func Test_SendMsgByGroup(t *testing.T) {
	// 用户A 向 用户B 发送一条消息
	userA := "2935954421"
	group := "1347996360"
	msg := fmt.Sprintf("%v send to %v a message", userA, group)
	t.Log("prefix ", msg)
	funcation.LoginOne(userA)
	// funcation.LoginOne(userB)
	operationID := utils.OperationIDGenerator()
	ctx := ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID: userA,
		Token:  funcation.AllLoginMgr[userA].Token,
	})
	ctx = ccontext.WithOperationID(ctx, operationID)
	ctx = ccontext.WithSendMessageCallback(ctx, funcation.TestSendMsgCallBack{
		OperationID: operationID,
	})

	for i := 1; i <= 100000; i++ {
		// time.Sleep(time.Duration(100) * time.Millisecond)
		msg := fmt.Sprintf("%v send to %v message by %d ", userA, group, i)
		funcation.SendMsg(ctx, userA, "", group, msg)
		t.Log("prefix " + msg)
	}

	// conversationID := "si_" + userA + userB
	// count := 20
	//
	// // 获取会话查看消息是否抵达
	// funcation.GetConversation(userA, conversationID, message.ClientMsgID, count)
	time.Sleep(10 * time.Second)
}

func Test_SendMsgBatch(t *testing.T) {
	// 用户A 向 用户B 发送一条消息
	userA := "2935954421"
	userB := "4587931616"

	funcation.LoginOne(userA)
	// funcation.LoginOne(userB)
	operationID := utils.OperationIDGenerator()
	ctx := ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID: userA,
		Token:  funcation.AllLoginMgr[userA].Token,
	})
	ctx = ccontext.WithOperationID(ctx, operationID)
	ctx = ccontext.WithSendMessageCallback(ctx, funcation.TestSendMsgCallBack{
		OperationID: operationID,
	})
	for i := 0; i < 100000; i++ {
		msg := fmt.Sprintf("%v send to %v message by %d ", userA, userB, i)
		// funcation.SendMsg(ctx, userA, userB, "", msg)
		t.Log("prefix " + msg)
		time.Sleep(10 * time.Millisecond)
	}
	// recvTime := funcation.SendSuccAllMsg[userA]
	// t.Log(recvTime)
	time.Sleep(1000 * time.Second)
}
