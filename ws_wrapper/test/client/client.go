package client

import (
	"encoding/json"
	"fmt"
	"open_im_sdk/pkg/network"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"

	"github.com/gorilla/websocket"
)

type IMClient struct {
	Token    string
	UserID   string
	ApiURL   string
	JssdkURL string
	Platform int
	Conn     *websocket.Conn
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

func (i *IMClient) GetUserIDList() ([]string, error) {
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

func (i *IMClient) SendTestMsg(conn *websocket.Conn) error {
	return conn.WriteMessage(1, []byte("s"))
}

func (i *IMClient) WsLogin() error {
	loginContent := struct {
		UserID string `json:"userID"`
		Token  string `json:"token"`
	}{UserID: i.UserID, Token: i.Token}
	return i.Conn.WriteMessage(1, i.getWsReq("Login", 1, loginContent))
}

func (i *IMClient) WsLogout() error {
	return i.Conn.WriteMessage(1, i.getWsReq("Logout", 0, nil))
}

func (i *IMClient) GetLoginStatus() error {
	return i.Conn.WriteMessage(1, i.getWsReq("GetLoginStatus", 0, nil))
}

func (i *IMClient) SendMsg(userID string) error {
	msg := server_api_params.MsgData{
		SendID:           i.UserID,
		RecvID:           "MTc3MjYzNzg0Mjg=",
		SenderPlatformID: int32(i.Platform),
		ClientMsgID:      utils.GetMsgID(i.UserID),
		CreateTime:       utils.GetCurrentTimestampByMill(),
		SessionType:      1,
		MsgFrom:          100,
		ContentType:      101,
	}
	return i.Conn.WriteMessage(1, i.getWsReq("SendMessage", 1, msg))
}

func (i *IMClient) GetSelfUserInfo() error {
	return i.Conn.WriteMessage(1, i.getWsReq("GetSelfUserInfo", 0, nil))
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
	fmt.Println(string(bytes))
	return bytes
}
