package open_im_sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"

	"time"
)

func closeListenerCh() {
	if ConversationCh != nil {
		close(ConversationCh)
		ConversationCh = nil
	}
}

func (im *IMManager) initSDK(config string, cb IMSDKListener) bool {
	if err := json.Unmarshal([]byte(config), &SvrConf); err != nil {
		sdkLog("initSDK failed, config: ", SvrConf, err.Error())
		return false
	}

	im.cb = cb
	initListenerCh()
	initAddr()

	sdkLog("init success, ", config)

	go doListener(&ConListener)
	return true
}

func (im *IMManager) unInitSDK() {
	unInitAll()
	closeListenerCh()
}

func (im *IMManager) getVersion() string {
	return "v-1.0"
}

func (im *IMManager) getServerTime() int64 {
	return 0
}

func (im *IMManager) logout(cb Base) {
	stateMutex.Lock()
	defer stateMutex.Unlock()
	im.LoginState = LogoutCmd

	err := im.closeConn()
	if err != nil {
		cb.OnError(ErrCodeInitLogin, err.Error())
		return
	}
	sdkLog("close conn ok")

	err = closeDB()
	if err != nil {
		cb.OnError(ErrCodeInitLogin, err.Error())
		return
	}
	sdkLog("close db ok")

	LoginUid = ""
	token = ""
	cb.OnSuccess("")
}

func (im *IMManager) login(uid, tk string, cb Base) {
	token = tk
	LoginUid = uid

	err := initDBX(LoginUid)
	if err != nil {
		token = ""
		LoginUid = ""
		cb.OnError(ErrCodeInitLogin, err.Error())
		return
	}
	sdkLog("init db ok, uid: ", LoginUid)

	err = im.syncSeq2Msg()
	if err != nil {
		token = ""
		LoginUid = ""
		cb.OnError(ErrCodeInitLogin, err.Error())
		return
	}
	sdkLog("login sync msg ok")

	im.conn, err = im.reConn(im.conn)
	if err != nil {
		token = ""
		LoginUid = ""
		cb.OnError(ErrCodeInitLogin, err.Error())
		return
	}
	sdkLog("login ws conn ok")
	go im.forcedSynchronization() //todo:coroutine
	go im.run()
	sdkLog("ws coroutine run")
	sdkLog("login ok, ", uid, tk)
	cb.OnSuccess("")
}

func (im *IMManager) closeConn() error {
	if im.conn != nil {
		err := im.conn.Close()
		if err != nil {
			sdkLog("close conn failed, ", err.Error())
			return err
		}
	}
	return nil
}

func (im *IMManager) getLoginUser() string {
	if im.LoginState == LoginSuccess {
		return LoginUid
	} else {
		return ""
	}
}

func (im *IMManager) getLoginStatus() int {
	return im.LoginState
}

func (im *IMManager) forcedSynchronization() {
	ForceSyncFriend()
	ForceSyncBlackList()
	ForceSyncFriendApplication()
	ForceSyncLoginUerInfo()

	ForceSyncMsg()

	ForceSyncJoinedGroup()
	ForceSyncGroupRequest()
	ForceSyncJoinedGroupMember()
	ForceSyncApplyGroupRequest()
	sdkLog("sync friend blacklist friendapplication userinfo  msg ok")
}

