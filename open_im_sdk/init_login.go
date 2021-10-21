package open_im_sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"

	"time"
)

func (u *UserRelated) closeListenerCh() {
	if u.ConversationCh != nil {
		close(u.ConversationCh)
		u.ConversationCh = nil
	}
}

func (u *UserRelated) initSDK(config string, cb IMSDKListener) bool {
	u.cb = cb
	u.initListenerCh()
	sdkLog("init success, ", config)

	go doListener(u)
	return true
}

func (u *UserRelated) unInitSDK() {
	u.unInitAll()
	u.closeListenerCh()
}

func (im *IMManager) getVersion() string {
	return "v-1.0"
}

func (im *IMManager) getServerTime() int64 {
	return 0
}

func (u *UserRelated) logout(cb Base) {
	u.stateMutex.Lock()
	defer u.stateMutex.Unlock()
	u.LoginState = LogoutCmd

	err := u.closeConn()
	if err != nil {
		cb.OnError(ErrCodeInitLogin, err.Error())
		return
	}
	sdkLog("close conn ok")

	err = u.closeDB()
	if err != nil {
		cb.OnError(ErrCodeInitLogin, err.Error())
		return
	}
	sdkLog("close db ok")

	u.LoginUid = ""
	u.token = ""
	cb.OnSuccess("")
}

func (u *UserRelated) login(uid, tk string, cb Base) {
	u.token = tk
	u.LoginUid = uid

	err := u.initDBX(u.LoginUid)
	if err != nil {
		u.token = ""
		u.LoginUid = ""
		cb.OnError(ErrCodeInitLogin, err.Error())
		return
	}

	//read from db
	seq, err := u.getLocalMaxConSeq()
	if err != nil {
		sdkLog("GetLocalMaxConSeq failed ", err.Error())
		cb.OnError(ErrCodeInitLogin, err.Error())
		return
	}
	u.SetMinSeqSvr(seq)
	sdkLog("getLocalMaxConSeq ", seq)

	sdkLog("init db ok, uid: ", u.LoginUid, seq)

	err = u.syncSeq2Msg()
	if err != nil {
		sdkLog("syncSeq2Msg failed ", err.Error(), uid, tk)
		u.token = ""
		u.LoginUid = ""
		cb.OnError(ErrCodeInitLogin, err.Error())
		return
	}
	sdkLog("login sync msg ok")

	u.conn, err = u.reConn(u.conn)
	if err != nil {
		u.token = ""
		u.LoginUid = ""
		cb.OnError(ErrCodeInitLogin, err.Error())
		return
	}
	sdkLog("login ws conn ok")
	go u.forcedSynchronization() //todo:coroutine
	go u.run()
	sdkLog("ws coroutine run")
	sdkLog("login ok, ", uid, tk)
	cb.OnSuccess("")
}

func (u *UserRelated) closeConn() error {
	if u.conn != nil {
		err := u.conn.Close()
		if err != nil {
			sdkLog("close conn failed, ", err.Error())
			return err
		}
	}
	return nil
}

func (u *UserRelated) getLoginUser() string {
	if u.LoginState == LoginSuccess {
		return u.LoginUid
	} else {
		return ""
	}
}

func (im *IMManager) getLoginStatus() int {
	return im.LoginState
}

func (u *UserRelated) forcedSynchronization() {
	u.ForceSyncFriend()
	u.ForceSyncBlackList()
	u.ForceSyncFriendApplication()
	u.ForceSyncLoginUserInfo()

	u.ForceSyncMsg()

	u.ForceSyncJoinedGroup()
	u.ForceSyncGroupRequest()
	u.ForceSyncJoinedGroupMember()
	u.ForceSyncApplyGroupRequest()
	sdkLog("sync friend blacklist friendapplication userinfo  msg ok")
}

