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

package client

import (
	"encoding/json"
	"open_im_sdk/pkg/network"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"sync"

	"github.com/gorilla/websocket"
)

type IMClient struct {
	Token    string
	UserID   string
	ApiURL   string
	JssdkURL string
	Platform int
	Conn     *websocket.Conn
	Lock     sync.Mutex
}

func NewIMClient(token, userID, apiURL, jssdkURL string, platform int) *IMClient {
	return &IMClient{
		Token:    token,
		UserID:   userID,
		ApiURL:   apiURL,
		JssdkURL: jssdkURL,
		Platform: platform,
	}
}

func (i *IMClient) GetToken() (string, error) {
	req := struct {
		Secret      string `json:"secret"`
		Platform    int32  `json:"platform"`
		UserID      string `json:"userID"`
		OperationID string `json:"operationID"`
	}{Secret: "tuoyun", Platform: int32(i.Platform), UserID: i.UserID, OperationID: utils.OperationIDGenerator()}
	content, err := network.Post2Api(i.ApiURL+"/auth/user_token", req, "")
	if err != nil {
		return "", err
	}
	type respToken struct {
		Data struct {
			ExpiredTime int64  `json:"expiredTime"`
			Token       string `json:"token"`
			Uid         string `json:"uid"`
		}
		ErrCode int    `json:"errCode"`
		ErrMsg  string `json:"errMsg"`
	}
	resp := respToken{}
	err = json.Unmarshal(content, &resp)
	return resp.Data.Token, err
}

func (i *IMClient) GetALLUserIDList() ([]string, error) {
	req := struct {
		OperationID string `json:"operationID"`
	}{OperationID: utils.OperationIDGenerator()}
	content, err := network.Post2Api(i.ApiURL+"/user/get_all_users_uid", req, i.Token)
	if err != nil {
		return nil, err
	}
	type Resp struct {
		UserIDList []string `json:"data"`
	}
	resp := Resp{}
	err = json.Unmarshal(content, &resp)
	return resp.UserIDList, err
}

func (i *IMClient) WsLogin() error {
	loginContent := struct {
		UserID string `json:"userID"`
		Token  string `json:"token"`
	}{UserID: i.UserID, Token: i.Token}
	return i.writeMessage(i.getWsReq("Login", 1, loginContent))
}

func (i *IMClient) WsLogout() error {
	return i.writeMessage(i.getWsReq("Logout", 0, nil))
}

func (i *IMClient) GetLoginStatus() error {
	return i.writeMessage(i.getWsReq("GetLoginStatus", 0, nil))
}

func (i *IMClient) CreateTextMessage(userID string) error {
	msg := server_api_params.MsgData{
		SendID:           i.UserID,
		RecvID:           "MTc3MjYzNzg0Mjg=",
		GroupID:          "",
		SenderPlatformID: int32(i.Platform),
		ClientMsgID:      utils.GetMsgID(i.UserID),
		CreateTime:       utils.GetCurrentTimestampByMill(),
		SendTime:         utils.GetCurrentTimestampByMill(),
		SessionType:      1,
		MsgFrom:          100,
		ContentType:      101,
		OfflinePushInfo:  &server_api_params.OfflinePushInfo{Title: "offlinePush"},
	}
	return i.writeMessage(i.getWsReq("CreateTextMessage", 1, msg))
}

func (i *IMClient) SendMsg(userID string, msg interface{}) error {
	return i.writeMessage(i.getWsReq("SendMessage", 1, msg))
}

func (i *IMClient) GetSelfUserInfo() error {
	return i.writeMessage(i.getWsReq("GetSelfUserInfo", 0, nil))
}

func (i *IMClient) GetAllConversationList() error {
	return i.writeMessage(i.getWsReq("GetAllConversationList", 0, nil))
}

func (i *IMClient) GetFriendList() error {
	return i.writeMessage(i.getWsReq("GetFriendList", 0, nil))
}

func (i *IMClient) GetRecvFriendApplicationList() error {
	return i.writeMessage(i.getWsReq("GetRecvFriendApplicationList", 0, nil))
}

func (i *IMClient) GetSendFriendApplicationList() error {
	return i.writeMessage(i.getWsReq("GetSendFriendApplicationList", 0, nil))
}

func (i *IMClient) GetJoinedGroupList() error {
	return i.writeMessage(i.getWsReq("GetJoinedGroupList", 0, nil))
}

func (i *IMClient) GetRecvGroupApplicationList() error {
	return i.writeMessage(i.getWsReq("GetRecvFriendApplicationList", 0, nil))
}

func (i *IMClient) GetSendGroupApplicationList() error {
	return i.writeMessage(i.getWsReq("GetSendFriendApplicationList", 0, nil))
}

func (i *IMClient) GetBlackList() error {
	return i.writeMessage(i.getWsReq("GetJoinedGroupList", 0, nil))
}

func (i *IMClient) writeMessage(bytes []byte) error {
	i.Lock.Lock()
	defer i.Lock.Unlock()
	return i.Conn.WriteMessage(1, bytes)
}

func (i *IMClient) getWsReq(event string, batch int, data interface{}) []byte {
	type Req struct {
		ReqFuncName string `json:"reqFuncName" `
		OperationID string `json:"operationID"`
		Data        string `json:"data"`
		UserID      string `json:"userID"`
		Batch       int    `json:"batchMsg,omitempty"`
	}
	req := Req{
		ReqFuncName: event,
		OperationID: i.UserID + utils.OperationIDGenerator(),
		UserID:      i.UserID,
	}
	req.Batch = batch
	var bytes []byte
	if data != nil {
		bytes, _ := json.Marshal(data)
		req.Data = string(bytes)
	}
	bytes, _ = json.Marshal(req)
	// log.Println("send:", string(bytes))
	return bytes
}