func (im *IMManager) doMsg(message []byte) {
	sdkLog("ws recv msg, do Msg: ", string(message))
	var msg Msg
	if err := json.Unmarshal(message, &msg); err != nil {
		sdkLog("Unmarshal failed, err: ", err.Error())
		return
	}
	if msg.ErrCode != 0 {
		sdkLog("errcode: ", msg.ErrCode, " errmsg: ", msg.ErrMsg)
		return
	}

	//local
	maxSeq, err := getLocalMaxSeq()
	if err != nil {
		sdkLog("getLocalMaxSeq failed, ", err.Error())
		return
	}
	sdkLog("getLocalMaxSeq ok, max seq: ", maxSeq)

	if maxSeq > msg.Data.Seq { // typing special handle
		sdkLog("warning seq ignore, do nothing", maxSeq, msg.Data.Seq)
	}

	if maxSeq == msg.Data.Seq {
		sdkLog("seq ignore, do nothing", maxSeq, msg.Data.Seq)
		return
	}

	//svr  17    15
	if msg.Data.Seq-maxSeq > 1 {
		pullOldMsgAndMergeNewMsg(maxSeq+1, msg.Data.Seq-1)
		sdkLog("pull msg: ", maxSeq+1, msg.Data.Seq-1)
	}

	if msg.Data.SessionType == SingleChatType {
		arrMsg := ArrMsg{}
		arrMsg.SingleData = append(arrMsg.SingleData, msg.Data)
		triggerCmdNewMsgCome(arrMsg)

		if msg.Data.ContentType > SingleTipBegin && msg.Data.ContentType < SingleTipEnd {
			im.doFriendMsg(msg.Data)
			sdkLog("doFriendMsg, ", msg.Data)
		} else if msg.Data.ContentType > GroupTipBegin && msg.Data.ContentType < GroupTipEnd {
			groupManager.doGroupMsg(msg.Data)
			sdkLog("doGroupMsg, SingleChat ", msg.Data)
		} else {
			sdkLog("type failed, ", msg.Data)
		}

	} else if msg.Data.SessionType == GroupChatType {
		arrMsg := ArrMsg{}
		arrMsg.GroupData = append(arrMsg.GroupData, msg.Data)
		triggerCmdNewMsgCome(arrMsg)
		if msg.Data.ContentType > GroupTipBegin && msg.Data.ContentType < GroupTipEnd {
			groupManager.doGroupMsg(msg.Data)
			sdkLog("doGroupMsg, ", msg.Data)
		} else {
			sdkLog("type failed, ", msg.Data)
		}
	} else {
		sdkLog("type failed, ", msg.Data)
	}
}

func (im *IMManager) syncSeq2Msg() error {
	//local
	maxSeq, err := getLocalMaxSeq()
	if err != nil {
		return err
	}

	//svr
	newestSeq, err := getUserNewestSeq()
	if err != nil {
		return err
	}

	if maxSeq >= newestSeq {
		sdkLog("SyncSeq2Msg LocalmaxSeq >= NewestSeq ", maxSeq, newestSeq)
		return nil
	}

	if newestSeq > maxSeq {
		sdkLog("syncSeq2Msg", maxSeq+1, newestSeq)
		pullOldMsgAndMergeNewMsg(maxSeq+1, newestSeq)
	}
	return nil
}

func (im *IMManager) syncLoginUserInfo() error {
	userSvr, err := getServerUserInfo()
	if err != nil {
		return err
	}
	sdkLog("getServerUserInfo ok, user: ", *userSvr)

	userLocal, err := getLoginUserInfoFromLocal()
	if err != nil {
		return err
	}
	sdkLog("getLoginUserInfoFromLocal ok, user: ", userLocal)

	if userSvr.Uid != userLocal.Uid ||
		userSvr.Name != userLocal.Name ||
		userSvr.Icon != userLocal.Icon ||
		userSvr.Gender != userLocal.Gender ||
		userSvr.Mobile != userLocal.Mobile ||
		userSvr.Birth != userLocal.Birth ||
		userSvr.Email != userLocal.Email ||
		userSvr.Ex != userLocal.Ex {
		bUserInfo, err := json.Marshal(userSvr)
		if err != nil {
			sdkLog("marshal failed, ", err.Error())
			return err
		}
		err = replaceIntoUser(userSvr)
		if err != nil {
			im.cb.OnSelfInfoUpdated(string(bUserInfo))
		}
	}
	return nil
}

/*

func (im *IMManager) doWsCmd(cmd cmd2Value) {
	switch cmd.Cmd {
	case CmdGeyLoginUserInfo:
		im.syncLoginUserInfo()
	case CmdReLogin:
		if im.conn != nil {
			im.conn.Close()
		}
	//	im.syncSeq2Msg()
	}
}
*/

