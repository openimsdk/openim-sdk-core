package open_im_sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"

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
	if cb == nil {
		sdkLog("callback == nil")
		return false
	}
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
	sdkLog("closeConn ok")

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
	sdkLog("login start, ", uid, tk)
	u.token = tk
	u.LoginUid = uid

	err := u.initDBX(u.LoginUid)
	if err != nil {
		u.token = ""
		u.LoginUid = ""
		cb.OnError(ErrCodeInitLogin, err.Error())
		sdkLog("initDBX failed, ", err.Error())
		return
	}
	sdkLog("initDBX ok ", uid)

	u.conn, err = u.reConn(nil)
	if err != nil {
		u.token = ""
		u.LoginUid = ""
		cb.OnError(ErrCodeInitLogin, err.Error())
		sdkLog("reConn failed ", err.Error())
		return
	}
	sdkLog("ws conn ok ")

	seq, err := u.getLocalMaxConSeqFromDB()
	if err != nil {
		sdkLog("getLocalMaxConSeqFromDB failed ", err.Error())
		cb.OnError(ErrCodeInitLogin, err.Error())
		return
	}
	u.SetMinSeqSvr(seq)
	sdkLog("getLocalMaxConSeqFromDB SetMinSeqSvr ok ", seq)
	go u.run()

	//err = u.syncSeq2Msg()
	//if err != nil {
	//	sdkLog("syncSeq2Msg failed ", err.Error(), uid, tk)
	//	u.token = ""
	//	u.LoginUid = ""
	//	cb.OnError(ErrCodeInitLogin, err.Error())
	//	return
	//}
	//sdkLog("syncSeq2Msg ok ", uid, tk)

	sdkLog("ws, forcedSynchronization heartbeat coroutine timedCloseDB run ...")
	go u.forcedSynchronization()
	go u.heartbeat()
	go u.timedCloseDB()
	cb.OnSuccess("")
	sdkLog("login end, ", uid, tk)
}

func (u *UserRelated) timedCloseDB() {
	timeTicker := time.NewTicker(time.Second * 300)
	for {
		<-timeTicker.C
		sdkLog("closeDBSetNil begin")
		u.closeDBSetNil()
		sdkLog("closeDBSetNil end")

		u.stateMutex.Lock()
		if u.LoginState == LogoutCmd {
			sdkLog("logout timedCloseDB return", LogoutCmd)
			u.stateMutex.Unlock()
			return
		}
		u.stateMutex.Unlock()
	}
}

func (u *UserRelated) closeConn() error {
	LogBegin()
	if u.conn != nil {
		err := u.conn.Close()
		if err != nil {
			LogFReturn(err.Error())
			return err
		}
	}
	LogSReturn(nil)
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
	LogBegin()
	u.ForceSyncFriend()
	u.ForceSyncBlackList()
	u.ForceSyncFriendApplication()
	u.ForceSyncLoginUserInfo()

	//u.ForceSyncMsg()

	u.ForceSyncJoinedGroup()
	u.ForceSyncGroupRequest()
	u.ForceSyncJoinedGroupMember()
	u.ForceSyncApplyGroupRequest()
	LogSReturn()
}

func (u *UserRelated) doWsMsg(message []byte) {
	LogBegin()
	LogBegin("decodeBinaryWs")
	wsResp, err := u.decodeBinaryWs(message)
	LogEnd("decodeBinaryWs ", wsResp.OperationID, wsResp.ReqIdentifier)
	if err != nil {
		LogFReturn(err.Error())
		return
	}

	switch wsResp.ReqIdentifier {
	case WSGetNewestSeq:
		u.doWSGetNewestSeq(*wsResp)
	case WSPullMsgBySeqList:
		u.doWSPullMsg(*wsResp)
	case WSPushMsg:
		u.doWSPushMsg(*wsResp)
	case WSSendMsg:
		u.doWSSendMsg(*wsResp)
	default:
		LogFReturn("type failed, ", wsResp.ReqIdentifier, wsResp.OperationID)
		return
	}

	LogSReturn()
	return
}