func (u *UserRelated) doWsMsg(message []byte) {
	sdkLog("doWsMsg Msg: ", string(message))
	var wsResp GeneralWsResp
	if err := json.Unmarshal(message, &wsResp); err != nil {
		sdkLog("Unmarshal failed, err: ", err.Error())
		return
	}

	switch wsResp.ReqIdentifier {
	case WSGetNewestSeq:
		u.doWSGetNewestSeq()
	case WSPullMsg:
		u.doWSPullMsg()
	case WSPushMsg:
		u.doWSPushMsg(message)
	case WSSendMsg:
		u.doWSSendMsg(wsResp)
	default:
		sdkLog("err ReqIdentifier ", wsResp.ReqIdentifier)
	}

}

func (u *UserRelated) doWSGetNewestSeq() {
	sdkLog("doWSGetNewestSeq  failed")
}

func (u *UserRelated) doWSPullMsg() {
	sdkLog("doWSPullMsg  failed")
}

func (u *UserRelated) doWSSendMsg(wsResp GeneralWsResp) {
	sdkLog("GetCh start")
	ch := u.GetCh(wsResp.MsgIncr)
	if ch == nil {
		sdkLog("failed, no chan ")
		return
	}
	sdkLog("GetCh end, ", ch)

	sdkLog("notify ch start", wsResp.OperationID)

	err := notifyCh(ch, wsResp, 1)
	if err != nil {
		sdkLog("notifyCh failed, ", err.Error(), ch, wsResp)
	}
	sdkLog("notify ch end", wsResp.OperationID)
}

func (u *UserRelated) doWSPushMsg(message []byte) {
	sdkLog("openim ws  recv push msg")
	u.doMsg(message)
}

func (u *UserRelated) delSeqMsg(beginSeq, endSeq int64) {
	u.seqMsgMutex.Lock()
	defer u.seqMsgMutex.Unlock()
	for i := beginSeq; i <= endSeq; i++ {
		sdkLog("del seqMsg, seq: ", i)
		delete(u.seqMsg, i)
	}

}

func (u *UserRelated) doMsg(message []byte) {
	sdkLog("openim ws  recv msg, do Msg: ", string(message))
	var msg Msg
	if err := json.Unmarshal(message, &msg); err != nil {
		sdkLog("Unmarshal failed, err: ", err.Error())
		return
	}
	if msg.ErrCode != 0 {
		sdkLog("errcode: ", msg.ErrCode, " errmsg: ", msg.ErrMsg)
		return
	}

	sdkLog("getConsequentLocalMaxSeq start")
	//local
	/*
		maxSeq, err := u.getConsequentLocalMaxSeq()
		if err != nil {
			sdkLog("getConsequentLocalMaxSeq failed, ", err.Error())
			return
		}
		sdkLog("getConsequentLocalMaxSeq ok, max seq: ", maxSeq)
		u.delSeqMsg(atomic.LoadInt64(&u.minSeqSvr), maxSeq)
		u.setLocalMaxConSeq(int(maxSeq))
		u.SetMinSeqSvr(int64(maxSeq))
		if maxSeq > msg.Data.Seq { // typing special handle
			sdkLog("warning seq ignore, do nothing", maxSeq, msg.Data.Seq)
		}

		if maxSeq == msg.Data.Seq {
			sdkLog("seq ignore, do nothing", maxSeq, msg.Data.Seq)
			return
		}

		//svr  17    15
		if msg.Data.Seq-maxSeq > 1 {
			//	u.pullOldMsgAndMergeNewMsg(maxSeq+1, msg.Data.Seq-1)
			u.pullBySplit(maxSeq+1, msg.Data.Seq-1)
			sdkLog("pull msg: ", maxSeq+1, msg.Data.Seq-1)
		}
	*/
	if msg.Data.SessionType == SingleChatType {
		arrMsg := ArrMsg{}
		arrMsg.SingleData = append(arrMsg.SingleData, msg.Data)

		err := u.triggerCmdNewMsgCome(arrMsg)
		sdkLog("recv push msg, trigger cmd |", msg.Data.ClientMsgID, "|", err)

		if msg.Data.ContentType > SingleTipBegin && msg.Data.ContentType < SingleTipEnd {
			u.doFriendMsg(msg.Data)
			sdkLog("doFriendMsg, ", msg.Data)
		} else if msg.Data.ContentType > GroupTipBegin && msg.Data.ContentType < GroupTipEnd {
			u.doGroupMsg(msg.Data)
			sdkLog("doGroupMsg, SingleChat ", msg.Data)
		} else {
			sdkLog("type no process, ", msg.Data)
		}

	} else if msg.Data.SessionType == GroupChatType {
		arrMsg := ArrMsg{}
		arrMsg.GroupData = append(arrMsg.GroupData, msg.Data)
		u.triggerCmdNewMsgCome(arrMsg)
		if msg.Data.ContentType > GroupTipBegin && msg.Data.ContentType < GroupTipEnd {
			u.doGroupMsg(msg.Data)
			sdkLog("doGroupMsg, ", msg.Data)
		} else {
			sdkLog("type failed, ", msg.Data)
		}
	} else {
		sdkLog("type failed, ", msg.Data)
	}
}

