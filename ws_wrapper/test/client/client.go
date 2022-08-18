package client

import (
	"encoding/json"
	"log"
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

func (i *IMClient) SendMsg(userID string) error {
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
	return i.writeMessage(i.getWsReq("SendMessage", 1, msg))
}

// data: "{\"recvID\":\"4266290636\",\"groupID\":\"\",\"offlinePushInfo\":\"{\\\"title\\\":\\\"你有一条新消息\\\",\\\"desc\\\":\\\"\\\",\\\"ex\\\":\\\"\\\",\\\"iOSPushSound\\\":\\\"+1\\\",\\\"iOSBadgeCount\\\":true}\",\"message\":\"{\\\"clientMsgID\\\":\\\"f1dac895a848b1f2c1e061b14e62cf00\\\",\\\"createTime\\\":1660843983746,\\\"sendTime\\\":1660843983746,\\\"sessionType\\\":0,\\\"sendID\\\":\\\"4266290636\\\",\\\"msgFrom\\\":100,\\\"contentType\\\":101,\\\"platformID\\\":5,\\\"senderNickname\\\":\\\"kernal在\\\",\\\"senderFaceUrl\\\":\\\"ic_avatar_06\\\",\\\"content\\\":\\\"1\\\",\\\"seq\\\":0,\\\"isRead\\\":false,\\\"status\\\":1,\\\"offlinePush\\\":{},\\\"pictureElem\\\":{\\\"sourcePicture\\\":{\\\"size\\\":0,\\\"width\\\":0,\\\"height\\\":0},\\\"bigPicture\\\":{\\\"size\\\":0,\\\"width\\\":0,\\\"height\\\":0},\\\"snapshotPicture\\\":{\\\"size\\\":0,\\\"width\\\":0,\\\"height\\\":0}},\\\"soundElem\\\":{\\\"dataSize\\\":0,\\\"duration\\\":0},\\\"videoElem\\\":{\\\"videoSize\\\":0,\\\"duration\\\":0,\\\"snapshotSize\\\":0,\\\"snapshotWidth\\\":0,\\\"snapshotHeight\\\":0},\\\"fileElem\\\":{\\\"fileSize\\\":0},\\\"mergeElem\\\":{},\\\"atElem\\\":{\\\"isAtSelf\\\":false},\\\"faceElem\\\":{\\\"index\\\":0},\\\"locationElem\\\":{\\\"longitude\\\":0,\\\"latitude\\\":0},\\\"customElem\\\":{},\\\"quoteElem\\\":{},\\\"notificationElem\\\":{},\\\"messageEntityElem\\\":{},\\\"attachedInfoElem\\\":{\\\"groupHasReadInfo\\\":{\\\"hasReadCount\\\":0,\\\"groupMemberCount\\\":0},\\\"isPrivateChat\\\":false,\\\"hasReadTime\\\":0,\\\"notSenderNotificationPush\\\":false}}\"}"

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
	log.Println("send:", string(bytes))
	return bytes
}