func (u *UserRelated) notifyResp(wsResp GeneralWsResp) {
	LogBegin(wsResp.OperationID)
	u.wsMutex.Lock()
	defer u.wsMutex.Unlock()

	ch := u.GetCh(wsResp.MsgIncr)
	if ch == nil {
		sdkLog("failed, no chan ", wsResp.MsgIncr, wsResp.OperationID)
		return
	}
	sdkLog("GetCh end, ", ch)

	sdkLog("notify ch start", wsResp.OperationID)

	err := notifyCh(ch, wsResp, 1)
	if err != nil {
		sdkLog("notifyCh failed, ", err.Error(), ch, wsResp)
	}
	sdkLog("notify ch end", wsResp.OperationID)
	LogSReturn(nil)
}

func (u *UserRelated) doWSGetNewestSeq(wsResp GeneralWsResp) {
	LogBegin(wsResp.OperationID)
	u.notifyResp(wsResp)
	LogSReturn(wsResp.OperationID)
}

func (u *UserRelated) doWSPullMsg(wsResp GeneralWsResp) {
	LogBegin(wsResp.OperationID)
	u.notifyResp(wsResp)
	LogSReturn(wsResp.OperationID)
}

func (u *UserRelated) doWSSendMsg(wsResp GeneralWsResp) {
	LogBegin(wsResp.OperationID)
	u.notifyResp(wsResp)
	LogSReturn(wsResp.OperationID)
}

func (u *UserRelated) doWSPushMsg(wsResp GeneralWsResp) {
	LogBegin()
	u.doMsg(wsResp)
	LogSReturn()
}

func (u *UserRelated) delSeqMsg(beginSeq, endSeq int64) {
	sdkLog("delSeqMsg, seq begin: ", beginSeq, " end: ", endSeq)
	u.seqMsgMutex.Lock()
	defer u.seqMsgMutex.Unlock()
	for i := beginSeq; i <= endSeq; i++ {
		delete(u.seqMsg, i)
	}
}

func (u *UserRelated) doMsg(wsResp GeneralWsResp) {
	LogBegin(wsResp.OperationID)
	var msg MsgData

	if wsResp.ErrCode != 0 {
		sdkLog("errcode: ", wsResp.ErrCode, " errmsg: ", wsResp.ErrMsg)
		LogFReturn()
		return
	}

	err := proto.Unmarshal(wsResp.Data, &msg)
	if err != nil {
		sdkLog("Unmarshal failed", err.Error())
		LogFReturn()
		return
	}

	if msg.SessionType == SingleChatType {
		arrMsg := ArrMsg{}
		arrMsg.SingleData = append(arrMsg.SingleData, msg)

		err := u.triggerCmdNewMsgCome(arrMsg)
		sdkLog("recv push msg, trigger cmd |", msg.ClientMsgID, "|", err)

		if msg.ContentType > SingleTipBegin && msg.ContentType < SingleTipEnd {
			u.doFriendMsg(msg)
			sdkLog("doFriendMsg, ", msg)
		} else if msg.ContentType > GroupTipBegin && msg.ContentType < GroupTipEnd {
			u.doGroupMsg(msg)
			sdkLog("doGroupMsg, SingleChat ", msg)
		} else {
			sdkLog("type no process, ", msg)
		}

	} else if msg.SessionType == GroupChatType {
		arrMsg := ArrMsg{}
		arrMsg.GroupData = append(arrMsg.GroupData, msg)
		u.triggerCmdNewMsgCome(arrMsg)
		if msg.ContentType > GroupTipBegin && msg.ContentType < GroupTipEnd {
			u.doGroupMsg(msg)
			sdkLog("doGroupMsg, ", msg)
		} else {
			sdkLog("type failed, ", msg)
		}
	} else {
		sdkLog("type failed, ", msg)
	}

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

}

func (u *UserRelated) GetMinSeqSvr() int64 {
	u.minSeqSvrRWMutex.RLock()
	min := u.minSeqSvr
	u.minSeqSvrRWMutex.RUnlock()
	return min
}

