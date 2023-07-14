// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// @Author BanTanger 2023/7/10 15:30:00
package testv3

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/testv3/funcation"
	"testing"
	"time"
)

// 定点发送消息（好友关系）
func Test_SendMsg(t *testing.T) {
	// 用户A 向 用户B 发送一条消息
	userA := "2935954421"
	userB := "4587931616"

	funcation.LoginOne(userA)
	operationID := utils.OperationIDGenerator()
	ctx := ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID: userA,
		Token:  funcation.AllLoginMgr[userA].Token,
	})
	ctx = ccontext.WithOperationID(ctx, operationID)
	ctx = ccontext.WithSendMessageCallback(ctx, funcation.TestSendMsgCallBack{
		OperationID: operationID,
	})

	msg := fmt.Sprintf("%v send to %v a message", userA, userB)
	t.Log("prefix ", msg)
	funcation.SendMsg(ctx, userA, userB, "", msg)
}

// 定点发送消息（非好友关系）
func Test_SendMsgNoFriend(t *testing.T) {
	// 用户A 向 用户B 发送一条消息
	userA := "register_test_526"
	userB := "2935954421"
	funcation.LoginOne(userA)
	operationID := utils.OperationIDGenerator()
	ctx := ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID:   userA,
		Token:    funcation.AllLoginMgr[userA].Token,
		IMConfig: funcation.Config,
	})
	ctx = ccontext.WithOperationID(ctx, operationID)
	ctx = ccontext.WithSendMessageCallback(ctx, funcation.TestSendMsgCallBack{
		OperationID: operationID,
	})

	funcation.AllLoginMgr[userA].Mgr.Conversation().GetOneConversation(ctx, 1, userB)
	msg := fmt.Sprintf("%v send to %v a message", userA, userB)
	t.Log("prefix ", msg)
	funcation.SendMsg(ctx, userA, userB, "", msg)
}

// 定点发送消息（非好友关系）:为千人会话发送一条消息
func Test_SendMsgNoFriendBatch(t *testing.T) {
	// 批量注册用户
	count := 1000
	var users []funcation.Users
	for i := 0; i < count; i++ {
		users = append(users, funcation.Users{
			Uid:      fmt.Sprintf("register_test_%d", i),
			Nickname: fmt.Sprintf("register_test_%d", i),
			FaceUrl:  "",
		})
	}
	funcation.RegisterBatch(users)
	time.Sleep(500 * time.Millisecond)
	var userIDList []string
	for i := range users {
		userIDList = append(userIDList, users[i].Uid)
	}
	funcation.LoginBatch(userIDList)
	time.Sleep(500 * time.Millisecond)

	// 用户A 向 大量用户 发送一条消息
	userA := "6506148011"
	funcation.LoginOne(userA)
	operationID := utils.OperationIDGenerator()
	ctx := ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID:   userA,
		Token:    funcation.AllLoginMgr[userA].Token,
		IMConfig: funcation.Config,
	})
	ctx = ccontext.WithOperationID(ctx, operationID)
	ctx = ccontext.WithSendMessageCallback(ctx, funcation.TestSendMsgCallBack{
		OperationID: operationID,
	})

	for i := range users {
		msg := fmt.Sprintf("%v send to %v a message", userA, users[i].Uid)
		t.Log("prefix ", msg)
		funcation.AllLoginMgr[userA].Mgr.Conversation().GetOneConversation(ctx, 1, users[i].Uid)
		funcation.SendMsg(ctx, userA, users[i].Uid, "", msg)
	}
}

// 定点发送消息（非好友关系）:为千人会话发送千条消息
func Test_SendMsgNoFriendBatch2(t *testing.T) {
	// 用户A 向 大量用户 发送一条消息
	userA := "bantanger"
	funcation.LoginOne(userA)
	operationID := utils.OperationIDGenerator()
	ctx := ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID:   userA,
		Token:    funcation.AllLoginMgr[userA].Token,
		IMConfig: funcation.Config,
	})
	ctx = ccontext.WithOperationID(ctx, operationID)
	ctx = ccontext.WithSendMessageCallback(ctx, funcation.TestSendMsgCallBack{
		OperationID: operationID,
	})

	for i := range funcation.AllLoginMgr {
		users := fmt.Sprintf("register_test_%v", i)
		funcation.AllLoginMgr[userA].Mgr.Conversation().GetOneConversation(ctx, 1, users)
		for i := 1; i <= 1000; i++ {
			// time.Sleep(200 * time.Millisecond)
			msg := fmt.Sprintf("%v send to %v message by %d ", userA, users, i)
			log.ZInfo(ctx, "send message", "userA", userA, "userB", users, "i", i)
			funcation.SendMsg(ctx, userA, users, "", msg)
			log.ZDebug(ctx, "msg prefix", "msg", msg)
		}
	}
}

// 批量发送消息：P2P
func Test_SendMsgBatch(t *testing.T) {
	// 用户A 向 用户B 发送多条消息
	userA := "register_test_526"
	userB := "2935954421"

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
	log.ZInfo(ctx, "start send message")
	funcation.AllLoginMgr[userA].Mgr.Conversation().GetOneConversation(ctx, 1, userB)
	for i := 1; i <= 100000; i++ {
		// time.Sleep(200 * time.Millisecond)
		msg := fmt.Sprintf("%v send to %v message by %d ", userA, userB, i)
		log.ZInfo(ctx, "send message", "userA", userA, "userB", userB, "i", i)
		funcation.SendMsg(ctx, userA, userB, "", msg)
		log.ZDebug(ctx, "msg prefix", "msg", msg)

	}
	log.ZInfo(ctx, "end send message")
	time.Sleep(1000 * time.Second)
}

// 批量发送消息：群聊消息
func Test_SendMsgByGroup(t *testing.T) {
	// 用户A 向 群组 发送多条消息
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
		time.Sleep(time.Duration(200) * time.Millisecond)
		msg := fmt.Sprintf("%v send to %v message by %d ", userA, group, i)
		funcation.SendMsg(ctx, userA, "", group, msg)
		log.ZDebug(ctx, "msg prefix", "msg", msg)
	}

	// conversationID := "si_" + userA + userB
	// count := 20
	//
	// // 获取会话查看消息是否抵达
	// funcation.GetConversation(userA, conversationID, message.ClientMsgID, count)
	time.Sleep(10 * time.Second)
}