func (u *UserRelated) SetMinSeqSvr(minSeqSvr int64) {
	old := atomic.LoadInt64(&u.minSeqSvr)
	if minSeqSvr > old {
		atomic.StoreInt64(&u.minSeqSvr, minSeqSvr)
	}
}

func (u *UserRelated) syncSeq2Msg() error {
	//local
	sdkLog("getConsequentLocalMaxSeq start")
	maxSeq, err := u.getConsequentLocalMaxSeq()
	if err != nil {
		return err
	}
	sdkLog("getConsequentLocalMaxSeq ,max", maxSeq)

	u.delSeqMsg(atomic.LoadInt64(&u.minSeqSvr), maxSeq)
	u.setLocalMaxConSeq(int(maxSeq))
	u.SetMinSeqSvr(int64(maxSeq))
	sdkLog("setLocalMaxConSeq ", maxSeq)
	//svr
	newestSeq, err := u.getUserNewestSeq()
	if err != nil {
		sdkLog("getUserNewestSeq failed ", err.Error())
		return err
	}
	sdkLog("svr seq: ", newestSeq, ", local seq: ", maxSeq)
	if maxSeq >= newestSeq {
		sdkLog("SyncSeq2Msg LocalmaxSeq >= NewestSeq ", maxSeq, newestSeq)
		return nil
	}

	if newestSeq > maxSeq {
		sdkLog("syncSeq2Msg ", maxSeq+1, newestSeq)
		return u.pullBySplit(maxSeq+1, newestSeq) //u.pullOldMsgAndMergeNewMsg(maxSeq+1, newestSeq)
	}
	return nil
}

func (u *UserRelated) syncLoginUserInfo() error {
	userSvr, err := u.getServerUserInfo()
	if err != nil {
		return err
	}
	sdkLog("getServerUserInfo ok, user: ", *userSvr)

	userLocal, err := u.getLoginUserInfoFromLocal()
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
		err = u.replaceIntoUser(userSvr)
		if err != nil {
			u.cb.OnSelfInfoUpdated(string(bUserInfo))
		}
	}
	return nil
}

func (u *UserRelated) reConn(conn *websocket.Conn) (*websocket.Conn, error) {
	if conn != nil {
		conn.Close()
		conn = nil
	}
	u.stateMutex.Lock()
	u.LoginState = Logining
	u.stateMutex.Unlock()
	u.IMManager.cb.OnConnecting()
	url := fmt.Sprintf("%s?sendID=%s&token=%s&platformID=%d", SvrConf.IpWsAddr, u.LoginUid, u.token, SvrConf.Platform)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		//	im.LoginState = LoginFailed
		u.cb.OnConnectFailed(ErrCodeInitLogin, err.Error())
		sdkLog("websocket dial failed, ", SvrConf.IpWsAddr, u.LoginUid, u.token, SvrConf.Platform, err.Error())
		return nil, err
	}
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

	u.cb.OnConnectSuccess()
	u.stateMutex.Lock()
	u.LoginState = LoginSuccess
	u.stateMutex.Unlock()

	return conn, nil
}