func (u *UserRelated) SetMinSeqSvr(minSeqSvr int64) {

	u.minSeqSvrRWMutex.Lock()
	if minSeqSvr > u.minSeqSvr {
		u.minSeqSvr = minSeqSvr
	}
	u.minSeqSvrRWMutex.Unlock()

	/*
		old := atomic.LoadInt64(&u.minSeqSvr)
		if minSeqSvr > old {
			atomic.StoreInt64(&u.minSeqSvr, minSeqSvr)
		}
	*/
}

func (u *UserRelated) syncMsg2ServerMaxSeq(serverMaxSeq int64) error {
	LogBegin(serverMaxSeq)
	LogBegin("getConsequentLocalMaxSeq")
	maxSeq, err := u.getConsequentLocalMaxSeq() //local
	LogEnd("getConsequentLocalMaxSeq", maxSeq, err)
	if err != nil {
		LogFReturn(err.Error())
		return err
	}
	minSeqSvr := u.GetMinSeqSvr()
	//	LogBegin("delSeqMsg setLocalMaxConSeq SetMinSeqSvr", atomic.LoadInt64(&u.minSeqSvr), maxSeq)
	LogBegin("delSeqMsg setLocalMaxConSeq SetMinSeqSvr", minSeqSvr, maxSeq)
	//u.delSeqMsg(atomic.LoadInt64(&u.minSeqSvr), maxSeq)
	u.delSeqMsg(minSeqSvr, maxSeq)
	u.setLocalMaxConSeq(int(maxSeq))
	u.SetMinSeqSvr(int64(maxSeq))
	LogEnd("elSeqMsg setLocalMaxConSeq SetMinSeqSvr")

	if maxSeq >= serverMaxSeq {
		sdkLog("don't sync,  LocalmaxSeq >= NewestSeq ", maxSeq, serverMaxSeq)
		LogSReturn(nil)
		return nil
	}

	LogBegin("pullBySplit", maxSeq+1, serverMaxSeq)
	err = u.pullBySplit(maxSeq+1, serverMaxSeq)
	LogEnd("pullBySplit", maxSeq+1, serverMaxSeq)
	if err != nil {
		LogFReturn(err.Error())
		return err
	} else {
		LogSReturn(nil)
		return nil
	}
}

func (u *UserRelated) syncSeq2Msg() error {
	svrMaxSeq, err := u.getUserNewestSeq()
	if err != nil {
		sdkLog("getUserNewestSeq failed ", err.Error())
		return err
	}
	return u.syncMsg2ServerMaxSeq(svrMaxSeq)
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
	LogBegin(conn)
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
		u.cb.OnConnectFailed(ErrCodeInitLogin, err.Error())
		LogFReturn(nil, err.Error(), url)
		return nil, err
	}
	sdkLog("ws connect ok, ", url)
	u.cb.OnConnectSuccess()
	u.stateMutex.Lock()
	u.LoginState = LoginSuccess
	u.stateMutex.Unlock()
	LogSReturn(conn, nil)
	return conn, nil
}