func (im *IMManager) reConn(conn *websocket.Conn) (*websocket.Conn, error) {
	if conn != nil {
		conn.Close()
		conn = nil
	}
	stateMutex.Lock()
	im.LoginState = Logining
	stateMutex.Unlock()
	SdkInitManager.cb.OnConnecting()
	url := fmt.Sprintf("%s?sendID=%s&token=%s&platformID=%d", SvrConf.IpWsAddr, LoginUid, token, SvrConf.Platform)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		//	im.LoginState = LoginFailed
		SdkInitManager.cb.OnConnectFailed(ErrCodeInitLogin, err.Error())
		sdkLog("websocket dial failed, ", SvrConf.IpWsAddr, LoginUid, token, SvrConf.Platform, err.Error())
		return nil, err
	}
	SdkInitManager.cb.OnConnectSuccess()
	stateMutex.Lock()
	im.LoginState = LoginSuccess
	stateMutex.Unlock()

	return conn, nil
}

/*
func (im *IMManager) runCmd() {
	for {
		select {
		case v := <-im.ch:
			if v.Cmd == CmdUnInit {
				if im.conn != nil {
					im.conn.Close()
					im.conn = nil
					atomic.SwapInt32(&WsState, 100)
				}
				return
			} else {

				im.doWsCmd(v)
			}
		}
	}
}
*/

func (im *IMManager) run() {

	for {
		if im.conn == nil {
			re, err := im.reConn(nil)
			im.conn = re

			sdkLog("ws reconn ", err)
		}

		if im.conn != nil {
			sdkLog("conn ws ok, start read message")
			im.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			sdkLog("current read message conn:", im.conn)
			msgType, message, err := im.conn.ReadMessage()
			sdkLog("read one message")
			if err != nil {
				stateMutex.Lock()
				if im.LoginState == LogoutCmd {
					sdkLog("logout, ws close, return ", LogoutCmd, err)
					im.conn = nil
					stateMutex.Unlock()
					return
				}
				sdkLog("ws read message failed ", err.Error())
				/*
					if netErr, ok := err.(net.Error); ok {
						if netErr.Timeout() {

							sdkLog("ws read message failed ", err.Error())
							stateMutex.Unlock()
							continue
						}
					}
				*/
				stateMutex.Unlock()
				time.Sleep(time.Duration(2) * time.Second)
				im.conn, err = im.reConn(im.conn)
				sdkLog("ws reconn, ", err)

				err = im.syncSeq2Msg()
				sdkLog("sync newest msg, ", err)

			} else {
				if msgType == websocket.CloseMessage {
					sdkLog("websocket.CloseMessage, reconn to ws ")
					im.conn, _ = im.reConn(im.conn)
				} else if msgType == websocket.TextMessage {
					sdkLog("recv text message, domsg ", string(message))
					im.doMsg(message)
				} else {
					sdkLog("recv msg: type ", msgType)
				}
			}
		} else {
			sdkLog("ws failed, sleep 2s, reconn... ")
			time.Sleep(time.Duration(2) * time.Second)
		}
	}
}