func (u *UserRelated) run() {
	for {
		if u.conn == nil {
			re, err := u.reConn(nil)
			u.conn = re

			sdkLog("ws reconn ", err)
		}

		if u.conn != nil {
			sdkLog("conn ws ok, start read message")
			u.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			//			u.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			sdkLog("current read message conn:", u.conn)
			msgType, message, err := u.conn.ReadMessage()
			sdkLog("read one message")
			if err != nil {
				u.stateMutex.Lock()
				if u.LoginState == LogoutCmd {
					sdkLog("logout, ws close, return ", LogoutCmd, err)
					u.conn = nil
					u.stateMutex.Unlock()
					return
				}
				sdkLog("ws read message failed ", err.Error(), u.LoginState)

				u.stateMutex.Unlock()
				time.Sleep(time.Duration(2) * time.Second)
				u.conn, err = u.reConn(u.conn)
				sdkLog("ws reconn, ", err)

				err = u.syncSeq2Msg()
				sdkLog("sync newest msg, ", err)

			} else {
				if msgType == websocket.CloseMessage {
					sdkLog("websocket.CloseMessage, reconn to ws ")
					u.conn, _ = u.reConn(u.conn)
				} else if msgType == websocket.TextMessage {
					sdkLog("recv text message, domsg start", string(message), u.LoginState)
					//	u.doMsg(message)
					go u.doWsMsg(message)
					sdkLog("recv text message, domsg end", string(message), u.LoginState)

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

func (u *UserRelated) pullBySplit(beginSeq int64, endSeq int64) error {
	if beginSeq > endSeq {
		return nil
	}
	/*
		for k := beginSeq; k <= endSeq; k++ {
			err := u.pullOldMsgAndMergeNewMsg(k, k)
			if err != nil {
				sdkLog("pullOldMsgAndMergeNewMsg failed ", k, k, err.Error())
				return err
			}
		}

		return nil
	*/
	sdkLog("pullBySplit: begin ,end ", beginSeq, endSeq)
	var SPLIT int64 = 1000
	var bSeq, eSeq int64
	if endSeq-beginSeq > SPLIT {
		bSeq = beginSeq
		for i := 0; int64(i) < (endSeq-beginSeq)/SPLIT; i++ {
			eSeq = (int64(i)+1)*SPLIT + bSeq
			sdkLog("pull args: ", bSeq, eSeq, (int64(i)+1)*SPLIT)
			err := u.pullOldMsgAndMergeNewMsg(bSeq, eSeq)
			if err != nil {
				sdkLog("pullOldMsgAndMergeNewMsg failed ", bSeq, eSeq, err.Error())
				return err
			}

			bSeq = eSeq + 1
		}
		if bSeq <= endSeq {
			sdkLog("pull args: ", bSeq, endSeq)
			err := u.pullOldMsgAndMergeNewMsg(bSeq, endSeq)
			if err != nil {
				sdkLog("pullOldMsgAndMergeNewMsg failed ", bSeq, endSeq, err.Error())
				return err
			}

		}

	} else {
		sdkLog("pull args: ", beginSeq, endSeq)
		err := u.pullOldMsgAndMergeNewMsg(beginSeq, endSeq)
		if err != nil {
			sdkLog("pullOldMsgAndMergeNewMsg failed ", beginSeq, endSeq, err.Error())

			return err
		}

	}
	return nil
}

func (u *UserRelated) getNoInSeq(beginSeq int64, endSeq int64) (seqList []int64) {
	u.seqMsgMutex.RLock()
	defer u.seqMsgMutex.RUnlock()

	for i := beginSeq; i <= endSeq; i++ {
		_, ok := u.seqMsg[i]
		if !ok {
			seqList = append(seqList, i)
		}
	}
	return seqList
}

func (u *UserRelated) pullOldMsgAndMergeNewMsg(beginSeq int64, endSeq int64) (err error) {
	if beginSeq > endSeq {
		return nil
	}
	sdkLog("pullOldMsgAndMergeNewMsg", beginSeq, endSeq)

	seqList := u.getNoInSeq(beginSeq, endSeq)

	//todo

	sdkLog("seq list: ", seqList)
	//data := paramsPullUserMsgDataReq{SeqBegin: beginSeq, SeqEnd: endSeq}

	//	resp, err := post2Api(pullUserMsgRouter, paramsPullUserMsg{ReqIdentifier: 1002, OperationID: operationIDGenerator(), SendID: u.LoginUid, Data: data}, u.token)
	resp, err := post2Api(pullUserMsgBySeqRouter, paramsPullUserMsgBySeq{ReqIdentifier: 1002, OperationID: operationIDGenerator(), SendID: u.LoginUid, SeqList: seqList}, u.token)

	if err != nil {
		sdkLog("post2Api failed, ", pullUserMsgRouter, u.LoginUid, beginSeq, endSeq, err.Error())
		return err
	}
	sdkLog("pull ok begin end:", beginSeq, endSeq, string(resp))

	var pullMsg PullUserMsgResp
	err = json.Unmarshal(resp, &pullMsg)
	if err != nil {
		sdkLog("unmarshal failed, ", err.Error())
		return err
	}

	if pullMsg.ErrCode == 0 {
		u.seqMsgMutex.Lock()

		arrMsg := ArrMsg{}
		for i := 0; i < len(pullMsg.Data.Single); i++ {
			for j := 0; j < len(pullMsg.Data.Single[i].List); j++ {
				sdkLog("open_im pull one msg: |", pullMsg.Data.Single[i].List[j].ClientMsgID, "|")
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

				u.seqMsg[pullMsg.Data.Single[i].List[j].Seq] = singleMsg
				sdkLog("into map, seq: ", pullMsg.Data.Single[i].List[j].Seq)
				//	arrMsg.SingleData = append(arrMsg.SingleData, singleMsg)
				if pullMsg.Data.Single[i].List[j].ContentType > SingleTipBegin &&
					pullMsg.Data.Single[i].List[j].ContentType < SingleTipEnd {
					var msgRecv MsgData
					msgRecv.ContentType = pullMsg.Data.Single[i].List[j].ContentType
					msgRecv.Content = pullMsg.Data.Single[i].List[j].Content
					msgRecv.SendID = pullMsg.Data.Single[i].List[j].SendID
					msgRecv.RecvID = pullMsg.Data.Single[i].List[j].RecvID
					sdkLog("doFriendMsg ", msgRecv)
					u.doFriendMsg(msgRecv)
				}
			}
		}

		u.seqMsgMutex.Unlock()

		u.seqMsgMutex.RLock()
		for i := beginSeq; i <= endSeq; i++ {
			v, ok := u.seqMsg[i]
			if ok {
				arrMsg.SingleData = append(arrMsg.SingleData, v)
			} else {
				sdkLog("seq no in map, error, seq: ", i, u.LoginUid)
			}
		}
		u.seqMsgMutex.RUnlock()

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
					u.doGroupMsg(groupMsg)
					sdkLog("doGroupMsg ", groupMsg)
				}
			}
		}

		sdkLog("triggerCmdNewMsgCome len: ", len(arrMsg.SingleData))
		err = u.triggerCmdNewMsgCome(arrMsg)
		if err != nil {
			sdkLog("triggerCmdNewMsgCome failed, ", err.Error())
		}
	} else {
		sdkLog("pull failed, code, msg: ", pullMsg.ErrCode, pullMsg.ErrMsg)
	}
	return nil
}

func (u *UserRelated) getUserNewestSeq() (int64, error) {
	resp, err := post2Api(newestSeqRouter, paramsNewestSeqReq{ReqIdentifier: 1001, OperationID: operationIDGenerator(), SendID: u.LoginUid, MsgIncr: 1}, u.token)
	if err != nil {
		sdkLog("post2Api failed, ", newestSeqRouter, u.LoginUid, err.Error())
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

func (u *UserRelated) getServerUserInfo() (*userInfo, error) {
	var uidList []string
	uidList = append(uidList, u.LoginUid)
	resp, err := post2Api(getUserInfoRouter, paramsGetUserInfo{OperationID: operationIDGenerator(), UidList: uidList}, u.token)
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
		sdkLog("failed, no user : ", u.LoginUid)
		return nil, errors.New("no user")
	}
	return &userResp.Data[0], nil
}
func (u *UserRelated) getUserNameAndFaceUrlByUid(uid string) (faceUrl, name string, err error) {
	friendInfo, err := u.getFriendInfoByFriendUid(uid)
	if err != nil {
		return "", "", err
	}
	if friendInfo.UID == "" {
		userInfo, err := u.getUserInfoByUid(uid)
		if err != nil {
			return "", "", err
		} else {
			return userInfo.Icon, userInfo.Name, nil
		}
	} else {
		return friendInfo.Icon, friendInfo.Name, nil
	}
}
func (u *UserRelated) getUserInfoByUid(uid string) (*userInfo, error) {
	var uidList []string
	uidList = append(uidList, uid)
	resp, err := post2Api(getUserInfoRouter, paramsGetUserInfo{OperationID: operationIDGenerator(), UidList: uidList}, u.token)
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

func (u *UserRelated) doFriendMsg(msg MsgData) {
	if u.cb == nil || u.friendListener == nil {
		sdkLog("listener is null")
		return
	}

	if msg.SendID == u.LoginUid && msg.SenderPlatformID == u.SvrConf.Platform {
		sdkLog("sync msg ", msg)
		return
	}

	go func() {
		switch msg.ContentType {
		case AddFriendTip:
			u.addFriendNew(&msg) //
		case AcceptFriendApplicationTip:
			u.acceptFriendApplicationNew(&msg)
		case RefuseFriendApplicationTip:
			u.refuseFriendApplicationNew(&msg)
		case SetSelfInfoTip:
			u.setSelfInfo(&msg)
		case KickOnlineTip:
			u.kickOnline(&msg)
		default:
			sdkLog("type failed, ", msg)
		}
	}()
}

func (u *UserRelated) acceptFriendApplicationNew(msg *MsgData) {
	u.syncFriendList()
	sdkLog(msg.SendID, msg.RecvID)
	fmt.Println("sendID: ", msg.SendID, msg)

	fInfoList, err := u.getServerFriendList()
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
			u.friendListener.OnFriendListAdded(string(jData))
			u.friendListener.OnFriendApplicationListAccept(string(jData))
			return
		}
	}
}

func (u *UserRelated) refuseFriendApplicationNew(msg *MsgData) {
	sdkLog(msg.SendID, msg.RecvID)
	applyList, err := u.getServerSelfApplication()

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
			u.friendListener.OnFriendApplicationListReject(string(jData))
			return
		}
	}

}

