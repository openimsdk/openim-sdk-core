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

package testv2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/OpenIMSDK/protocol/constant"
	"io"
	"math/rand"
	"net/http"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/ccontext"
	"strconv"
	"time"

	"github.com/OpenIMSDK/tools/log"
)

var (
	ctx context.Context
)

func init() {
	fmt.Println("------------------------>>>>>>>>>>>>>>>>>>> test v2 init funcation <<<<<<<<<<<<<<<<<<<------------------------")
	rand.Seed(time.Now().UnixNano())
	listner := &OnConnListener{}
	config := getConf(APIADDR, WSADDR)
	config.DataDir = ""
	configData, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	isInit := open_im_sdk.InitSDK(listner, "test", string(configData))
	if !isInit {
		panic("init sdk failed")
	}
	ctx = open_im_sdk.UserForSDK.Context()
	ctx = ccontext.WithOperationID(ctx, "initOperationID_"+strconv.Itoa(int(time.Now().UnixMilli())))
	token, err := GetUserToken(ctx, UserID)
	if err != nil {
		panic(err)
	}
	if err := open_im_sdk.UserForSDK.Login(ctx, UserID, token); err != nil {
		panic(err)
	}
	open_im_sdk.UserForSDK.SetListenerForService(&onListenerForService{ctx: ctx})
	open_im_sdk.UserForSDK.SetConversationListener(&onConversationListener{ctx: ctx})
	open_im_sdk.UserForSDK.SetGroupListener(&onGroupListener{ctx: ctx})
	open_im_sdk.UserForSDK.SetAdvancedMsgListener(&onAdvancedMsgListener{ctx: ctx})
	open_im_sdk.UserForSDK.SetFriendListener(&onFriendListener{ctx: ctx})
	open_im_sdk.UserForSDK.SetUserListener(&onUserListener{ctx: ctx})
	time.Sleep(time.Second * 10)
}

func GetUserToken(ctx context.Context, userID string) (string, error) {
	jsonReqData, err := json.Marshal(map[string]any{
		"userID":     userID,
		"platformID": constant.LinuxPlatformID,
		"secret":     "openIM123",
	})
	if err != nil {
		return "", err
	}
	path := APIADDR + "/auth/user_token"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(jsonReqData))
	if err != nil {
		return "", err
	}
	req.Header.Set("operationID", ctx.Value("operationID").(string))
	client := http.Client{Timeout: time.Second * 3}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	type Result struct {
		ErrCode int    `json:"errCode"`
		ErrMsg  string `json:"errMsg"`
		ErrDlt  string `json:"errDlt"`
		Data    struct {
			Token             string `json:"token"`
			ExpireTimeSeconds int    `json:"expireTimeSeconds"`
		} `json:"data"`
	}
	var result Result
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if result.ErrCode != 0 {
		return "", fmt.Errorf("errCode:%d, errMsg:%s, errDlt:%s", result.ErrCode, result.ErrMsg, result.ErrDlt)
	}
	return result.Data.Token, nil
}

type onListenerForService struct {
	ctx context.Context
}

func (o *onListenerForService) OnGroupApplicationAdded(groupApplication string) {
	log.ZInfo(o.ctx, "OnGroupApplicationAdded", "groupApplication", groupApplication)
}

func (o *onListenerForService) OnGroupApplicationAccepted(groupApplication string) {
	log.ZInfo(o.ctx, "OnGroupApplicationAccepted", "groupApplication", groupApplication)
}

func (o *onListenerForService) OnFriendApplicationAdded(friendApplication string) {
	log.ZInfo(o.ctx, "OnFriendApplicationAdded", "friendApplication", friendApplication)
}

func (o *onListenerForService) OnFriendApplicationAccepted(groupApplication string) {
	log.ZInfo(o.ctx, "OnFriendApplicationAccepted", "groupApplication", groupApplication)
}

func (o *onListenerForService) OnRecvNewMessage(message string) {
	log.ZInfo(o.ctx, "OnRecvNewMessage", "message", message)
}

type onConversationListener struct {
	ctx context.Context
}