func pullOldMsgAndMergeNewMsg(beginSeq int64, endSeq int64) (err error) {
	sdkLog("pullOldMsgAndMergeNewMsg", beginSeq, endSeq)
	data := paramsPullUserMsgDataReq{SeqBegin: beginSeq, SeqEnd: endSeq}

	resp, err := post2Api(pullUserMsgRouter, paramsPullUserMsg{ReqIdentifier: 1002, OperationID: operationIDGenerator(), SendID: LoginUid, Data: data}, token)
	if err != nil {
		sdkLog("post2Api failed, ", pullUserMsgRouter, LoginUid, beginSeq, endSeq, err.Error())
		return err
	}
	sdkLog("pull ok begin end:", beginSeq, endSeq)

	var pullMsg PullUserMsgResp
	err = json.Unmarshal(resp, &pullMsg)
	if err != nil {
		sdkLog("unmarshal failed, ", err.Error())
		return err
	}

	if pullMsg.ErrCode == 0 {
		arrMsg := ArrMsg{}
		for i := 0; i < len(pullMsg.Data.Single); i++ {
			for j := 0; j < len(pullMsg.Data.Single[i].List); j++ {
				sdkLog("pull msg: ", pullMsg.Data.Single[i].List[j])
				singleMsg := MsgData{
					SendID:           pullMsg.Data.Single[i].List[j].SendID,
					RecvID:           pullMsg.Data.Single[i].List[j].RecvID,
					SessionType:      SingleChatType,
					MsgFrom:          pullMsg.Data.Single[i].List[j].MsgFrom,
					ContentType:      pullMsg.Data.Single[i].List[j].ContentType,
					ServerMsgID:      pullMsg.Data.Single[i].List[j].ServerMsgID,
					Content:          pullMsg.Data.Single[i].List[j].Content,
					SendTime:         pullMsg.Data.Single[i].List[j].SendTime,
					Seq:              pullMsg.Data.Single[i].List[j].Seq,
					SenderNickName:   pullMsg.Data.Single[i].List[j].SenderNickName,
					SenderFaceURL:    pullMsg.Data.Single[i].List[j].SenderFaceURL,
					ClientMsgID:      pullMsg.Data.Single[i].List[j].ClientMsgID,
					SenderPlatformID: pullMsg.Data.Single[i].List[j].SenderPlatformID,
				}
				arrMsg.SingleData = append(arrMsg.SingleData, singleMsg)
				if pullMsg.Data.Single[i].List[j].ContentType > SingleTipBegin &&
					pullMsg.Data.Single[i].List[j].ContentType < SingleTipEnd {
					var msgRecv MsgData
					msgRecv.ContentType = pullMsg.Data.Single[i].List[j].ContentType
					msgRecv.Content = pullMsg.Data.Single[i].List[j].Content
					msgRecv.SendID = pullMsg.Data.Single[i].List[j].SendID
					msgRecv.RecvID = pullMsg.Data.Single[i].List[j].RecvID
					sdkLog("doFriendMsg ", msgRecv)
					SdkInitManager.doFriendMsg(msgRecv)
				}
			}
		}

		for i := 0; i < len(pullMsg.Data.Group); i++ {
			for j := 0; j < len(pullMsg.Data.Group[i].List); j++ {
				groupMsg := MsgData{
					SendID:           pullMsg.Data.Group[i].List[j].SendID,
					RecvID:           pullMsg.Data.Group[i].List[j].RecvID,
					SessionType:      GroupChatType,
					MsgFrom:          pullMsg.Data.Group[i].List[j].MsgFrom,
					ContentType:      pullMsg.Data.Group[i].List[j].ContentType,
					ServerMsgID:      pullMsg.Data.Group[i].List[j].ServerMsgID,
					Content:          pullMsg.Data.Group[i].List[j].Content,
					SendTime:         pullMsg.Data.Group[i].List[j].SendTime,
					Seq:              pullMsg.Data.Group[i].List[j].Seq,
					SenderNickName:   pullMsg.Data.Group[i].List[j].SenderNickName,
					SenderFaceURL:    pullMsg.Data.Group[i].List[j].SenderFaceURL,
					ClientMsgID:      pullMsg.Data.Group[i].List[j].ClientMsgID,
					SenderPlatformID: pullMsg.Data.Group[i].List[j].SenderPlatformID,
				}
				arrMsg.GroupData = append(arrMsg.GroupData, groupMsg)

				ctype := pullMsg.Data.Group[i].List[j].ContentType
				if ctype > GroupTipBegin && ctype < GroupTipEnd {
					groupManager.doGroupMsg(groupMsg)
					sdkLog("doGroupMsg ", groupMsg)
				}
			}
		}

		sdkLog("triggerCmdNewMsgCome len: ", len(arrMsg.SingleData))
		err = triggerCmdNewMsgCome(arrMsg)
		if err != nil {
			sdkLog("triggerCmdNewMsgCome failed, ", err.Error())
		}
	} else {
		sdkLog("pull failed, code, msg: ", pullMsg.ErrCode, pullMsg.ErrMsg)
	}
	return nil
}

func getUserNewestSeq() (int64, error) {
	resp, err := post2Api(newestSeqRouter, paramsNewestSeqReq{ReqIdentifier: 1001, OperationID: operationIDGenerator(), SendID: LoginUid, MsgIncr: 1}, token)
	if err != nil {
		sdkLog("post2Api failed, ", newestSeqRouter, LoginUid, err.Error())
		return 0, err
	}
	var seqResp paramsNewestSeqResp
	err = json.Unmarshal(resp, &seqResp)
	if err != nil {
		sdkLog("UnMarshal failed, ", err.Error())
		return 0, err
	}

	if seqResp.ErrCode != 0 {
		sdkLog("errcode: ", seqResp.ErrCode, "errmsg: ", seqResp.ErrMsg)
		return 0, errors.New(seqResp.ErrMsg)
	}

	return seqResp.Data.Seq, nil
}

