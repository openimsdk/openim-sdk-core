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
	"fmt"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/testv3/funcation"
	"testing"
	"time"

	"github.com/OpenIMSDK/tools/log"
)

// 定点发送消息（好友关系）
func Test_SendMsg(t *testing.T) {
	// 用户A 向 用户B 发送一条消息
	userA := "2935954421"
	userB := "4587931616"

	funcation.LoginOne(userA)
	ctx := funcation.CreateCtx(userA)

	msg := fmt.Sprintf("%v send to %v a message", userA, userB)
	t.Log("prefix ", msg)
	funcation.SendMsg(ctx, userA, userB, "", msg)
}

// 定点发送消息（非好友关系）
func Test_SendMsgNoFriend(t *testing.T) {
	// 用户A 向 用户B 发送一条消息
	userA := "bantanger"
	userB := "9003169405"
	funcation.LoginOne(userA)
	ctx := funcation.CreateCtx(userA)

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
	ctx := funcation.CreateCtx(userA)

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
	ctx := funcation.CreateCtx(userA)

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
	userB := "8423809271"

	funcation.LoginOne(userA)
	ctx := funcation.CreateCtx(userA)
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
	ctx := funcation.CreateCtx(userA)

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

// 十万人群聊测试
func Test_SendMsgByGroup2(t *testing.T) {
	funcation.AllLoginMgr = make(map[string]*funcation.CoreNode)
	count := 100000
	var uid []string
	for i := 1; i <= count; i++ {
		uid = append(uid, fmt.Sprintf("register_test_%v", i))
	}
	funcation.LoginBatch(uid)
	groupID := "3747979639"
	for i := 1; i <= count; i++ {
		ctx := funcation.CreateCtx(uid[i])

		time.Sleep(time.Duration(200) * time.Millisecond)
		msg := fmt.Sprintf("%v send to %v message by %d ", uid[i], groupID, i)
		funcation.SendMsg(ctx, uid[i], "", groupID, msg)
		log.ZDebug(ctx, "msg prefix", "msg", msg)
	}
}

// 十万人群聊测试
func Test_SendMsgByGroup3(t *testing.T) {
	funcation.AllLoginMgr = make(map[string]*funcation.CoreNode)
	count := 1
	// groupID := "3747979639"
	groupID := "3686091035"
	for i := 1; i <= count; i++ {
		uid := fmt.Sprintf("register_test_%v", i)
		funcation.LoginOne(uid)
		ctx := funcation.CreateCtx(uid)

		funcation.AllLoginMgr[uid].Mgr.Conversation().GetOneConversation(ctx, 2, uid)

		time.Sleep(time.Duration(200) * time.Millisecond)
		msg := fmt.Sprintf("%v send to %v message by %d ", uid, groupID, i)
		funcation.SendMsg(ctx, uid, "", groupID, msg)
		log.ZDebug(ctx, "msg prefix", "msg", msg)
	}
}

// 管理员邀请群成员进行群聊测试
func Test_SendMsgByGroup_One(t *testing.T) {
	funcation.AllLoginMgr = make(map[string]*funcation.CoreNode)
	uid := "register_test_1002"
	groupID := "226506521"
	// groupID := "4238406538"
	// groupID := "2592040611"
	funcation.LoginOne(uid)

	ctx := funcation.CreateCtx(uid)
	time.Sleep(time.Duration(200) * time.Millisecond)

	// 管理员邀请进群
	funcation.AllLoginMgr[uid].Mgr.Group().JoinGroup(ctx, groupID, "", constant.JoinBySearch)

	msg := fmt.Sprintf("%v send to %v message", uid, groupID)
	funcation.SendMsg(ctx, uid, "", groupID, msg)
	log.ZDebug(ctx, "msg prefix", "msg", msg)
}

// 管理员邀请大量群成员进行群聊测试
func Test_SendMsgByGroup_Batch(t *testing.T) {
	count := 1000

	// groupID := "780048154"
	// groupID := "3159824577"
	groupID := "1938828611"
	var uidList []string
	for i := 0; i <= count; i++ {
		uid := fmt.Sprintf("register_test_%v", i+1)
		uidList = append(uidList, uid)
	}
	// 管理员批量邀请进群
	adminUID := "openIM123456"
	funcation.LoginOne(adminUID)
	ctx := funcation.CreateCtx(adminUID)
	err := funcation.AllLoginMgr[adminUID].Mgr.Group().InviteUserToGroup(ctx, groupID, "", uidList)
	if err != nil {
		t.Error("invite user fails")
		return
	}

	for i := 0; i <= count; i++ {
		uid := uidList[i]
		funcation.LoginOne(uid)
		ctx := funcation.CreateCtx(uid)
		// time.Sleep(time.Duration(200) * time.Millisecond)

		msg := fmt.Sprintf("%v send to %v message", uid, groupID)
		funcation.SendMsg(ctx, uid, "", groupID, msg)
		log.ZDebug(ctx, "msg prefix", "msg", msg)
	}
}
