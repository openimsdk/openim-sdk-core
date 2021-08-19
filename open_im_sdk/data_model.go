package open_im_sdk

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
)

func closeDB() error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	if initDB != nil {
		if err := initDB.Close(); err != nil {
			sdkLog("close db failed, ", err.Error())
			return err
		}
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
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	if initDB != nil {
		initDB.Close()
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
		sdkLog(err.Error())
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
		sdkLog(err.Error())
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
		sdkLog(err.Error())
		return err
	}

	//Apply by yourself to add other people's friends form
	table = `
      create table if not exists self_apply_to_friend_request (
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
		sdkLog(err.Error())
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
		sdkLog(err.Error())
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
      sender_face_url varchar(255) DEFAULT NULL,
      sender_nick_name varchar(255) DEFAULT NULL,
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
		sdkLog(err.Error())
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
		sdkLog(err.Error())
		return err
	}

	table = `create table if not exists  group_info (
    	group_id varchar(255) NOT NULL,
    	name varchar(255) DEFAULT NULL,
    	introduction varchar(255) DEFAULT NULL,
    	notification varchar(255) DEFAULT NULL,
    	face_url varchar(255) DEFAULT NULL,
    	create_time INTEGER(255) DEFAULT NULL,
    	ex varchar(255) DEFAULT NULL,
    	PRIMARY KEY (group_id)
	)`

	_, err = db.Exec(table)
	if err != nil {
		sdkLog(err.Error())
		return err
	}

	table = `create table if not exists  group_member (
		group_id varchar(255) NOT NULL,
		uid varchar(255) NOT NULL,
		nickname varchar(255) DEFAULT NULL,
		user_group_face_url varchar(255) DEFAULT NULL,
		administrator_level int(11) NOT NULL,
		join_time INTEGER(255) NOT NULL,
		PRIMARY KEY (group_id,uid)
	)`

	_, err = db.Exec(table)
	if err != nil {
		sdkLog(err.Error())
		return err
	}

	table = `create table if not exists  group_request (
		id int(11) NOT NULL,
		group_id varchar(255) NOT NULL,
		from_user_id varchar(255) NOT NULL,
		to_user_id varchar(255) NOT NULL,
		flag int(10) NOT NULL DEFAULT '0',
		req_msg varchar(255) DEFAULT '',
		handled_msg varchar(255) DEFAULT '',
		create_time INTEGER(255) NOT NULL,
		from_user_nickname varchar(255) DEFAULT '',
		to_user_nickname varchar(255) DEFAULT NULL,
		from_user_face_url varchar(255) DEFAULT '',
		to_user_face_url varchar(255) DEFAULT '',
		handled_user varchar(255) DEFAULT '',
    	is_read int(10) NOT NULL DEFAULT '0',
		PRIMARY KEY (id)
	)`

	_, err = db.Exec(table)
	if err != nil {
		sdkLog(err.Error())
		return err
	}

	table = `create table if not exists  self_apply_to_group_request (
		group_id varchar(255) NOT NULL,
		flag int(10) NOT NULL DEFAULT '0',
		req_msg varchar(255) DEFAULT '',
		create_time INTEGER(255) NOT NULL,
		primary key (group_id)
	)`

	_, err = db.Exec(table)
	if err != nil {
		sdkLog(err.Error())
		return err
	}

	return nil
}

//func getLocalMaxSeq(uid string) (int64, error) {
//	type MaxSeq struct {
//		Seq int64
//	}
//
//	mRWMutex.Lock()
//	defer mRWMutex.Unlock()
//
//	var maxSeq MaxSeq
//
//	rows, err := initDB.Query(fmt.Sprintf("select IFNULL(max(seq), 0) from chat_log where send_id = '%s' or recv_id = '%s'", uid, uid))
//	if err != nil {
//		sdkLog(fmt.Sprintf("1111 getLocalMaxSeq err = %s", err.Error()))
//		return 0, nil
//	}
//
//	for rows.Next() {
//		err = rows.Scan(&maxSeq.Seq)
//		if err != nil {
//			sdkLog(fmt.Sprintf("getLocalMaxSeq rows.Scan err = %s", err.Error()))
//			continue
//		}
//	}
//
//	return maxSeq.Seq, nil
//}

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

func getAllConversationListModel() (err error, list []*ConversationStruct) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	rows, err := initDB.Query("SELECT * FROM conversation where latest_msg_send_time!=0 order by  case when is_pinned=1 then 0 else 1 end,latest_msg_send_time DESC")
	for rows.Next() {
		c := new(ConversationStruct)
		err = rows.Scan(&c.ConversationID, &c.ConversationType, &c.UserID, &c.GroupID, &c.ShowName,
			&c.FaceURL, &c.RecvMsgOpt, &c.UnreadCount, &c.LatestMsg, &c.LatestMsgSendTime, &c.DraftText, &c.DraftTimestamp, &c.IsPinned)
		if err != nil {
			sdkLog("getAllConversationListModel ,err:", err.Error())
			continue
		} else {
			list = append(list, c)
		}
	}
	return nil, list
}

func insertConversationModel(c *ConversationStruct) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("INSERT INTO conversation(conversation_id, conversation_type, " +
		"user_id,group_id,show_name,face_url,recv_msg_opt,unread_count,latest_msg,latest_msg_send_time,draft_text,draft_timestamp,is_pinned) values(?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(c.ConversationID, c.ConversationType, c.UserID, c.GroupID, c.ShowName, c.FaceURL, c.RecvMsgOpt, c.UnreadCount, c.LatestMsg, c.LatestMsgSendTime, c.DraftText, c.DraftTimestamp, c.IsPinned)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil
}
func getConversationLatestMsgModel(conversationID string) (err error, latestMsg string) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	var s string
	rows, err := initDB.Query("SELECT latest_msg FROM conversation where  conversation_id=?", conversationID)
	if err != nil {
		log(err.Error())
		return err, ""
	}
	for rows.Next() {
		err = rows.Scan(&s)
		if err != nil {
			sdkLog("getConversationLatestMsgModel ,err:", err.Error())
			continue
		}
	}
	return nil, s
}
func setConversationLatestMsgModel(c *ConversationStruct, conversationID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update conversation set latest_msg=?,latest_msg_send_time=? where conversation_id=?")
	if err != nil {
		sdkLog(err.Error())
		return err
	}

	_, err = stmt.Exec(c.LatestMsg, c.LatestMsgSendTime, conversationID)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil

}
func setConversationFaceUrlAndNickName(c *ConversationStruct, conversationID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update conversation set show_name=?,face_url=? where conversation_id=?")
	if err != nil {
		sdkLog(err.Error())
		return err
	}

	_, err = stmt.Exec(c.ShowName, c.FaceURL, conversationID)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil

}
func judgeConversationIfExists(conversationID string) bool {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	var count int
	rows, err := initDB.Query("select count(*) from conversation where  conversation_id=?", conversationID)
	if err != nil {
		fmt.Println("judge err")
		sdkLog(err.Error())
		return false
	}
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			sdkLog("getConversationLatestMsgModel ,err:", err.Error())
			continue
		}
	}
	if count == 1 {
		return true
	} else {
		return false
	}

}
func addConversationOrUpdateLatestMsg(c *ConversationStruct, conversationID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("INSERT INTO conversation(conversation_id, conversation_type, user_id,group_id,show_name,face_url,recv_msg_opt,unread_count,latest_msg,latest_msg_send_time,draft_text,draft_timestamp,is_pinned)" +
		" values(?,?,?,?,?,?,?,?,?,?,?,?,?)" +
		"ON CONFLICT(conversation_id) DO UPDATE SET latest_msg = ?,latest_msg_send_time=?")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(c.ConversationID, c.ConversationType, c.UserID, c.GroupID, c.ShowName, c.FaceURL, c.RecvMsgOpt, c.UnreadCount, c.LatestMsg, c.LatestMsgSendTime, c.DraftText, c.DraftTimestamp, c.IsPinned, c.LatestMsg, c.LatestMsgSendTime)

	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil

}
func getOneConversationModel(conversationID string) (err error, c ConversationStruct) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	rows, err := initDB.Query("SELECT * FROM conversation where  conversation_id=?", conversationID)
	if err != nil {
		sdkLog("getOneConversationModel ,err:", err.Error())
		sdkLog(err.Error())
		return err, c
	}
	for rows.Next() {
		err = rows.Scan(&c.ConversationID, &c.ConversationType, &c.UserID, &c.GroupID, &c.ShowName,
			&c.FaceURL, &c.RecvMsgOpt, &c.UnreadCount, &c.LatestMsg, &c.LatestMsgSendTime, &c.DraftText, &c.DraftTimestamp, &c.IsPinned)
		if err != nil {
			sdkLog("getOneConversationModel ,err:", err.Error())
			continue
		}
	}
	return nil, c

}
func deleteConversationModel(conversationID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("delete from conversation where conversation_id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(conversationID)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil
}
func ResetConversation(conversationID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update  conversation set unread_count=?,latest_msg=?,latest_msg_send_time=?," +
		"draft_text=?,draft_timestamp=?,is_pinned=? where conversation_id=?")
	if err != nil {
		sdkLog("ResetConversation", err.Error())
		return err
	}
	_, err = stmt.Exec(0, "", 0, "", 0, 0, conversationID)
	if err != nil {
		sdkLog("ResetConversation err:", err.Error())
		return err
	}
	return nil

}
func setConversationDraftModel(conversationID, draftText string, DraftTimestamp int64) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update conversation set draft_text=?,latest_msg_send_time=?,draft_timestamp=?where conversation_id=?")
	if err != nil {
		sdkLog("setConversationDraftModel err:", err.Error())
		return err
	}

	_, err = stmt.Exec(draftText, DraftTimestamp, DraftTimestamp, conversationID)
	if err != nil {
		sdkLog("setConversationDraftModel err:", err.Error())
		return err
	}
	return nil

}
func pinConversationModel(conversationID string, isPinned int) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update conversation set is_pinned=? where conversation_id=?")
	if err != nil {
		sdkLog(err.Error())
		return err
	}

	_, err = stmt.Exec(isPinned, conversationID)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil

}
func setConversationUnreadCount(unreadCount int, conversationID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update conversation set unread_count=? where conversation_id=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(unreadCount, conversationID)
	if err != nil {
		return err
	}
	return nil

}
func incrConversationUnreadCount(conversationID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update conversation set unread_count = unread_count+1 where conversation_id=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(conversationID)
	if err != nil {
		return err
	}
	return nil

}
func getTotalUnreadMsgCountModel() (totalUnreadCount int32, err error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	rows, err := initDB.Query("SELECT IFNULL(SUM(unread_count), 0) FROM conversation")
	if err != nil {
		sdkLog(err.Error())
		return totalUnreadCount, err
	}
	for rows.Next() {
		err = rows.Scan(&totalUnreadCount)
		if err != nil {
			sdkLog("getTotalUnreadMsgCountModel ,err:", err.Error())
			continue
		}
	}
	return totalUnreadCount, err

}
func getMultipleConversationModel(conversationIDList []string) (err error, list []*ConversationStruct) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	rows, err := initDB.Query("SELECT * FROM conversation where conversation_id in (" + sqlStringHandle(conversationIDList) + ")")
	fmt.Println("SELECT * FROM conversation where conversation_id in (" + sqlStringHandle(conversationIDList) + ")")
	for rows.Next() {
		temp := new(ConversationStruct)
		err = rows.Scan(&temp.ConversationID, &temp.ConversationType, &temp.UserID, &temp.GroupID, &temp.ShowName,
			&temp.FaceURL, &temp.RecvMsgOpt, &temp.UnreadCount, &temp.LatestMsg, &temp.LatestMsgSendTime, &temp.DraftText, &temp.DraftTimestamp, &temp.IsPinned)
		if err != nil {
			sdkLog("getMultipleConversationModel err:", err.Error())
			return err, nil
		} else {
			list = append(list, temp)
		}
	}
	return nil, list
}
func sqlStringHandle(ss []string) (s string) {
	for i := 0; i < len(ss); i++ {
		s += "'" + ss[i] + "'"
		if i < len(ss)-1 {
			s += ","
		}
	}
	return s
}

func insertIntoTheFriendToFriendInfo(uid, name, comment, icon string, gender int32, mobile, birth, email, ex string) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("insert into friend_info(uid,name,comment,icon,gender,mobile,birth,email,ex) values (?,?,?,?,?,?,?,?,?)")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(uid, name, comment, icon, gender, mobile, birth, email, ex)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil
}

func delTheFriendFromFriendInfo(uid string) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("delete from friend_info where uid=?")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(uid)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil
}
func updateTheFriendInfo(uid, name, comment, icon string, gender int32, mobile, birth, email, ex string) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("replace into friend_info(uid,name,comment,icon,gender,mobile,birth,email,ex) values (?,?,?,?,?,?,?,?,?)")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(uid, name, comment, icon, gender, mobile, birth, email, ex)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil
}

func updateFriendInfo(uid, name, icon string, gender int32, mobile, birth, email, ex string) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update friend_info set `name` = ?, icon = ?, gender = ?, mobile = ?, birth = ?, email = ?, ex = ? where uid = ?")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(name, icon, gender, mobile, birth, email, ex, uid)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil
}

func insertIntoTheUserToBlackList(info userInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("insert into black_list(uid,name,icon,gender,mobile,birth,email,ex) values (?,?,?,?,?,?,?,?)")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(info.Uid, info.Name, info.Icon, info.Gender, info.Mobile, info.Birth, info.Email, info.Ex)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil
}

func updateBlackList(info userInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("replace into black_list(uid,name,icon,gender,mobile,birth,email,ex) values (?,?,?,?,?,?,?,?)")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(info.Uid, info.Name, info.Icon, info.Gender, info.Mobile, info.Birth, info.Email, info.Ex)
	if err != nil {
		fmt.Println(err)
		sdkLog(err.Error())
		return err
	}
	return nil
}

func delTheUserFromBlackList(uid string) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("delete from black_list where uid=?")
	if err != nil {
		fmt.Println(err)
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(uid)
	if err != nil {
		fmt.Println(err)
		sdkLog(err.Error())
		return err
	}
	return nil
}

func insertIntoTheUserToApplicationList(appUserInfo applyUserInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("insert into friend_request(uid,name,icon,gender,mobile,birth,email,ex,flag,req_message,create_time) values (?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		sdkLog("Prepare failed ", err.Error())
		return err
	}
	_, err = stmt.Exec(appUserInfo.Uid, appUserInfo.Name, appUserInfo.Icon, appUserInfo.Gender, appUserInfo.Mobile, appUserInfo.Birth, appUserInfo.Email, appUserInfo.Ex, appUserInfo.Flag, appUserInfo.ReqMessage, appUserInfo.ApplyTime)
	if err != nil {
		sdkLog("Exec failed, ", err.Error())
		return err
	}
	return nil
}

func delTheUserFromApplicationList(uid string) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("delete from friend_request where uid=?")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(uid)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil
}

func updateApplicationList(info applyUserInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("replace into friend_request(uid,name,icon,gender,mobile,birth,email,ex,flag,req_message,create_time) values (?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = stmt.Exec(info.Uid, info.Name, info.Icon, info.Gender, info.Mobile, info.Birth, info.Email, info.Ex, info.Flag, info.ReqMessage, info.ApplyTime)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func getFriendInfoByFriendUid(friendUid string) (*friendInfo, error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	stmt, err := initDB.Query("select * from friend_info  where uid=? ", friendUid)
	if err != nil {
		sdkLog("query failed, ", err.Error())
		return nil, err
	}
	var (
		uid           string
		name          string
		icon          string
		gender        int32
		mobile        string
		birth         string
		email         string
		ex            string
		comment       string
		isInBlackList int32
	)
	for stmt.Next() {
		err = stmt.Scan(&uid, &name, &comment, &icon, &gender, &mobile, &birth, &email, &ex)
		if err != nil {
			sdkLog("scan failed, ", err.Error())
			continue
		}
	}

	return &friendInfo{uid, name, icon, gender, mobile, birth, email, ex, comment, isInBlackList}, nil
}

func getLocalFriendList() ([]friendInfo, error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	stmt, err := initDB.Query("select * from friend_info")
	if err != nil {
		return nil, err
	}
	friends := make([]friendInfo, 0)
	for stmt.Next() {
		var (
			uid           string
			name          string
			icon          string
			gender        int32
			mobile        string
			birth         string
			email         string
			ex            string
			comment       string
			isInBlackList int32
		)
		err = stmt.Scan(&uid, &name, &comment, &icon, &gender, &mobile, &birth, &email, &ex)
		if err != nil {
			sdkLog("scan failed, ", err.Error())
			continue
		}

		friends = append(friends, friendInfo{uid, name, icon, gender, mobile, birth, email, ex, comment, isInBlackList})
	}
	return friends, nil
}

func getLocalFriendApplication() ([]applyUserInfo, error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	stmt, err := initDB.Query("select * from friend_request order by create_time desc")
	if err != nil {
		println(err.Error())
		return nil, err
	}
	applyUsersInfo := make([]applyUserInfo, 0)
	for stmt.Next() {
		var (
			uid        string
			name       string
			icon       string
			gender     int32
			mobile     string
			birth      string
			email      string
			ex         string
			reqMessage string
			applyTime  string
			flag       int32
		)
		err = stmt.Scan(&uid, &name, &icon, &gender, &mobile, &birth, &email, &ex, &flag, &reqMessage, &applyTime)
		if err != nil {
			sdkLog("sqlite scan failed", err.Error())
			continue
		}
		applyUsersInfo = append(applyUsersInfo, applyUserInfo{uid, name, icon, gender, mobile, birth, email, ex, reqMessage, applyTime, flag})
	}
	return applyUsersInfo, nil
}

func getLocalBlackList() ([]userInfo, error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	stmt, err := initDB.Query("select * from black_list")
	if err != nil {
		return nil, err
	}
	usersInfo := make([]userInfo, 0)
	for stmt.Next() {
		var (
			uid    string
			name   string
			icon   string
			gender int32
			mobile string
			birth  string
			email  string
			ex     string
		)
		err = stmt.Scan(&uid, &name, &icon, &gender, &mobile, &birth, &email, &ex)
		if err != nil {
			sdkLog("sqlite scan failed", err.Error())
			continue
		}
		usersInfo = append(usersInfo, userInfo{uid, name, icon, gender, mobile, birth, email, ex})
	}
	return usersInfo, nil
}

func getBlackUsInfoByUid(blackUid string) (*userInfo, error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	stmt, err := initDB.Query("select * from black_list where uid=?", blackUid)
	if err != nil {
		return nil, err
	}

	var (
		uid    string
		name   string
		icon   string
		gender int32
		mobile string
		birth  string
		email  string
		ex     string
	)
	for stmt.Next() {
		err = stmt.Scan(&uid, &name, &icon, &gender, &mobile, &birth, &email, &ex)
		if err != nil {
			sdkLog("scan failed, ", err.Error())
			continue
		}
	}

	return &userInfo{uid, name, icon, gender, mobile, birth, email, ex}, nil
}

func updateLocalTransferGroupOwner(transfer *TransferGroupOwnerReq) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()

	stmt, err := initDB.Prepare("update group_member set administrator_level = ? where group_id = ? and uid = ?")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(0, transfer.GroupID, transfer.OldOwner)
	if err != nil {
		sdkLog(err.Error())
		return err
	}

	stmt, err = initDB.Prepare("update group_member set administrator_level = ? where group_id = ? and uid = ?")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(1, transfer.GroupID, transfer.NewOwner)
	if err != nil {
		sdkLog(err.Error())
		return err
	}

	return nil
}

func insertLocalAcceptGroupApplication(addMem *groupMemberFullInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("insert into group_member(group_id,uid,nickname,user_group_face_url,administrator_level,join_time) values (?,?,?,?,?,?)")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(addMem.GroupId, addMem.UserId, addMem.NickName, addMem.FaceUrl, addMem.Role, addMem.JoinTime)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil
}
func insertIntoLocalGroupInfo(info groupInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("insert into group_info(group_id,name,introduction,notification,face_url,create_time,ex) values (?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.GroupId, info.GroupName, info.Introduction, info.Notification, info.FaceUrl, info.CreateTime, info.Ex)
	if err != nil {
		return err
	}
	return nil
}
func delLocalGroupInfo(groupId string) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("delete from group_info where group_id=?")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(groupId)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil
}
func replaceLocalGroupInfo(info groupInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("replace into group_info(group_id,name,introduction,notification,face_url,create_time,ex) values (?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.GroupId, info.GroupName, info.Introduction, info.Notification, info.FaceUrl, info.CreateTime, info.Ex)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil
}

func updateLocalGroupInfo(info groupInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update group_info set name=?,introduction=?,notification=?,face_url=? where group_id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.GroupName, info.Introduction, info.Notification, info.FaceUrl, info.GroupId)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil
}

func getLocalGroupsInfo() ([]groupInfo, error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	stmt, err := initDB.Query("select * from group_info")

	if err != nil {
		return nil, err
	}
	groupsInfo := make([]groupInfo, 0)
	for stmt.Next() {
		var (
			groupId      string
			name         string
			introduction string
			notification string
			faceUrl      string
			createTime   uint64
			ex           string
		)
		err = stmt.Scan(&groupId, &name, &introduction, &notification, &faceUrl, &createTime, &ex)
		if err != nil {
			sdkLog("sqlite scan failed", err.Error())
			continue
		}
		groupsInfo = append(groupsInfo, groupInfo{groupId, name, notification, introduction, faceUrl, ex, "", createTime, 0})
	}
	return groupsInfo, nil
}

func getLocalGroupsInfoByGroupID(groupID string) (*groupInfo, error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	stmt, err := initDB.Query("select * from group_info where group_id=?", groupID)
	if err != nil {
		return nil, err
	}
	var gInfo groupInfo
	for stmt.Next() {
		var (
			groupId      string
			name         string
			introduction string
			notification string
			faceUrl      string
			createTime   uint64
			ex           string
		)
		err = stmt.Scan(&groupId, &name, &introduction, &notification, &faceUrl, &createTime, &ex)
		if err != nil {
			sdkLog("sqlite scan failed", err.Error())
			continue
		}
		gInfo.GroupId = groupID
		gInfo.GroupName = name
		gInfo.Introduction = introduction
		gInfo.Notification = notification
		gInfo.FaceUrl = faceUrl
		gInfo.CreateTime = createTime
	}
	return &gInfo, nil
}

func findLocalGroupOwnerByGroupId(groupId string) (uid string, err error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	stmt, err := initDB.Query("select uid from group_member where group_id=? and administrator_level=?", groupId, 1)
	if err != nil {
		return "", err
	}
	for stmt.Next() {
		err = stmt.Scan(&uid)
		if err != nil {
			sdkLog(err.Error())
			continue
		}
	}

	return uid, nil
}

func getLocalGroupMemberNumByGroupId(groupId string) (num int, err error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	stmt, err := initDB.Query("select count(*) from group_member where group_id=?", groupId)
	if err != nil {
		return 0, err
	}
	for stmt.Next() {
		err = stmt.Scan(&num)
		if err != nil {
			sdkLog("getLocalGroupMemberNumByGroupId query failed, err", err.Error())
			continue
		}
	}
	return num, err
}

func getLocalGroupMemberInfoByGroupIdUserId(groupId string, uid string) (*groupMemberFullInfo, error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	stmt, err := initDB.Query("select * from group_member where group_id=? and uid=?", groupId, uid)
	if err != nil {
		sdkLog("query failed, ", err.Error())
		return nil, err
	}
	var member groupMemberFullInfo
	for stmt.Next() {
		err = stmt.Scan(&member.GroupId, &member.UserId, &member.NickName, &member.FaceUrl, &member.Role, &member.JoinTime)
		if err != nil {
			sdkLog("sqlite scan failed", err.Error())
			continue
		}
	}
	return &member, nil
}

func getLocalGroupMemberList() ([]groupMemberFullInfo, error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	stmt, err := initDB.Query("select * from group_member")

	if err != nil {
		return nil, err
	}
	groupMemberList := make([]groupMemberFullInfo, 0)
	for stmt.Next() {
		var (
			groupId          string
			uid              string
			nickname         string
			userGroupFaceUrl string
			administrator    int
			joinTime         uint64
		)
		err = stmt.Scan(&groupId, &uid, &nickname, &userGroupFaceUrl, &administrator, &joinTime)
		if err != nil {
			sdkLog("sqlite scan failed", err.Error())
			continue
		}
		groupMemberList = append(groupMemberList, groupMemberFullInfo{groupId, uid, administrator, joinTime, nickname, userGroupFaceUrl})
	}
	return groupMemberList, nil
}

func getLocalGroupMemberListByGroupID(groupId string) ([]groupMemberFullInfo, error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	stmt, err := initDB.Query("select * from group_member where group_id=?", groupId)
	if err != nil {
		return nil, err
	}
	groupMemberList := make([]groupMemberFullInfo, 0)
	for stmt.Next() {
		var (
			groupId          string
			uid              string
			nickname         string
			userGroupFaceUrl string
			administrator    int
			joinTime         uint64
		)
		err = stmt.Scan(&groupId, &uid, &nickname, &userGroupFaceUrl, &administrator, &joinTime)
		if err != nil {
			sdkLog("sqlite scan failed", err.Error())
			continue
		}
		groupMemberList = append(groupMemberList, groupMemberFullInfo{groupId, uid, administrator, joinTime, nickname, userGroupFaceUrl})
	}
	return groupMemberList, nil
}

func insertIntoLocalGroupMember(info groupMemberFullInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("insert into group_member(group_id, uid, nickname,user_group_face_url,administrator_level, join_time) values (?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.GroupId, info.UserId, info.NickName, info.FaceUrl, info.Role, info.JoinTime)
	if err != nil {
		return err
	}
	return nil
}

func delLocalGroupMember(info groupMemberFullInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("delete from group_member where group_id=? and uid=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.GroupId, info.UserId)
	if err != nil {
		return err
	}
	return nil
}
func replaceLocalGroupMemberInfo(info groupMemberFullInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("replace into group_member(group_id,uid,nickname,user_group_face_url,administrator_level, join_time) values (?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.GroupId, info.UserId, info.NickName, info.FaceUrl, info.Role, info.JoinTime)
	if err != nil {
		return err
	}
	return nil
}

func updateLocalGroupMemberInfo(info groupMemberFullInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update group_member set nickname=?,user_group_face_url=? where group_id=? and uid=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.NickName, info.FaceUrl, info.GroupId, info.UserId)
	if err != nil {
		return err
	}
	return nil
}

func insertIntoSelfApplyToGroupRequest(groupId, message string) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("replace into self_apply_to_group_request(group_id,flag,req_msg,create_time) values (?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(groupId, 0, message, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func insertMessageToLocalOrUpdateContent(message *MsgStruct) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("INSERT INTO chat_log(msg_id, send_id, is_read," +
		" seq,status, session_type, recv_id, content_type, sender_face_url,sender_nick_name,msg_from, content, remark,sender_platform_id, send_time,create_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)" +
		"ON CONFLICT(msg_id) DO UPDATE SET content = ?")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(message.ClientMsgID, message.SendID,
		getIsRead(message.IsRead), message.Seq, message.Status, message.SessionType, message.RecvID, message.ContentType, message.SenderFaceURL, message.SenderNickName,
		message.MsgFrom, message.Content, message.Remark, message.PlatformID, message.SendTime, message.CreateTime, message.Content)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil
}

func insertPushMessageToChatLog(message *MsgStruct) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("INSERT INTO chat_log(msg_id, send_id, is_read," +
		" seq,status, session_type, recv_id, content_type, sender_face_url,sender_nick_name,msg_from, content, remark,sender_platform_id, send_time,create_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)" +
		"ON CONFLICT(msg_id) DO UPDATE SET seq = ?")
	if err != nil {
		sdkLog("Prepare failed, ", err.Error())
		return err
	}
	_, err = stmt.Exec(message.ClientMsgID, message.SendID,
		getIsRead(message.IsRead), message.Seq, message.Status, message.SessionType, message.RecvID, message.ContentType, message.SenderFaceURL, message.SenderNickName,
		message.MsgFrom, message.Content, message.Remark, message.PlatformID, message.SendTime, message.CreateTime, message.Seq)
	if err != nil {
		sdkLog("Exec failed, ", err.Error())
		return err
	}
	return nil
}
func updateMessageSeq(message *MsgStruct) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update chat_log set seq=? where msg_id=?")
	if err != nil {
		sdkLog("Prepare failed, ", err.Error())
		return err
	}
	_, err = stmt.Exec(message.Seq, message.ClientMsgID)
	if err != nil {
		sdkLog("Exec failed ", err.Error())
		return err
	}
	return nil
}
func judgeMessageIfExists(message *MsgStruct) bool {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	var count int
	rows, err := initDB.Query("select count(*) from chat_log where  msg_id=?", message.ClientMsgID)
	if err != nil {
		sdkLog("Query failed, ", err.Error())
		return false
	}
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			sdkLog(err.Error())
			continue
		}
	}
	if count == 1 {
		return true
	} else {
		return false
	}
}
func getOneMessage(msgID string) (m *MsgStruct, err error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	// query
	rows, err := initDB.Query("SELECT * FROM chat_log where msg_id = ?", msgID)
	if err != nil {
		sdkLog("getOneMessage err:", err.Error(), msgID)
		return nil, err
	}
	temp := new(MsgStruct)
	for rows.Next() {
		err = rows.Scan(&temp.ClientMsgID, &temp.SendID, &temp.IsRead,
			&temp.Seq, &temp.Status, &temp.SessionType, &temp.RecvID, &temp.ContentType, &temp.SenderFaceURL, &temp.SenderNickName,
			&temp.MsgFrom, &temp.Content, &temp.Remark, &temp.PlatformID, &temp.SendTime, &temp.CreateTime)
		if err != nil {
			sdkLog("getOneMessage,err:", err.Error())
			continue
		}
	}
	return temp, nil
}

//func GetHistoryMessage(recvID string,count int)  {
//	rows, err := defInitDB.Query("SELECT * FROM chat_log where seq = 0 And recv_id =" +recvID+"And send_id = "+LoginUid )
//	if err != nil {
//		return nil, err
//	}
//	rows.c
//	for rows.Next() {
//		t = c
//		err = rows.Scan(&t.ConversationID, &t.ConversationType, &t.UserID, &t.GroupID, &t.ShowName,
//			&t.FaceURL, &t.RecvMsgOpt, &t.UnreadCount, &t.LatestMsg, &t.DraftText, &t.DraftTimestamp, &t.IsPinned)
//		if err != nil {
//			return err, nil
//		} else {
//			list = append(list, &t)
//		}
//	}
//
//}
func setSingleMessageHasRead(sendID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update chat_log set is_read=? where send_id=?And is_read=?AND session_type=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(HasRead, sendID, NotRead, SingleChatType)
	if err != nil {
		return err
	}
	return nil

}
func setGroupMessageHasRead(groupID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update chat_log set is_read=? where recv_id=?And is_read=?AND session_type=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(HasRead, groupID, NotRead, GroupChatType)
	if err != nil {
		return err
	}
	return nil
}
func setMessageStatus(msgID string, status int) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update chat_log set status=? where msg_id=?")
	if err != nil {
		sdkLog("setMessageStatus prepare failed, err: ", err.Error())
		return err
	}
	_, err = stmt.Exec(status, msgID)
	if err != nil {
		sdkLog("setMessageStatus exec failed, err: ", err.Error())
		return err
	}
	return nil
}
func setMessageStatusBySourceID(sourceID string, status, sessionType int) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update chat_log set status=? where (send_id=? or recv_id=?)AND session_type=?")
	if err != nil {
		sdkLog("prepare failed, err: ", err.Error())
		return err
	}
	_, err = stmt.Exec(status, sourceID, sourceID, sessionType)
	if err != nil {
		sdkLog("exec failed, err: ", err.Error())
		return err
	}
	return nil

}
func setMessageHasReadByMsgID(msgID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update chat_log set is_read=? where msg_id=?And is_read=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(HasRead, msgID, NotRead)
	if err != nil {
		return err
	}
	return nil

}
func getHistoryMessage(sourceConversationID string, startTime int64, count int, sessionType int) (err error, list MsgFormats) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	rows, err := initDB.Query("select * from chat_log WHERE (send_id = ? OR recv_id =? )AND (content_type<=? and content_type not in (?)or (content_type >=? and content_type <=?  and content_type not in(?,?)  ))AND status not in(?,?)AND session_type=?AND send_time<?  order by send_time DESC  LIMIT ? OFFSET 0 ",
		sourceConversationID, sourceConversationID, AcceptFriendApplicationTip, HasReadReceipt, GroupTipBegin, GroupTipEnd, SetGroupInfoTip, JoinGroupTip, MsgStatusHasDeleted, MsgStatusRevoked, sessionType, startTime, count)
	for rows.Next() {
		temp := new(MsgStruct)
		err = rows.Scan(&temp.ClientMsgID, &temp.SendID, &temp.IsRead, &temp.Seq, &temp.Status, &temp.SessionType,
			&temp.RecvID, &temp.ContentType, &temp.SenderFaceURL, &temp.SenderNickName, &temp.MsgFrom, &temp.Content, &temp.Remark, &temp.PlatformID, &temp.SendTime, &temp.CreateTime)
		if err != nil {
			sdkLog("getHistoryMessage,err:", err.Error())
			continue
		} else {
			err = msgHandleByContentType(temp)
			if err != nil {
				sdkLog("getHistoryMessage,err:", err.Error())
				continue
			}
			list = append(list, temp)
		}
	}
	return nil, list
}
func deleteMessageByConversationModel(sourceConversationID string, maxSeq int64) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("delete from chat_log  where send_id=? or recv_id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(sourceConversationID, sourceConversationID)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil
}
func deleteMessageByMsgID(msgID string) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("delete from chat_log  where msg_id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(msgID)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil
}
func updateMessageTimeAndMsgIDStatus(ClientMsgID string, sendTime int64, status int) (err error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("update chat_log set send_time=?, status=? where msg_id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(sendTime, status, ClientMsgID)
	if err != nil {
		return err
	}
	return nil
}
func getMultipleMessageModel(messageIDList []string) (err error, list []*MsgStruct) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	rows, err := initDB.Query("SELECT * FROM chat_log where msg_id in (" + sqlStringHandle(messageIDList) + ")")
	fmt.Println("SELECT * FROM conversation where conversation_id in (" + sqlStringHandle(messageIDList) + ")")
	defer rows.Close()
	for rows.Next() {
		temp := new(MsgStruct)
		err = rows.Scan(&temp.ClientMsgID, &temp.SendID, &temp.IsRead, &temp.Seq, &temp.Status, &temp.SessionType,
			&temp.RecvID, &temp.ContentType, &temp.SenderFaceURL, &temp.SenderNickName, &temp.MsgFrom, &temp.Content, &temp.Remark, &temp.PlatformID, &temp.SendTime, &temp.CreateTime)
		if err != nil {
			sdkLog("getMultipleMessageModel,err:", err.Error())
			continue
		} else {
			err = msgHandleByContentType(temp)
			if err != nil {
				sdkLog("getMultipleMessageModel,err:", err.Error())
				continue
			}
			list = append(list, temp)
		}
	}
	return nil, list
}
func getLocalMaxSeq() (seq int64, err error) {
	mRWMutex.RLock()
	defer mRWMutex.RUnlock()
	rows, err := initDB.Query("SELECT IFNULL(MAX(seq), 0) FROM chat_log")
	if err != nil {
		sdkLog("getLocalMaxSeqModel,Query err:", err.Error())
		return seq, err
	}
	for rows.Next() {
		err = rows.Scan(&seq)
		if err != nil {
			sdkLog("getLocalMaxSeqModel ,err:", err.Error())
			continue
		}
	}
	return seq, err
}

func getLoginUserInfoFromLocal() (userInfo, error) {
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

func getOwnLocalGroupApplicationList(groupId string) (*groupApplicationResult, error) {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()

	sql := "select id, group_id, from_user_id, to_user_id, flag, req_msg, handled_msg, create_time, from_user_nickname, " +
		" to_user_nickname, from_user_face_url, to_user_face_url, handled_user  from `group_request` "

	if len(groupId) > 0 {
		sql = fmt.Sprintf("%s where group_id = %s", sql, groupId)
	}

	rows, err := initDB.Query(sql)
	if err != nil {
		sdkLog("db Query getOwnLocalGroupApplicationList faild, ", err.Error())
		return nil, err
	}

	reply := groupApplicationResult{}

	for rows.Next() {
		glInfo := GroupReqListInfo{}
		err = rows.Scan(&glInfo.ID, &glInfo.GroupID, &glInfo.FromUserID, &glInfo.ToUserID, &glInfo.Flag, &glInfo.RequestMsg,
			&glInfo.HandledMsg, &glInfo.AddTime, &glInfo.FromUserNickname, &glInfo.ToUserNickname,
			&glInfo.FromUserFaceUrl, &glInfo.ToUserFaceUrl, &glInfo.HandledUser)
		if err != nil {
			sdkLog("rows.Scan failed, ", err.Error())
			continue
		}
		//if glInfo.IsRead == 0 {
		//	reply.UnReadCount++
		//}

		if glInfo.ToUserID != "0" {
			glInfo.Type = 1
		}

		if len(glInfo.HandledUser) > 0 {
			if glInfo.HandledUser == LoginUid {
				glInfo.HandleStatus = 2
			} else {
				glInfo.HandleStatus = 1
			}
		}

		if glInfo.Flag == 1 {
			glInfo.HandleResult = 1
		}

		reply.GroupApplicationList = append(reply.GroupApplicationList, glInfo)
	}

	//sql = "update group_request set is_read = 1 "
	//if len(groupId) > 0 {
	//	sql = fmt.Sprintf("%s where group_id = %s", sql, groupId)
	//}
	//stmt, err := initDB.Prepare(sql)
	//if err != nil {
	//	sdkLog("db update  faild, ", err.Error(), sql)
	//	return nil, err
	//}
	//_, err = stmt.Exec()
	//if err != nil {
	//	sdkLog("db update  faild, ", err.Error(), sql)
	//	return nil, err
	//}

	return &reply, nil
}

func insertIntoRequestToGroupRequest(info GroupReqListInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("insert into group_request(id, group_id, from_user_id, to_user_id,flag,req_msg,handled_msg, create_time,from_user_nickname,to_user_nickname,from_user_face_url,to_user_face_url,handled_user) values (?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.ID, info.GroupID, info.FromUserID, info.ToUserID, info.Flag, info.RequestMsg, info.HandledMsg, info.AddTime, info.FromUserNickname, info.ToUserNickname, info.FromUserFaceUrl, info.ToUserFaceUrl, info.HandledUser)
	if err != nil {
		return err
	}
	return nil
}

func delRequestFromGroupRequest(info GroupReqListInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("delete from group_request where id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.ID)
	if err != nil {
		return err
	}
	return nil
}

func replaceIntoRequestToGroupRequest(info GroupReqListInfo) error {
	mRWMutex.Lock()
	defer mRWMutex.Unlock()
	stmt, err := initDB.Prepare("replace into group_request(id,group_id, from_user_id, to_user_id,flag,req_msg,handled_msg, create_time,from_user_nickname,to_user_nickname,from_user_face_url,to_user_face_url,handled_user) values (?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.ID, info.GroupID, info.FromUserID, info.ToUserID, info.Flag, info.RequestMsg, info.HandledMsg, info.AddTime, info.FromUserNickname, info.ToUserNickname, info.FromUserFaceUrl, info.ToUserFaceUrl, info.HandledUser)
	if err != nil {
		return err
	}
	return nil
}