func getServerUserInfo() (*userInfo, error) {
	var uidList []string
	uidList = append(uidList, LoginUid)
	resp, err := post2Api(getUserInfoRouter, paramsGetUserInfo{OperationID: operationIDGenerator(), UidList: uidList}, token)
	if err != nil {
		sdkLog("post2Api failed, ", getUserInfoRouter, uidList, err.Error())
		return nil, err
	}
	var userResp getUserInfoResp
	err = json.Unmarshal(resp, &userResp)
	if err != nil {
		sdkLog("Unmarshal failed, ", resp, err.Error())
		return nil, err
	}

	if userResp.ErrCode != 0 {
		sdkLog("errcode: ", userResp.ErrCode, "errmsg:", userResp.ErrMsg)
		return nil, errors.New(userResp.ErrMsg)
	}

	if len(userResp.Data) == 0 {
		sdkLog("failed, no user : ", LoginUid)
		return nil, errors.New("no user")
	}
	return &userResp.Data[0], nil
}
func getUserNameAndFaceUrlByUid(uid string) (faceUrl, name string, err error) {
	friendInfo, err := getFriendInfoByFriendUid(uid)
	if err != nil {
		return "", "", err
	}
	if friendInfo.UID == "" {
		userInfo, err := getUserInfoByUid(uid)
		if err != nil {
			return "", "", err
		} else {
			return userInfo.Icon, userInfo.Name, nil
		}
	} else {
		return friendInfo.Icon, friendInfo.Name, nil
	}
}
func getUserInfoByUid(uid string) (*userInfo, error) {
	var uidList []string
	uidList = append(uidList, uid)
	resp, err := post2Api(getUserInfoRouter, paramsGetUserInfo{OperationID: operationIDGenerator(), UidList: uidList}, token)
	if err != nil {
		sdkLog("post2Api failed, ", getUserInfoRouter, uidList, err.Error())
		return nil, err
	}
	sdkLog("post api: ", getUserInfoRouter, paramsGetUserInfo{OperationID: operationIDGenerator(), UidList: uidList}, "uid ", uid)
	var userResp getUserInfoResp
	err = json.Unmarshal(resp, &userResp)
	if err != nil {
		sdkLog("Unmarshal failed, ", resp, err.Error())
		return nil, err
	}

	if userResp.ErrCode != 0 {
		sdkLog("errcode: ", userResp.ErrCode, "errmsg:", userResp.ErrMsg)
		return nil, errors.New(userResp.ErrMsg)
	}

	if len(userResp.Data) == 0 {
		sdkLog("failed, no user :", uid)
		return nil, errors.New("no user")
	}
	return &userResp.Data[0], nil
}

func (im *IMManager) doFriendMsg(msg MsgData) {
	if im.cb == nil || FriendObj.friendListener == nil {
		sdkLog("listener is null")
		return
	}

	go func() {
		switch msg.ContentType {
		case AddFriendTip:
			im.addFriend(&msg)
		case AcceptFriendApplicationTip:
			im.acceptFriendApplication(&msg)
		case RefuseFriendApplicationTip:
			im.refuseFriendApplication(&msg)
		case SetSelfInfoTip:
			im.setSelfInfo(&msg)
		case KickOnlineTip:
			im.kickOnline(&msg)
		default:
			sdkLog("type failed, ", msg)
		}
	}()
}

func (im *IMManager) acceptFriendApplication(msg *MsgData) {
	FriendObj.syncFriendList()
	sdkLog(msg.SendID, msg.RecvID)
	fmt.Println("sendID: ", msg.SendID, msg)

	fInfoList, err := FriendObj.getServerFriendList()
	if err != nil {
		return
	}
	for _, fInfo := range fInfoList {
		if fInfo.UID == msg.SendID {
			jData, err := json.Marshal(fInfo)
			if err != nil {
				sdkLog("err: ", err.Error())
				return
			}
			FriendObj.friendListener.OnFriendListAdded(string(jData))
			FriendObj.friendListener.OnFriendApplicationListAccept(string(jData))
			return
		}
	}
}