func (o *onConversationListener) OnSyncServerStart() {
	log.ZInfo(o.ctx, "OnSyncServerStart")
}

func (o *onConversationListener) OnSyncServerFinish() {
	log.ZInfo(o.ctx, "OnSyncServerFinish")
}

func (o *onConversationListener) OnSyncServerFailed() {
	log.ZInfo(o.ctx, "OnSyncServerFailed")
}

func (o *onConversationListener) OnNewConversation(conversationList string) {
	log.ZInfo(o.ctx, "OnNewConversation", "conversationList", conversationList)
}

func (o *onConversationListener) OnConversationChanged(conversationList string) {
	log.ZInfo(o.ctx, "OnConversationChanged", "conversationList", conversationList)
}

func (o *onConversationListener) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
	log.ZInfo(o.ctx, "OnTotalUnreadMessageCountChanged", "totalUnreadCount", totalUnreadCount)
}

type onGroupListener struct {
	ctx context.Context
}

func (o *onGroupListener) OnGroupDismissed(groupInfo string) {
	log.ZInfo(o.ctx, "OnGroupDismissed", "groupInfo", groupInfo)
}

func (o *onGroupListener) OnJoinedGroupAdded(groupInfo string) {
	log.ZInfo(o.ctx, "OnJoinedGroupAdded", "groupInfo", groupInfo)
}

func (o *onGroupListener) OnJoinedGroupDeleted(groupInfo string) {
	log.ZInfo(o.ctx, "OnJoinedGroupDeleted", "groupInfo", groupInfo)
}

func (o *onGroupListener) OnGroupMemberAdded(groupMemberInfo string) {
	log.ZInfo(o.ctx, "OnGroupMemberAdded", "groupMemberInfo", groupMemberInfo)
}

func (o *onGroupListener) OnGroupMemberDeleted(groupMemberInfo string) {
	log.ZInfo(o.ctx, "OnGroupMemberDeleted", "groupMemberInfo", groupMemberInfo)
}

func (o *onGroupListener) OnGroupApplicationAdded(groupApplication string) {
	log.ZInfo(o.ctx, "OnGroupApplicationAdded", "groupApplication", groupApplication)
}

func (o *onGroupListener) OnGroupApplicationDeleted(groupApplication string) {
	log.ZInfo(o.ctx, "OnGroupApplicationDeleted", "groupApplication", groupApplication)
}

func (o *onGroupListener) OnGroupInfoChanged(groupInfo string) {
	log.ZInfo(o.ctx, "OnGroupInfoChanged", "groupInfo", groupInfo)
}

func (o *onGroupListener) OnGroupMemberInfoChanged(groupMemberInfo string) {
	log.ZInfo(o.ctx, "OnGroupMemberInfoChanged", "groupMemberInfo", groupMemberInfo)
}

func (o *onGroupListener) OnGroupApplicationAccepted(groupApplication string) {
	log.ZInfo(o.ctx, "OnGroupApplicationAccepted", "groupApplication", groupApplication)
}

func (o *onGroupListener) OnGroupApplicationRejected(groupApplication string) {
	log.ZInfo(o.ctx, "OnGroupApplicationRejected", "groupApplication", groupApplication)
}

type onAdvancedMsgListener struct {
	ctx context.Context
}

func (o *onAdvancedMsgListener) OnRecvOfflineNewMessage(message string) {
	//TODO implement me
	panic("implement me")
}

func (o *onAdvancedMsgListener) OnMsgDeleted(message string) {
	log.ZInfo(o.ctx, "OnMsgDeleted", "message", message)
}

//funcation (o *onAdvancedMsgListener) OnMsgDeleted(messageList string) {
//	log.ZInfo(o.ctx, "OnRecvOfflineNewMessages", "messageList", messageList)
//}
//
//funcation (o *onAdvancedMsgListener) OnMsgDeleted(message string) {
//	log.ZInfo(o.ctx, "OnMsgDeleted", "message", message)
//}

func (o *onAdvancedMsgListener) OnRecvOfflineNewMessages(messageList string) {
	log.ZInfo(o.ctx, "OnRecvOfflineNewMessages", "messageList", messageList)
}