func (u *UserRelated) addFriendNew(msg *MsgData) {
	sdkLog("addFriend start ")
	u.syncFriendApplication()

	var ui2GetUserInfo ui2ClientCommonReq
	ui2GetUserInfo.UidList = append(ui2GetUserInfo.UidList, msg.SendID)
	resp, err := post2Api(getUserInfoRouter, paramsGetUserInfo{UidList: ui2GetUserInfo.UidList, OperationID: operationIDGenerator()}, u.token)
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
	u.friendListener.OnFriendApplicationListAdded(string(jsonInfo))
}

func (im *IMManager) kickOnline(msg *MsgData) {
	im.cb.OnKickedOffline()
}

func (u *UserRelated) setSelfInfo(msg *MsgData) {
	var uidList []string
	uidList = append(uidList, msg.SendID)
	resp, err := post2Api(getUserInfoRouter, paramsGetUserInfo{OperationID: operationIDGenerator(), UidList: uidList}, u.token)
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
		sdkLog("failed, no user : ", u.LoginUid)
		return
	}

	err = u.updateFriendInfo(userResp.Data[0].Uid, userResp.Data[0].Name, userResp.Data[0].Icon, userResp.Data[0].Gender, userResp.Data[0].Mobile, userResp.Data[0].Birth, userResp.Data[0].Email, userResp.Data[0].Ex)
	if err != nil {
		sdkLog("  db change failed", err.Error())
		return
	}

	jsonInfo, err := json.Marshal(userResp.Data[0])
	if err != nil {
		sdkLog("  marshal failed", err.Error())
		return
	}

	u.friendListener.OnFriendInfoChanged(string(jsonInfo))
}