func (u *UserRelated) heartbeat() {
	for {
		u.stateMutex.Lock()
		if u.LoginState == LogoutCmd {
			sdkLog("logout, ws close, heartbeat return ", LogoutCmd)
			u.conn = nil
			u.stateMutex.Unlock()
			return
		}
		u.stateMutex.Unlock()

		LogBegin()
		LogBegin("AddCh")
		msgIncr, ch := u.AddCh()
		LogEnd("AddCh")

		var wsReq GeneralWsReq
		wsReq.ReqIdentifier = WSGetNewestSeq
		wsReq.OperationID = operationIDGenerator()
		wsReq.SendID = u.LoginUid
		wsReq.MsgIncr = msgIncr
		var connSend *websocket.Conn
		LogBegin("WriteMsg", wsReq.OperationID, wsReq.MsgIncr)
		err, connSend := u.WriteMsg(wsReq)
		LogEnd("WriteMsg", wsReq.OperationID, wsReq.MsgIncr)
		if err != nil {
			sdkLog("WriteMsg failed ", err.Error(), msgIncr, wsReq.OperationID)
			LogBegin("closeConn DelCh", msgIncr, wsReq.OperationID)
			u.closeConn()
			u.DelCh(msgIncr)
			LogEnd("closeConn DelCh continue", wsReq.OperationID)
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}

		timeout := 30
		breakFlag := 0
		for {
			if breakFlag == 1 {
				sdkLog("break ", wsReq.OperationID)
				break
			}
			select {
			case r := <-ch:
				sdkLog("ws ch recvMsg success: ", wsReq.OperationID, "seq cache map size: ", len(u.seqMsg))
				if r.ErrCode != 0 {
					sdkLog("heartbeat response faield ", r.ErrCode, r.ErrMsg, wsReq.OperationID)
					LogBegin("closeConn DelCh", msgIncr, wsReq.OperationID)
					u.closeConn()
					//u.DelCh(msgIncr)
					LogEnd("closeConn DelCh continue", wsReq.OperationID)

				} else {
					sdkLog("heartbeat response success ", wsReq.OperationID)
					var wsSeqResp GetMaxAndMinSeqResp
					err = proto.Unmarshal(r.Data, &wsSeqResp)
					if err != nil {
						sdkLog("Unmarshal failed, ", err.Error(), wsReq.OperationID)
						LogBegin("closeConn DelCh", msgIncr, wsReq.OperationID)
						u.closeConn()
						//	u.DelCh(msgIncr)
						LogEnd("closeConn DelCh continue")
					} else {
						//if wsSeqResp.MinSeq > atomic.LoadInt64(&u.minSeqSvr) {
						minSeqSvr := u.GetMinSeqSvr()
						if wsSeqResp.MinSeq > minSeqSvr {
							LogBegin("setLocalMaxConSeq SetMinSeqSvr ", wsSeqResp.MinSeq, minSeqSvr)
							u.SetMinSeqSvr(wsSeqResp.MinSeq)
							LogEnd("setLocalMaxConSeq SetMinSeqSvr ")
						}
						LogBegin("syncMsg2ServerMaxSeq", wsSeqResp.MaxSeq, wsSeqResp.MinSeq, wsReq.OperationID)
						u.syncMsg2ServerMaxSeq(wsSeqResp.MaxSeq)
						LogEnd("syncMsg2ServerMaxSeq")
					}
				}
				breakFlag = 1

			case <-time.After(time.Second * time.Duration(timeout)):
				var flag bool
				sdkLog("ws ch recvMsg err: ", wsReq.OperationID)
				if connSend != u.conn {
					sdkLog("old conn != current conn  ", connSend, u.conn)
					flag = false // error
				} else {
					flag = false //error
					for tr := 0; tr < 3; tr++ {
						err = u.sendPingMsg()
						if err != nil {
							sdkLog("sendPingMsg failed ", wsReq.OperationID, err.Error(), tr)
							time.Sleep(time.Duration(30) * time.Second)
						} else {
							sdkLog("sendPingMsg ok ", wsReq.OperationID)
							flag = true //wait continue
							break
						}
					}
				}
				if flag == false {
					sdkLog("ws ch recvMsg timeout ", timeout, "s ", wsReq.OperationID)
					LogBegin("closeConn", wsReq.OperationID)
					u.closeConn()
					LogEnd("closeConn", wsReq.OperationID)
					breakFlag = 1
					break
				} else {
					sdkLog("wait resp continue", wsReq.OperationID)
					breakFlag = 0
					continue
				}
			}
		}

		u.DelCh(msgIncr)
		LogEnd("DelCh", wsReq.OperationID)
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func (u *UserRelated) run() {
	for {
		LogStart()
		if u.conn == nil {
			LogBegin("reConn", nil)
			re, _ := u.reConn(nil)
			LogEnd("reConn", re)
			u.conn = re
		}
		if u.conn != nil {
			msgType, message, err := u.conn.ReadMessage()
			sdkLog("ReadMessage message ", msgType, err)
			if err != nil {
				u.stateMutex.Lock()
				sdkLog("ws read message failed ", err.Error(), u.LoginState)
				if u.LoginState == LogoutCmd {
					sdkLog("logout, ws close, return ", LogoutCmd, err)
					u.conn = nil
					u.stateMutex.Unlock()
					return
				}
				u.stateMutex.Unlock()
				time.Sleep(time.Duration(5) * time.Second)
				sdkLog("ws  ReadMessage failed, sleep 5s, reconn, ", err)
				LogBegin("reConn", u.conn)
				u.conn, err = u.reConn(u.conn)
				LogEnd("reConn", u.conn)
			} else {
				if msgType == websocket.CloseMessage {
					u.conn, _ = u.reConn(u.conn)
				} else if msgType == websocket.TextMessage {
					sdkLog("type failed, recv websocket.TextMessage ", string(message))
				} else if msgType == websocket.BinaryMessage {
					go u.doWsMsg(message)
				} else {
					sdkLog("recv other msg: type ", msgType)
				}
			}
		} else {
			sdkLog("ws failed, sleep 5s, reconn... ")
			time.Sleep(time.Duration(5) * time.Second)
		}
	}
}

func (u *UserRelated) pullBySplit(beginSeq int64, endSeq int64) error {
	LogBegin(beginSeq, endSeq)
	if beginSeq > endSeq {
		LogFReturn("beginSeq > endSeq")
		return nil
	}
	var SPLIT int64 = 100
	var bSeq, eSeq int64
	if endSeq-beginSeq > SPLIT {
		bSeq = beginSeq
		//1, 118 117/10 = 11  i: 0- 10  1-> 11 12->22 23->33  34->44   111->121
		for i := 0; int64(i) < (endSeq-beginSeq)/SPLIT; i++ {
			eSeq = bSeq + SPLIT - 1
			sdkLog("pull args: ", i, bSeq, eSeq, (endSeq-beginSeq)/SPLIT)
			err := u.pullOldMsgAndMergeNewMsg(bSeq, eSeq)
			if err != nil {
				LogFReturn(err.Error())
				return err
			}
			bSeq = eSeq + 1
		}
		if bSeq <= endSeq {
			sdkLog("pull remainder args: ", bSeq, endSeq)
			err := u.pullOldMsgAndMergeNewMsg(bSeq, endSeq)
			if err != nil {
				LogFReturn(err.Error())
				return err
			}
		}
	} else {
		err := u.pullOldMsgAndMergeNewMsg(beginSeq, endSeq)
		if err != nil {
			LogFReturn(err.Error())
			return err
		}
	}
	return nil
}

func (u *UserRelated) getNoInSeq(beginSeq int64, endSeq int64) (seqList []int64) {
	LogBegin(beginSeq, endSeq)
	u.seqMsgMutex.RLock()
	defer u.seqMsgMutex.RUnlock()

	for i := beginSeq; i <= endSeq; i++ {
		_, ok := u.seqMsg[i]
		if !ok {
			seqList = append(seqList, i)
		}
	}
	LogSReturn(seqList)
	return seqList
}

func (u *UserRelated) pullOldMsgAndMergeNewMsgByWs(beginSeq int64, endSeq int64) (err error) {
	LogBegin(beginSeq, endSeq)
	if beginSeq > endSeq {
		LogSReturn(nil)
		return nil
	}
	LogBegin("AddCh")
	msgIncr, ch := u.AddCh()
	LogEnd("AddCh")

	var wsReq GeneralWsReq
	wsReq.ReqIdentifier = WSPullMsgBySeqList
	wsReq.OperationID = operationIDGenerator()
	wsReq.SendID = u.LoginUid
	//wsReq.Token = u.token
	wsReq.MsgIncr = msgIncr

	var pullMsgReq PullMessageBySeqListReq
	LogBegin("getNoInSeq ", beginSeq, endSeq)
	pullMsgReq.SeqList = u.getNoInSeq(beginSeq, endSeq)
	LogEnd("getNoInSeq ", pullMsgReq.SeqList)

	wsReq.Data, err = proto.Marshal(&pullMsgReq)
	if err != nil {
		sdkLog("Marshl failed")
		LogFReturn(err.Error())
		return err
	}
	LogBegin("WriteMsg ", wsReq.OperationID)
	err, _ = u.WriteMsg(wsReq)
	LogEnd("WriteMsg ", wsReq.OperationID, err)
	if err != nil {
		sdkLog("close conn, WriteMsg failed ", err.Error())
		u.DelCh(msgIncr)
		return err
	}

	timeout := 10
	select {
	case r := <-ch:
		sdkLog("ws ch recvMsg success: ", wsReq.OperationID)
		if r.ErrCode != 0 {
			sdkLog("pull msg failed ", r.ErrCode, r.ErrMsg, wsReq.OperationID)
			u.DelCh(msgIncr)
			return errors.New(r.ErrMsg)
		} else {
			sdkLog("pull msg success ", wsReq.OperationID)
			var pullMsg PullUserMsgResp

			pullMsg.ErrCode = 0

			var pullMsgResp PullMessageBySeqListResp
			err := proto.Unmarshal(r.Data, &pullMsgResp)
			if err != nil {
				sdkLog("Unmarshal failed ", err.Error())
				LogFReturn(err.Error())
				return err
			}
			pullMsg.Data.Group = pullMsgResp.GroupUserMsg
			pullMsg.Data.Single = pullMsgResp.SingleUserMsg
			pullMsg.Data.MaxSeq = pullMsgResp.MaxSeq
			pullMsg.Data.MinSeq = pullMsgResp.MinSeq

			u.seqMsgMutex.Lock()

			arrMsg := ArrMsg{}
			//	sdkLog("pullmsg data: ", pullMsgResp.SingleUserMsg, pullMsg.Data.Single)
			for i := 0; i < len(pullMsg.Data.Single); i++ {
				for j := 0; j < len(pullMsg.Data.Single[i].List); j++ {
					sdkLog("open_im pull one msg: |", pullMsg.Data.Single[i].List[j].ClientMsgID, "|")
					sdkLog("pull all: |", pullMsg.Data.Single[i].List[j].Seq, pullMsg.Data.Single[i].List[j])

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
					//	arrMsg.SingleData = append(arrMsg.SingleData, singleMsg)
					u.seqMsg[pullMsg.Data.Single[i].List[j].Seq] = singleMsg
					sdkLog("into map, seq: ", pullMsg.Data.Single[i].List[j].Seq, pullMsg.Data.Single[i].List[j].ClientMsgID, pullMsg.Data.Single[i].List[j].ServerMsgID)
					/*
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
					*/
				}
			}

			/*
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
			*/

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
					//	arrMsg.GroupData = append(arrMsg.GroupData, groupMsg)
					u.seqMsg[pullMsg.Data.Group[i].List[j].Seq] = groupMsg
					sdkLog("into map, seq: ", pullMsg.Data.Group[i].List[j].Seq, pullMsg.Data.Group[i].List[j].ClientMsgID, pullMsg.Data.Group[i].List[j].ServerMsgID)
					sdkLog("pull all: |", pullMsg.Data.Group[i].List[j].Seq, pullMsg.Data.Group[i].List[j])
					/*
						ctype := pullMsg.Data.Group[i].List[j].ContentType
						if ctype > GroupTipBegin && ctype < GroupTipEnd {
							u.doGroupMsg(groupMsg)
							sdkLog("doGroupMsg ", groupMsg)
						}
					*/
				}
			}
			u.seqMsgMutex.Unlock()

			u.seqMsgMutex.RLock()
			for i := beginSeq; i <= endSeq; i++ {
				v, ok := u.seqMsg[i]
				if ok {
					if v.SessionType == SingleChatType {
						arrMsg.SingleData = append(arrMsg.SingleData, v)
						if v.ContentType > SingleTipBegin && v.ContentType < SingleTipEnd {
							var msgRecv MsgData
							msgRecv.ContentType = v.ContentType
							msgRecv.Content = v.Content
							msgRecv.SendID = v.SendID
							msgRecv.RecvID = v.RecvID
							LogBegin("doFriendMsg ", msgRecv)
							u.doFriendMsg(msgRecv)
							LogEnd("doFriendMsg ", msgRecv)
						}
					} else if v.SessionType == GroupChatType {
						arrMsg.GroupData = append(arrMsg.GroupData, v)
						if v.ContentType > GroupTipBegin && v.ContentType < GroupTipEnd {
							LogBegin("doGroupMsg ", v)
							u.doGroupMsg(v)
							LogEnd("doGroupMsg ", v)
						}
					} else {
						sdkLog("type failed, ", v.SessionType, v)
					}
				} else {
					sdkLog("seq no in map, failed, seq: ", i)
				}
			}
			u.seqMsgMutex.RUnlock()

			sdkLog("triggerCmdNewMsgCome len: ", len(arrMsg.SingleData), len(arrMsg.GroupData))
			err = u.triggerCmdNewMsgCome(arrMsg)
			if err != nil {
				sdkLog("triggerCmdNewMsgCome failed, ", err.Error())
			}
			u.DelCh(msgIncr)
		}
	case <-time.After(time.Second * time.Duration(timeout)):
		sdkLog("ws ch recvMsg timeout,", wsReq.OperationID)
		u.DelCh(msgIncr)
	}
	return nil
}

