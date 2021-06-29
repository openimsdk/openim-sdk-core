package open_im_sdk

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func closeListenerCh() {
	if ConversationCh != nil {
		close(ConversationCh)
		ConversationCh = nil
	}
	if InitCh != nil {
		close(InitCh)
		InitCh = nil
	}
	if groupCh != nil {
		close(groupCh)
		groupCh = nil
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

	sdkLog("init linsten channel, init addr success, ", config)
	go im.run()
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

func (im *IMManager) login(uid, tk string, cb Base) {
	token = tk
	LoginUid = uid

	err := initDBX(LoginUid)
	if err != nil {
		cb.OnError(ErrCodeInitLogin, err.Error())
		return
	}

	_, err = getUserNewestSeq()
	if err != nil {
		cb.OnError(ErrCodeInitLogin, err.Error())
		return
	}

	err = triggerCmdReLogin()
	if err != nil {
		cb.OnError(ErrCodeInitLogin, err.Error())
		return
	}
	im.forcedSynchronization()
	cb.OnSuccess("")
	sdkLog("login ok", uid, tk)
}

func (im *IMManager) logout(cb Base) {

	atomic.SwapInt32(&WsState, 100)
	if im.conn != nil {
		im.conn.Close()
		im.conn = nil
	}
	closeDB()
	LoginUid = ""
	token = ""

	cb.OnSuccess("")
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
	ForceSyncMsg()
	ForceSyncLoginUerInfo()
}

func (im *IMManager) doMsg(message []byte) {
	sdkLog("do Msg: ", string(message))
	var msg Msg
	if err := json.Unmarshal(message, &msg); err != nil {
		sdkLog("Unmarshal failed, err: ", err.Error())
		return
	}
	if msg.ErrCode != 0 {
		sdkLog("msg errcode: ", msg.ErrCode, " errmsg: ", msg.ErrMsg)
		return
	}

	//local
	maxSeq, err := getLocalMaxSeq(msg.Data.RecvID)
	if err != nil {
		sdkLog("getLocalMaxSeq failed, ", err.Error())
		return
	}
	if maxSeq >= msg.Data.Seq {
		sdkLog("seq error, ", maxSeq, msg.Data.Seq)
		return
	}

	//svr  17    15
	if msg.Data.Seq-maxSeq > 1 {
		sdkLog("pull msg: ", maxSeq+1, msg.Data.Seq-1)
		pullOldMsgAndMergeNewMsg(maxSeq+1, msg.Data.Seq-1)
	}
	arrMsg := ArrMsg{}
	arrMsg.Data = append(arrMsg.Data, msg.Data)
	triggerCmdNewMsgCome(arrMsg)

	if msg.Data.ContentType == TransferGroupOwnerTip ||
		msg.Data.ContentType == CreateGroupTip ||
		msg.Data.ContentType == GroupApplicationResponseTip ||
		msg.Data.ContentType == JoinGroupTip ||
		msg.Data.ContentType == QuitGroupTip ||
		msg.Data.ContentType == SetGroupInfoTip ||
		msg.Data.ContentType == AcceptGroupApplicationTip ||
		msg.Data.ContentType == RefuseGroupApplicationTip ||
		msg.Data.ContentType == KickGroupMemberTip ||
		msg.Data.ContentType == InviteUserToGroupTip {
		groupManager.doGroupMsg(msg.Data)
	}

	if msg.Data.SessionType == SingleChatType {
		if msg.Data.ContentType == AddFriendTip {
			im.doFriendMsg(msg.Data)
			//	triggerCmdFriendApplication()
		} else if msg.Data.ContentType == AcceptFriendApplicationTip {
			//	triggerCmdAcceptFriend(msg.Data.SendID)
			//	triggerCmdFriend()
			im.doFriendMsg(msg.Data)
		} else if msg.Data.ContentType == RefuseFriendApplicationTip {
			//	triggerCmdRefuseFriend(msg.Data.SendID)
			im.doFriendMsg(msg.Data)
		} else if msg.Data.ContentType == KickOnlineTip {
			triggerCmdReLogin()
		} else if msg.Data.ContentType == SetSelfInfoTip {
			im.doFriendMsg(msg.Data)
		}
	} else if msg.Data.SessionType == GroupChatType {
		if msg.Data.ContentType == TransferGroupOwnerTip ||
			msg.Data.ContentType == CreateGroupTip ||
			msg.Data.ContentType == GroupApplicationResponseTip ||
			msg.Data.ContentType == JoinGroupTip ||
			msg.Data.ContentType == QuitGroupTip ||
			msg.Data.ContentType == SetGroupInfoTip ||
			msg.Data.ContentType == AcceptGroupApplicationTip ||
			msg.Data.ContentType == RefuseGroupApplicationTip ||
			msg.Data.ContentType == KickGroupMemberTip ||
			msg.Data.ContentType == InviteUserToGroupTip {
			//		groupManager.doGroupMsg(msg.Data)
		}
		// group

	}
}

func (im *IMManager) syncSeq2Msg() error {
	//local
	maxSeq, err := getLocalMaxSeq(LoginUid)
	if err != nil {
		return err
	}

	//svr
	newestSeq, err := getUserNewestSeq()
	if err != nil {
		return err
	}

	if maxSeq >= newestSeq {
		log(fmt.Sprintf("SyncSeq2Msg maxSeq[%d] >= msg.Data.Seq[%d] ", maxSeq, newestSeq))
		return nil
	}

	if newestSeq > maxSeq {
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

	userLocal, err := getUserInfoFromLocal()
	if err != nil {
		return err
	}
	sdkLog("getUserInfoFromLocal ok, user: ", userLocal)

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

func getUserInfoFromLocal() (userInfo, error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	var u userInfo
	rows, err := initDB.Query("select * from user limit 1 ")
	if err == nil {
		for rows.Next() {
			err = rows.Scan(&u.Uid, &u.Name, &u.Icon, &u.Gender, &u.Mobile, &u.Birth, &u.Email, &u.Ex)
			if err != nil {
				sdkLog("rows.Scan failed, ", err.Error())
				continue
			}
		}
		return u, nil
	} else {
		sdkLog("db Query faile, ", err.Error())
		return u, err
	}
}

func (im *IMManager) doWsCmd(cmd cmd2Value) {
	switch cmd.Cmd {
	case CmdGeyLoginUserInfo:
		im.syncLoginUserInfo()
	case CmdReLogin:
		if im.conn != nil {
			im.conn.Close()
		}
		im.syncSeq2Msg()
	}
}

func (im *IMManager) reConn(conn *websocket.Conn) *websocket.Conn {
	if conn != nil {
		conn.Close()
	}
	SdkInitManager.cb.OnConnecting()
	url := fmt.Sprintf("%s?sendID=%s&token=%s&platformID=%d", SvrConf.IpWsAddr, LoginUid, token, SvrConf.Platform)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		im.LoginState = LoginFailed
		SdkInitManager.cb.OnConnectFailed(ErrCodeInitLogin, err.Error())
		sdkLog("websocket dial failed, ", SvrConf.IpWsAddr, LoginUid, token, SvrConf.Platform, err.Error())
		return nil
	}
	im.syncSeq2Msg()
	SdkInitManager.cb.OnConnectSuccess()
	im.LoginState = LoginSuccess
	sdkLog("conn to ws ok, my ip: ", getMyIp())
	return conn
}

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

func (im *IMManager) run() {
	go im.runCmd()
	for {
		if token == "" {
			sdkLog("waiting login... ")
			time.Sleep(time.Duration(2) * time.Second)
			continue
		}

		if im.conn == nil {
			im.conn = im.reConn(nil)
		}
		if im.conn != nil {
			//im.syncSeq2Msg()
			sdkLog("conn ws ok, start readmessage")
			msgType, message, err := im.conn.ReadMessage()
			tvalue := atomic.LoadInt32(&WsState)
			if tvalue == 100 {
				atomic.CompareAndSwapInt32(&WsState, 100, 0)
				sdkLog("ws stop, return...")
				return
			}
			if err != nil {
				sdkLog("readmesage failed, reconn to ws, ", err.Error())
				im.conn = im.reConn(im.conn)
			} else {
				if msgType == websocket.CloseMessage {
					sdkLog("recv closemessage, reconn to ws, ")
					im.conn = im.reConn(im.conn)
				} else if msgType == websocket.TextMessage {
					sdkLog("recv textmessage, domsg")
					im.doMsg(message)
				}
			}
		} else {
			time.Sleep(time.Duration(2) * time.Second)
		}
	}
}

func pullOldMsgAndMergeNewMsg(beginSeq int64, endSeq int64) (err error) {
	data := paramsPullUserMsgDataReq{SeqBegin: beginSeq, SeqEnd: endSeq}
	resp, err := post2Api(pullUserMsgRouter, paramsPullUserMsg{ReqIdentifier: 1002, OperationID: operationIDGenerator(), SendID: LoginUid, MsgIncr: rand.Int(), Data: data}, token)
	if err != nil {
		sdkLog("post2Api failed, ", pullUserMsgRouter, LoginUid, beginSeq, endSeq, err.Error())
		return err
	}
	var pullMsg PullUserMsgResp
	err = json.Unmarshal(resp, &pullMsg)
	if err != nil {
		sdkLog("unmarshal failed, ", err.Error())
		return err
	}

	triggerCmd := make(map[int32]struct{})

	if pullMsg.ErrCode == 0 {
		arrMsg := ArrMsg{}
		for i := 0; i < len(pullMsg.Data.Single); i++ {
			for j := 0; j < len(pullMsg.Data.Single[i].List); j++ {
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
					IsEmphasize:      true,
					SenderPlatformID: pullMsg.Data.Single[i].List[j].SenderPlatformID,
				}
				arrMsg.Data = append(arrMsg.Data, singleMsg)
				if pullMsg.Data.Single[i].List[j].ContentType == AddFriendTip ||
					pullMsg.Data.Single[i].List[j].ContentType == AcceptFriendApplicationTip ||
					pullMsg.Data.Single[i].List[j].ContentType == RefuseFriendApplicationTip ||
					pullMsg.Data.Single[i].List[j].ContentType == SetSelfInfoTip {
					var msgRecv MsgData
					msgRecv.ContentType = pullMsg.Data.Single[i].List[j].ContentType
					msgRecv.Content = pullMsg.Data.Single[i].List[j].Content
					msgRecv.SendID = pullMsg.Data.Single[i].List[j].SendID
					msgRecv.RecvID = pullMsg.Data.Single[i].List[j].RecvID
					SdkInitManager.doFriendMsg(msgRecv)
				}

				triggerCmd[pullMsg.Data.Single[i].List[j].ContentType] = struct{}{}
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
					IsEmphasize:      true,
					SenderPlatformID: pullMsg.Data.Group[i].List[j].SenderPlatformID}

				ctype := pullMsg.Data.Group[i].List[j].ContentType
				if ctype == TransferGroupOwnerTip ||
					ctype == CreateGroupTip ||
					ctype == GroupApplicationResponseTip ||
					ctype == JoinGroupTip ||
					ctype == QuitGroupTip ||
					ctype == SetGroupInfoTip ||
					ctype == AcceptGroupApplicationTip ||
					ctype == RefuseGroupApplicationTip ||
					ctype == KickGroupMemberTip ||
					ctype == InviteUserToGroupTip {
					groupManager.doGroupMsg(groupMsg)
				}
			}
		}

		for i := 0; i < len(pullMsg.Data.Group); i++ {
			// 整理群消息
		}

		err = triggerCmdNewMsgCome(arrMsg)
		if err != nil {
			sdkLog("triggerCmdNewMsgCome failed, ", err.Error())
		}

		for k := range triggerCmd {
			if k == AddFriendTip {
				//	err := triggerCmdFriendApplication()
				if err != nil {
					sdkLog("triggerCmdFriendApplication failed, ", err.Error())
				}
			} else if k == AcceptFriendApplicationTip {
				//	err := triggerCmdFriend()
				if err != nil {
					sdkLog("triggerCmdFriend failed, ", err.Error())
				}
			} else if k == KickOnlineTip {
				//do nothing
			} else if k == SetSelfInfoTip {
				//FriendObj.doFriendList()
			} else if k == TransferGroupOwnerTip {
			} else if k == GroupApplicationResponseTip {
			}
		}
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
func getUserInfoByUid(uid string) (*userInfo, error) {
	var uidList []string
	uidList = append(uidList, uid)
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
		sdkLog("failed, no user :", uid)
		return nil, errors.New("no user")
	}
	return &userResp.Data[0], nil
}

var initDB *sql.DB
var mRWMutex *sync.RWMutex

func closeDB() error {
	if initDB != nil {
		if err := initDB.Close(); err != nil {
			return err
		}
		initDB = nil
	}
	return nil
}
func initDBX(uid string) error {
	if mRWMutex == nil {
		mRWMutex = new(sync.RWMutex)
	}
	if uid == "" {
		return errors.New("no uid")
	}
	if initDB != nil {
		return errors.New("db already opened")
	}
	db, err := sql.Open("sqlite3", SvrConf.DbDir+"OpenIM_"+uid+".db")
	sdkLog("open db:", SvrConf.DbDir+"OpenIM_"+uid+".db")
	if err != nil {
		sdkLog("failed open db:", SvrConf.DbDir+"OpenIM_"+uid+".db", err.Error())
		return err
	}
	initDB = db
	//(&u.Uid, &u.Name, &u.Icon, &u.Gender, &u.Mobile, &u.Birth, u.Email, &u.Ex)
	table := "CREATE TABLE if not exists `user` " +
		"(`uid` varchar(64) NOT NULL , " +
		"`name` varchar(64) DEFAULT NULL , " +
		"`icon` varchar(1024) DEFAULT NULL , " +
		"`gender` int(11) DEFAULT NULL , " +
		"`mobile` varchar(32) DEFAULT NULL , " +
		"`birth` varchar(16) DEFAULT NULL , " +
		"`email` varchar(64) DEFAULT NULL , " +
		"`ex` varchar(1024) DEFAULT NULL,  " +
		" PRIMARY KEY (uid) " +
		")"
	_, err = db.Exec(table)
	if err != nil {
		log(fmt.Sprintf("table user err = %s", err.Error()))
		return err
	}

	table = `create table if not exists  black_list (
   	 	uid VARCHAR (64) PRIMARY KEY  NOT NULL,
    	name VARCHAR(64) NULL ,
     	icon varchar(1024) DEFAULT NULL , 
     	gender int(11) DEFAULT NULL , 
     	mobile varchar(32) DEFAULT NULL ,
    	birth varchar(16) DEFAULT NULL , 
  	 	email varchar(64) DEFAULT NULL , 
  	 	ex varchar(1024) DEFAULT NULL
        )`
	_, err = db.Exec(table)
	if err != nil {
		log(fmt.Sprintf("table black_list err = %s", err.Error()))
		return err
	}

	table = `
      create table if not exists friend_request (
    	uid VARCHAR (64) PRIMARY KEY  NOT NULL,
    	name VARCHAR(64) NULL ,
     	icon varchar(1024) DEFAULT NULL , 
     	gender int(11) DEFAULT NULL , 
     	mobile varchar(32) DEFAULT NULL ,
    	birth varchar(16) DEFAULT NULL , 
  	 	email varchar(64) DEFAULT NULL , 
  	 	ex varchar(1024) DEFAULT NULL,
      	flag int(11) NOT NULL DEFAULT 0,
      	req_message varchar(255) DEFAULT NULL,
     	create_time  varchar(255) NOT NULL
      )`
	_, err = db.Exec(table)
	if err != nil {
		log(fmt.Sprintf("table friend_request err = %s", err.Error()))
		return err
	}

	//Apply by yourself to add other people's friends form
	table = `
      create table if not exists self_apply_to_other_request (
    	uid VARCHAR (64) PRIMARY KEY  NOT NULL,
    	name VARCHAR(64) NULL ,
     	icon varchar(1024) DEFAULT NULL , 
     	gender int(11) DEFAULT NULL , 
     	mobile varchar(32) DEFAULT NULL ,
    	birth varchar(16) DEFAULT NULL , 
  	 	email varchar(64) DEFAULT NULL , 
  	 	ex varchar(1024) DEFAULT NULL,
      	flag int(11) NOT NULL DEFAULT 0,
      	req_message varchar(255) DEFAULT NULL,
     	create_time  varchar(255) NOT NULL
      )`
	_, err = db.Exec(table)
	if err != nil {
		log(fmt.Sprintf("table friend_request err = %s", err.Error()))
		return err
	}

	table = ` CREATE TABLE IF NOT EXISTS friend_info(
     uid VARCHAR (64) PRIMARY KEY  NOT NULL,
     name VARCHAR(64) NULL ,
     comment varchar(255) DEFAULT NULL,
     icon varchar(1024) DEFAULT NULL , 
     gender int(11) DEFAULT NULL , 
     mobile varchar(32) DEFAULT NULL ,
     birth varchar(16) DEFAULT NULL , 
  	 email varchar(64) DEFAULT NULL , 
  	 ex varchar(1024) DEFAULT NULL
 	)`
	_, err = db.Exec(table)
	if err != nil {
		log(fmt.Sprintf("table friend_info err = %s", err.Error()))
		return err
	}

	table = `create table if not exists  chat_log (
      msg_id varchar(128)   NOT NULL,
	  send_id varchar(255)   NOT NULL ,
	  is_read int(255) NOT NULL ,
	  seq int(255) DEFAULT NULL ,
	  status int(11) NOT NULL ,
	  session_type int(11) NOT NULL ,
	  recv_id varchar(255)   NOT NULL ,
	  content_type int(11) NOT NULL ,
	  msg_from int(11) NOT NULL ,
	  content varchar(1000)   NOT NULL ,
	  remark varchar(100)    DEFAULT NULL ,
	  sender_platform_id int(11) NOT NULL ,
	  send_time INTEGER(255) DEFAULT NULL ,
	  create_time INTEGER (255) DEFAULT NULL,
	  PRIMARY KEY (msg_id) 
	)`
	_, err = db.Exec(table)
	if err != nil {
		log(fmt.Sprintf("table chat_log err = %s", err.Error()))
		return err
	}

	table = `create table if not exists  conversation (
	   conversation_id varchar(128) NOT NULL,
	  conversation_type int(11) NOT NULL,
	  user_id varchar(128)  DEFAULT NULL,
	  group_id varchar(128)  DEFAULT NULL,
	  show_name varchar(128)  NOT NULL,
	  face_url varchar(128)  NOT NULL,
	  recv_msg_opt int(11) NOT NULL ,
	  unread_count int(11) NOT NULL ,
	  latest_msg varchar(255)  NOT NULL ,
      latest_msg_send_time INTEGER(255)  NOT NULL ,
	  draft_text varchar(255)  DEFAULT NULL ,
	  draft_timestamp INTEGER(255)  DEFAULT NULL ,
	  is_pinned int(10) NOT NULL ,
	  PRIMARY KEY (conversation_id)
	)`

	_, err = db.Exec(table)
	if err != nil {
		log(fmt.Sprintf("table conversation err = %s", err.Error()))
		return err
	}

	return nil
}

func getLocalMaxSeq(uid string) (int64, error) {
	type MaxSeq struct {
		Seq int64
	}

	mRWMutex.Lock()
	defer mRWMutex.Unlock()

	var maxSeq MaxSeq

	rows, err := initDB.Query(fmt.Sprintf("select IFNULL(max(seq), 0) from chat_log where send_id = '%s' or recv_id = '%s'", uid, uid))
	defer rows.Close()
	if err != nil {
		log(fmt.Sprintf("1111 getLocalMaxSeq err = %s", err.Error()))
		return 0, nil
	}

	for rows.Next() {
		err = rows.Scan(&maxSeq.Seq)
		if err != nil {
			sdkLog(fmt.Sprintf("getLocalMaxSeq rows.Scan err = %s", err.Error()))
			continue
		}
	}

	return maxSeq.Seq, nil
}

func replaceIntoUser(info *userInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("replace into `user`(uid, `name`, icon, gender, mobile, birth, email, ex) values(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		sdkLog("db prepare failed, ", err.Error())
		return err
	}

	_, err = stmt.Exec(info.Uid, info.Name, info.Icon, info.Gender, info.Mobile, info.Birth, info.Email, info.Ex)
	if err != nil {
		sdkLog("db exec failed, ", err.Error())
		return err
	}
	return nil
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