func (im *IMManager) refuseFriendApplication(msg *MsgData) {
	sdkLog(msg.SendID, msg.RecvID)
	applyList, err := FriendObj.getServerSelfApplication()

	if err != nil {
		return
	}
	for _, v := range applyList {
		if v.Uid == msg.SendID {
			jData, err := json.Marshal(v)
			if err != nil {
				sdkLog("err: ", err.Error())
				return
			}
			FriendObj.friendListener.OnFriendApplicationListReject(string(jData))
			return
		}
	}

}

func (im *IMManager) addFriend(msg *MsgData) {
	sdkLog("addFriend start ")
	FriendObj.syncFriendApplication()

	var ui2GetUserInfo ui2ClientCommonReq
	ui2GetUserInfo.UidList = append(ui2GetUserInfo.UidList, msg.SendID)
	resp, err := post2Api(getUserInfoRouter, paramsGetUserInfo{UidList: ui2GetUserInfo.UidList, OperationID: operationIDGenerator()}, token)
	if err != nil {
		sdkLog("getUserInfo failed", err)
		return
	}
	var vgetUserInfoResp getUserInfoResp
	err = json.Unmarshal(resp, &vgetUserInfoResp)
	if err != nil {
		sdkLog("Unmarshal failed, ", err.Error())
		return
	}
	if vgetUserInfoResp.ErrCode != 0 {
		sdkLog(vgetUserInfoResp.ErrCode, vgetUserInfoResp.ErrMsg)
		return
	}
	if len(vgetUserInfoResp.Data) == 0 {
		sdkLog(vgetUserInfoResp.ErrCode, vgetUserInfoResp.ErrMsg, msg)
		return
	}
	var appUserNode applyUserInfo
	appUserNode.Uid = vgetUserInfoResp.Data[0].Uid
	appUserNode.Name = vgetUserInfoResp.Data[0].Name
	appUserNode.Icon = vgetUserInfoResp.Data[0].Icon
	appUserNode.Gender = vgetUserInfoResp.Data[0].Gender
	appUserNode.Mobile = vgetUserInfoResp.Data[0].Mobile
	appUserNode.Birth = vgetUserInfoResp.Data[0].Birth
	appUserNode.Email = vgetUserInfoResp.Data[0].Email
	appUserNode.Ex = vgetUserInfoResp.Data[0].Ex
	appUserNode.Flag = 0

	jsonInfo, err := json.Marshal(appUserNode)
	if err != nil {
		sdkLog("  marshal failed", err.Error())
		return
	}
	FriendObj.friendListener.OnFriendApplicationListAdded(string(jsonInfo))
}

func (im *IMManager) kickOnline(msg *MsgData) {
	im.cb.OnKickedOffline()
}

func (im *IMManager) setSelfInfo(msg *MsgData) {
	var uidList []string
	uidList = append(uidList, msg.SendID)
	resp, err := post2Api(getUserInfoRouter, paramsGetUserInfo{OperationID: operationIDGenerator(), UidList: uidList}, token)
	if err != nil {
		sdkLog("post2Api failed, ", getUserInfoRouter, uidList, err.Error())
		return
	}
	var userResp getUserInfoResp
	err = json.Unmarshal(resp, &userResp)
	if err != nil {
		sdkLog("Unmarshal failed, ", resp, err.Error())
		return
	}

	if userResp.ErrCode != 0 {
		sdkLog("errcode: ", userResp.ErrCode, "errmsg:", userResp.ErrMsg)
		return
	}

	if len(userResp.Data) == 0 {
		sdkLog("failed, no user : ", LoginUid)
		return
	}

	err = updateFriendInfo(userResp.Data[0].Uid, userResp.Data[0].Name, userResp.Data[0].Icon, userResp.Data[0].Gender, userResp.Data[0].Mobile, userResp.Data[0].Birth, userResp.Data[0].Email, userResp.Data[0].Ex)
	if err != nil {
		sdkLog("  db change failed", err.Error())
		return
	}

	jsonInfo, err := json.Marshal(userResp.Data[0])
	if err != nil {
		sdkLog("  marshal failed", err.Error())
		return
	}

	FriendObj.friendListener.OnFriendInfoChanged(string(jsonInfo))
}