func (u *UserRelated) pullOldMsgAndMergeNewMsg(beginSeq int64, endSeq int64) (err error) {
	return u.pullOldMsgAndMergeNewMsgByWs(beginSeq, endSeq)
}

func (u *UserRelated) getUserNewestSeq() (int64, error) {
	LogBegin()
	resp, err := post2Api(newestSeqRouter, paramsNewestSeqReq{ReqIdentifier: 1001, OperationID: operationIDGenerator(), SendID: u.LoginUid, MsgIncr: 1}, u.token)
	if err != nil {
		LogFReturn(0, err.Error())
		return 0, err
	}
	var seqResp paramsNewestSeqResp
	err = json.Unmarshal(resp, &seqResp)
	if err != nil {
		sdkLog("UnMarshal failed, ", err.Error())
		LogFReturn(0, err.Error())
		return 0, err
	}

	if seqResp.ErrCode != 0 {
		sdkLog("errcode: ", seqResp.ErrCode, "errmsg: ", seqResp.ErrMsg)
		LogFReturn(0, seqResp.ErrMsg)
		return 0, errors.New(seqResp.ErrMsg)
	}
	LogSReturn(seqResp.Data.Seq, nil)
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
	sdkLog("doFriendMsg ", msg)
	if u.cb == nil || u.friendListener == nil {
		sdkLog("listener is null")
		return
	}

	if msg.SendID == u.LoginUid && msg.SenderPlatformID == u.SvrConf.Platform {
		sdkLog("sync msg ", msg.ContentType)
		return
	}

	go func() {
		switch msg.ContentType {
		case AddFriendTip:
			sdkLog("addFriendNew ", msg)
			u.addFriendNew(&msg) //
		case AcceptFriendApplicationTip:
			sdkLog("acceptFriendApplicationNew ", msg)
			u.acceptFriendApplicationNew(&msg)
		case RefuseFriendApplicationTip:
			sdkLog("refuseFriendApplicationNew ", msg)
			u.refuseFriendApplicationNew(&msg)
		case SetSelfInfoTip:
			sdkLog("setSelfInfo ", msg)
			u.setSelfInfo(&msg)
		case KickOnlineTip:
			sdkLog("kickOnline ", msg)
			u.kickOnline(&msg)
		default:
			sdkLog("type failed, ", msg)
		}
	}()
}

func (u *UserRelated) acceptFriendApplicationNew(msg *MsgData) {
	LogBegin(msg.ContentType, msg.ServerMsgID, msg.ClientMsgID)
	u.syncFriendList()
	sdkLog(msg.SendID, msg.RecvID)
	sdkLog("acceptFriendApplicationNew", msg.ServerMsgID, msg)

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