func (o *onAdvancedMsgListener) OnRecvNewMessage(message string) {
	log.ZInfo(o.ctx, "OnRecvNewMessage", "message", message)
}

func (o *onAdvancedMsgListener) OnRecvC2CReadReceipt(msgReceiptList string) {
	log.ZInfo(o.ctx, "OnRecvC2CReadReceipt", "msgReceiptList", msgReceiptList)
}

func (o *onAdvancedMsgListener) OnRecvGroupReadReceipt(groupMsgReceiptList string) {
	log.ZInfo(o.ctx, "OnRecvGroupReadReceipt", "groupMsgReceiptList", groupMsgReceiptList)
}

func (o *onAdvancedMsgListener) OnRecvMessageRevoked(msgID string) {
	log.ZInfo(o.ctx, "OnRecvMessageRevoked", "msgID", msgID)
}

func (o *onAdvancedMsgListener) OnNewRecvMessageRevoked(messageRevoked string) {
	log.ZInfo(o.ctx, "OnNewRecvMessageRevoked", "messageRevoked", messageRevoked)
}

func (o *onAdvancedMsgListener) OnRecvMessageExtensionsChanged(msgID string, reactionExtensionList string) {
	log.ZInfo(o.ctx, "OnRecvMessageExtensionsChanged", "msgID", msgID, "reactionExtensionList", reactionExtensionList)
}

func (o *onAdvancedMsgListener) OnRecvMessageExtensionsDeleted(msgID string, reactionExtensionKeyList string) {
	log.ZInfo(o.ctx, "OnRecvMessageExtensionsDeleted", "msgID", msgID, "reactionExtensionKeyList", reactionExtensionKeyList)
}

func (o *onAdvancedMsgListener) OnRecvMessageExtensionsAdded(msgID string, reactionExtensionList string) {
	log.ZInfo(o.ctx, "OnRecvMessageExtensionsAdded", "msgID", msgID, "reactionExtensionList", reactionExtensionList)
}

type onFriendListener struct {
	ctx context.Context
}

func (o *onFriendListener) OnFriendApplicationAdded(friendApplication string) {
	log.ZDebug(context.Background(), "OnFriendApplicationAdded", "friendApplication", friendApplication)
}

func (o *onFriendListener) OnFriendApplicationDeleted(friendApplication string) {
	log.ZDebug(context.Background(), "OnFriendApplicationDeleted", "friendApplication", friendApplication)
}

func (o *onFriendListener) OnFriendApplicationAccepted(friendApplication string) {
	log.ZDebug(context.Background(), "OnFriendApplicationAccepted", "friendApplication", friendApplication)
}

func (o *onFriendListener) OnFriendApplicationRejected(friendApplication string) {
	log.ZDebug(context.Background(), "OnFriendApplicationRejected", "friendApplication", friendApplication)
}

func (o *onFriendListener) OnFriendAdded(friendInfo string) {
	log.ZDebug(context.Background(), "OnFriendAdded", "friendInfo", friendInfo)
}

func (o *onFriendListener) OnFriendDeleted(friendInfo string) {
	log.ZDebug(context.Background(), "OnFriendDeleted", "friendInfo", friendInfo)
}

func (o *onFriendListener) OnFriendInfoChanged(friendInfo string) {
	log.ZDebug(context.Background(), "OnFriendInfoChanged", "friendInfo", friendInfo)
}

func (o *onFriendListener) OnBlackAdded(blackInfo string) {
	log.ZDebug(context.Background(), "OnBlackAdded", "blackInfo", blackInfo)
}

func (o *onFriendListener) OnBlackDeleted(blackInfo string) {
	log.ZDebug(context.Background(), "OnBlackDeleted", "blackInfo", blackInfo)
}

type onUserListener struct {
	ctx context.Context
}

func (o *onUserListener) OnSelfInfoUpdated(userInfo string) {
	log.ZDebug(context.Background(), "OnBlackDeleted", "blackInfo", userInfo)
}
func (o *onUserListener) OnUserStatusChanged(statusMap string) {
	log.ZDebug(context.Background(), "OnUserStatusChanged", "OnUserStatusChanged", statusMap)
}
