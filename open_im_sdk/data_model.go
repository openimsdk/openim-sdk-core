package open_im_sdk

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (u *UserRelated) closeDB() error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	if u.db != nil {
		if err := u.db.Close(); err != nil {
			sdkLog("close db failed, ", err.Error())
			return err
		}
	}
	return nil
}

func (u *UserRelated) closeDBSetNil() error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	if u.db != nil {
		if err := u.db.Close(); err != nil {
			sdkLog("close db failed, ", err.Error())
			return err
		}
	}
	u.db = nil
	return nil
}

func (u *UserRelated) reOpenDB(uid string) error {
	db, err := sql.Open("sqlite3", SvrConf.DbDir+"OpenIM_"+uid+".db")
	sdkLog("open db:", SvrConf.DbDir+"OpenIM_"+uid+".db")
	if err != nil {
		sdkLog("failed open db:", SvrConf.DbDir+"OpenIM_"+uid+".db", err.Error())
		return err
	}
	u.db = db
	return nil
}

func (u *UserRelated) initDBX(uid string) error {
	//if u.mRWMutex == nil {
	//	u.mRWMutex = new(sync.RWMutex)
	//}
	if uid == "" {
		return errors.New("no uid")
	}
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	if u.db != nil {
		u.db.Close()
	}
	db, err := sql.Open("sqlite3", SvrConf.DbDir+"OpenIM_"+uid+".db")
	sdkLog("open db:", SvrConf.DbDir+"OpenIM_"+uid+".db")
	if err != nil {
		sdkLog("failed open db:", SvrConf.DbDir+"OpenIM_"+uid+".db", err.Error())
		return err
	}
	u.db = db
	//(&u.Uid, &u.Name, &u.Icon, &u.Gender, &u.Mobile, &u.Birth, u.Email, &u.Ex)
	table := `CREATE TABLE if not exists user(
        user_id char(64) PRIMARY KEY NOT NULL , 
		name varchar(64) DEFAULT NULL , 
		face_url varchar(100) DEFAULT NULL , 
		gender int(11) DEFAULT NULL , 
		phone_number varchar(32) DEFAULT NULL , 
		birth INTEGER DEFAULT NULL , 
		email varchar(64) DEFAULT NULL , 
		ex varchar(1024) DEFAULT NULL, 
         )`
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
	  is_read int(10) NOT NULL ,
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
      is_filter int(10) NOT NULL,
	  PRIMARY KEY (msg_id) 
	)`
	_, err = db.Exec(table)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	table = `create table if not exists  error_chat_log (
      seq int(255) NOT NULL ,
      msg_id varchar(128)   NOT NULL,
	  send_id varchar(255)   NOT NULL ,
	  is_read int(255) NOT NULL ,
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
	  PRIMARY KEY (seq) 
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

	table = `create table if not exists  my_local_data (
		user_id varchar(128)  DEFAULT NULL,
		seq int(10) NOT NULL DEFAULT '1',
		primary key (user_id)
	)`

	_, err = db.Exec(table)
	if err != nil {
		sdkLog(err.Error())
		return err
	}

	return nil
}

func (u *UserRelated) Prepare(query string) (*sql.Stmt, error) {
	if u.db == nil {
		err := u.reOpenDB(u.LoginUid)
		if err != nil {
			sdkLog("reOpenDB failed ", u.LoginUid)
			return nil, err
		}
	}
	return u.db.Prepare(query)
}

/*
func (u *UserRelated) setLocalMaxConSeq(seq int) (err error) {
	sdkLog("setLocalMaxConSeq start ", seq)
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()

	stmt, err := u.Prepare("replace into my_local_data(user_id, seq) values (?,?)")
	if err != nil {
		sdkLog("set failed", err.Error())
		return err
	}

	_, err = stmt.Exec(u.LoginUid, seq)
	if err != nil {
		sdkLog("stmt failed,", err.Error())
		return err
	}
	return nil
}
*/

func (u *UserRelated) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if u.db == nil {
		err := u.reOpenDB(u.LoginUid)
		if err != nil {
			sdkLog("reOpenDB failed ", u.LoginUid)
			return nil, err
		}
	}
	return u.db.Query(query, args...)
}

/*
func (u *UserRelated) getLocalMaxConSeqFromDB() (int64, error) {
	sdkLog("getLocalMaxConSeqFromDB start")
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	rows, err := u.Query("SELECT seq FROM my_local_data where  user_id=?", u.LoginUid)
	if err != nil {
		sdkLog("Query failed ", err.Error())
		return 0, err
	}
	var seq int
	for rows.Next() {
		err = rows.Scan(&seq)
		if err != nil {
			sdkLog("Scan, failed:", err.Error())
			continue
		}
	}
	sdkLog("getLocalMaxConSeqFromDB, seq: ", seq)
	return int64(seq), nil
}

*/
//1
func (u *UserRelated) getNeedSyncLocalMinSeq() int32 {
	sdkLog("getLocalMaxConSeqFromDB start")
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	rows, err := u.Query("SELECT seq FROM my_local_data where  user_id=?", u.LoginUid)
	if err != nil {
		sdkLog("Query failed ", err.Error())
		return 0
	}
	var seq int32
	for rows.Next() {
		err = rows.Scan(&seq)
		if err != nil {
			sdkLog("Scan, failed:", err.Error())
			continue
		}
	}
	sdkLog("getLocalMaxConSeqFromDB, seq: ", seq)
	return seq
}

//1
func (u *UserRelated) setNeedSyncLocalMinSeq(seq int32) {
	sdkLog("setLocalMaxConSeq start ", seq)
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()

	stmt, err := u.Prepare("replace into my_local_data(user_id, seq) values (?,?)")
	if err != nil {
		sdkLog("set failed", err.Error())
		return
	}
	_, err = stmt.Exec(u.LoginUid, seq)
	if err != nil {
		sdkLog("stmt failed,", err.Error())
	}
}

//1
func (u *UserRelated) replaceIntoUser(info *userInfo) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("replace into `user`(uid, `name`, icon, gender, mobile, birth, email, ex) values(?, ?, ?, ?, ?, ?, ?, ?)")
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

func (u *UserRelated) getAllConversationListModel() (err error, list []*ConversationStruct) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	rows, err := u.Query("SELECT * FROM conversation where latest_msg_send_time!=0 order by  case when is_pinned=1 then 0 else 1 end,max(latest_msg_send_time,draft_timestamp) DESC")
	if err != nil {
		sdkLog("Query failed ", err.Error())
		return err, nil
	}
	u.receiveMessageOptMutex.RLock()
	for rows.Next() {
		c := new(ConversationStruct)
		err = rows.Scan(&c.ConversationID, &c.ConversationType, &c.UserID, &c.GroupID, &c.ShowName,
			&c.FaceURL, &c.RecvMsgOpt, &c.UnreadCount, &c.LatestMsg, &c.LatestMsgSendTime, &c.DraftText, &c.DraftTimestamp, &c.IsPinned)
		if err != nil {
			sdkLog("getAllConversationListModel ,err:", err.Error())
			continue
		} else {
			if v, ok := u.receiveMessageOpt[c.ConversationID]; ok {
				c.RecvMsgOpt = int(v)
			}
			list = append(list, c)
		}
	}
	u.receiveMessageOptMutex.RUnlock()
	return nil, list
}
func (u *UserRelated) getConversationListSplitModel(offset, count int) (err error, list []*ConversationStruct) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	rows, err := u.Query("SELECT * FROM conversation where latest_msg_send_time!=0 order by  case when is_pinned=1 then 0 else 1 end,max(latest_msg_send_time,draft_timestamp) DESC LIMIT ? OFFSET ?", count, offset)
	if err != nil {
		sdkLog("Query failed ", err.Error())
		return err, nil
	}
	u.receiveMessageOptMutex.RLock()
	for rows.Next() {
		c := new(ConversationStruct)
		err = rows.Scan(&c.ConversationID, &c.ConversationType, &c.UserID, &c.GroupID, &c.ShowName,
			&c.FaceURL, &c.RecvMsgOpt, &c.UnreadCount, &c.LatestMsg, &c.LatestMsgSendTime, &c.DraftText, &c.DraftTimestamp, &c.IsPinned)
		if err != nil {
			sdkLog("getAllConversationListModel ,err:", err.Error())
			continue
		} else {
			if v, ok := u.receiveMessageOpt[c.ConversationID]; ok {
				c.RecvMsgOpt = int(v)
			}
			list = append(list, c)
		}
	}
	u.receiveMessageOptMutex.RUnlock()
	return nil, list
}
func convert(nanoSecond int64) string {
	if nanoSecond == 0 {
		return ""
	}
	return time.Unix(0, nanoSecond).Format("2006-01-02_15-04-05")
}

func (u *UserRelated) batchInsertConversationModel(conversations []*ConversationStruct) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	tx, err := u.db.Begin()
	if err != nil {
		sdkLog("start transaction err:", err.Error())
		return err
	}
	stmt, err := u.Prepare("INSERT INTO conversation(conversation_id, conversation_type, " +
		"user_id,group_id,show_name,face_url,recv_msg_opt,unread_count,latest_msg,latest_msg_send_time,draft_text,draft_timestamp,is_pinned) values(?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	defer stmt.Close()
	for _, c := range conversations {
		_, err = stmt.Exec(c.ConversationID, c.ConversationType, c.UserID, c.GroupID, c.ShowName, c.FaceURL, c.RecvMsgOpt, c.UnreadCount, c.LatestMsg, c.LatestMsgSendTime, c.DraftText, c.DraftTimestamp, c.IsPinned)
		if err != nil {
			sdkLog("Exec failed", err.Error(), c)
			continue
		}
	}
	err = tx.Commit()
	if err != nil {
		sdkLog("transaction commit failed", err.Error())
	}
	return nil
}
func (u *UserRelated) getConversationLatestMsgModel(conversationID string) (err error, latestMsg string) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var s string
	rows, err := u.Query("SELECT latest_msg FROM conversation where  conversation_id=?", conversationID)
	if err != nil {
		sdkLog("SELECT latest_msg FROM conversation where  conversation_id=", err.Error())
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

func (u *UserRelated) updateConversationLatestMsgModel(latestMsgSendTime int64, latestMsg, conversationID string) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update conversation set latest_msg=?,latest_msg_send_time=? where conversation_id=?")
	if err != nil {
		sdkLog(err.Error())
		return err
	}

	_, err = stmt.Exec(latestMsg, latestMsgSendTime, conversationID)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil

}
func (u *UserRelated) batchUpdateConversationLatestMsgModel(conversations []*ConversationStruct) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	tx, err := u.db.Begin()
	if err != nil {
		sdkLog("start transaction err:", err.Error())
		return err
	}
	stmt, err := u.Prepare("update conversation set latest_msg=?,latest_msg_send_time=?,unread_count = unread_count+? where conversation_id=?")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	defer stmt.Close()
	for _, c := range conversations {
		_, err = stmt.Exec(c.LatestMsg, c.LatestMsgSendTime, c.UnreadCount, c.ConversationID)
		if err != nil {
			sdkLog(err.Error())
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		sdkLog("transaction commit failed, ", err.Error())
		return err
	}
	return nil

}
func (u *UserRelated) setConversationFaceUrlAndNickName(c *ConversationStruct, conversationID string) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update conversation set show_name=?,face_url=? where conversation_id=?")
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
func (u *UserRelated) judgeConversationIfExists(conversationID string) bool {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var count int
	rows, err := u.Query("select count(*) from conversation where  conversation_id=? And latest_msg_send_time!=0", conversationID)
	if err != nil {
		sdkLog("judge err")
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
func (u *UserRelated) insertConOrUpdateLatestMsg(c *ConversationStruct, conversationID string) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("INSERT INTO conversation(conversation_id, conversation_type, user_id,group_id,show_name,face_url,recv_msg_opt,unread_count,latest_msg,latest_msg_send_time,draft_text,draft_timestamp,is_pinned)" +
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
func (u *UserRelated) getOneConversationModel(conversationID string) (err error, c ConversationStruct) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	rows, err := u.Query("SELECT * FROM conversation where  conversation_id=?", conversationID)
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
	if v, ok := u.receiveMessageOpt[c.ConversationID]; ok {
		c.RecvMsgOpt = int(v)
	}
	return nil, c

}
func (u *UserRelated) deleteConversationModel(conversationID string) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("delete from conversation where conversation_id=?")
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
func (u *UserRelated) ResetConversation(conversationID string) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update  conversation set unread_count=?,latest_msg=?,latest_msg_send_time=?," +
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
func (u *UserRelated) clearConversation(conversationID string) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update  conversation set unread_count=?,latest_msg=?," +
		"draft_text=?,draft_timestamp=? where conversation_id=?")
	if err != nil {
		sdkLog("ResetConversation", err.Error())
		return err
	}
	_, err = stmt.Exec(0, "", "", 0, conversationID)
	if err != nil {
		sdkLog("ResetConversation err:", err.Error())
		return err
	}
	return nil

}
func (u *UserRelated) setConversationDraftModel(conversationID, draftText string) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update conversation set draft_text=?,draft_timestamp=?,latest_msg_send_time=case when latest_msg_send_time=0 then ? else latest_msg_send_time  end where conversation_id=?")
	if err != nil {
		sdkLog("setConversationDraftModel err:", err.Error())
		return err
	}

	_, err = stmt.Exec(draftText, getCurrentTimestampByNano(), getCurrentTimestampByNano(), conversationID)
	if err != nil {
		sdkLog("setConversationDraftModel err:", err.Error())
		return err
	}
	return nil

}
func (u *UserRelated) removeConversationDraftModel(conversationID, draftText string) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update conversation set draft_text=?,draft_timestamp=?where conversation_id=?")
	if err != nil {
		sdkLog("setConversationDraftModel err:", err.Error())
		return err
	}

	_, err = stmt.Exec(draftText, 0, conversationID)
	if err != nil {
		sdkLog("setConversationDraftModel err:", err.Error())
		return err
	}
	return nil

}
func (u *UserRelated) pinConversationModel(conversationID string, isPinned int) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update conversation set is_pinned=?,draft_timestamp=? where conversation_id=?")
	if err != nil {
		sdkLog(err.Error())
		return err
	}

	_, err = stmt.Exec(isPinned, getCurrentTimestampByNano(), conversationID)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil

}
func (u *UserRelated) unPinConversationModel(conversationID string, isPinned int) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update conversation set is_pinned=?,draft_timestamp=case when draft_text='' then ? else draft_timestamp  end where conversation_id=?")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(isPinned, 0, conversationID)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil

}
func (u *UserRelated) setConversationUnreadCount(unreadCount int, conversationID string) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update conversation set unread_count=? where conversation_id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(unreadCount, conversationID)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserRelated) setConversationRecvMsgOpt(conversationID string, opt int) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update conversation set recv_msg_opt=? where conversation_id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(opt, conversationID)
	if err != nil {
		return err
	}
	return nil

}
func (u *UserRelated) incrConversationUnreadCount(conversationID string) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update conversation set unread_count = unread_count+1 where conversation_id=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(conversationID)
	if err != nil {
		return err
	}
	return nil

}
func (u *UserRelated) getTotalUnreadMsgCountModel() (totalUnreadCount int32, err error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	u.receiveMessageOptMutex.RLock()
	var uidList []string
	for key, v := range u.receiveMessageOpt {
		if v == ReceiveNotNotifyMessage {
			uidList = append(uidList, key)
		}
	}
	u.receiveMessageOptMutex.RUnlock()
	rows, err := u.Query("SELECT IFNULL(SUM(unread_count), 0) FROM conversation where recv_msg_opt!=? And conversation_id not in ("+sqlStringHandle(uidList)+")", ReceiveNotNotifyMessage)
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
func (u *UserRelated) setMultipleConversationRecvMsgOpt(conversationIDList []string, opt int) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update conversation set recv_msg_opt=? where conversation_id in (" + sqlStringHandle(conversationIDList) + ")")
	if err != nil {
		sdkLog("setMultipleConversationRecvMsgOpt err:", err.Error(), opt, sqlStringHandle(conversationIDList))
		return err
	}
	_, err = stmt.Exec(opt)
	if err != nil {
		return err
	}
	return nil

}
func (u *UserRelated) getMultipleConversationModel(conversationIDList []string) (err error, list []*ConversationStruct) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	rows, err := u.Query("SELECT * FROM conversation where conversation_id in (" + sqlStringHandle(conversationIDList) + ")")
	sdkLog("SELECT * FROM conversation where conversation_id in (" + sqlStringHandle(conversationIDList) + ")")
	if err != nil {
		sdkLog("getMultipleConversationModel err:", err.Error())
		return err, nil
	}
	for rows.Next() {
		temp := new(ConversationStruct)
		err = rows.Scan(&temp.ConversationID, &temp.ConversationType, &temp.UserID, &temp.GroupID, &temp.ShowName,
			&temp.FaceURL, &temp.RecvMsgOpt, &temp.UnreadCount, &temp.LatestMsg, &temp.LatestMsgSendTime, &temp.DraftText, &temp.DraftTimestamp, &temp.IsPinned)
		if err != nil {
			sdkLog("getMultipleConversationModel err:", err.Error())
			return err, nil
		} else {
			if v, ok := u.receiveMessageOpt[temp.ConversationID]; ok {
				temp.RecvMsgOpt = int(v)
			}
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

//1
func (u *UserRelated) insertIntoTheFriendToFriendInfo(uid, name, comment, icon string, gender int32, mobile, birth, email, ex string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("insert into friend_info(uid,name,comment,icon,gender,mobile,birth,email,ex) values (?,?,?,?,?,?,?,?,?)")
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

//1
func (u *UserRelated) delTheFriendFromFriendInfo(uid string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("delete from friend_info where uid=?")
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

//1
func (u *UserRelated) updateTheFriendInfo(uid, name, comment, icon string, gender int32, mobile, birth, email, ex string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("replace into friend_info(uid,name,comment,icon,gender,mobile,birth,email,ex) values (?,?,?,?,?,?,?,?,?)")
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

//1
func (u *UserRelated) updateFriendInfo(uid, name, icon string, gender int32, mobile, birth, email, ex string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update friend_info set `name` = ?, icon = ?, gender = ?, mobile = ?, birth = ?, email = ?, ex = ? where uid = ?")
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

//1
func (u *UserRelated) insertIntoTheUserToBlackList(info userInfo) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("insert into black_list(uid,name,icon,gender,mobile,birth,email,ex) values (?,?,?,?,?,?,?,?)")
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

//1
func (u *UserRelated) updateBlackList(info userInfo) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("replace into black_list(uid,name,icon,gender,mobile,birth,email,ex) values (?,?,?,?,?,?,?,?)")
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

//1
func (u *UserRelated) delTheUserFromBlackList(uid string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("delete from black_list where uid=?")
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

//1
func (u *UserRelated) insertIntoTheUserToApplicationList(appUserInfo applyUserInfo) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("insert into friend_request(uid,name,icon,gender,mobile,birth,email,ex,flag,req_message,create_time) values (?,?,?,?,?,?,?,?,?,?,?)")
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

//1
func (u *UserRelated) delTheUserFromApplicationList(uid string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("delete from friend_request where uid=?")
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

//1
func (u *UserRelated) updateApplicationList(info applyUserInfo) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("replace into friend_request(uid,name,icon,gender,mobile,birth,email,ex,flag,req_message,create_time) values (?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	_, err = stmt.Exec(info.Uid, info.Name, info.Icon, info.Gender, info.Mobile, info.Birth, info.Email, info.Ex, info.Flag, info.ReqMessage, info.ApplyTime)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	return nil
}

//1
func (u *UserRelated) getFriendInfoByFriendUid(friendUid string) (*friendInfo, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	stmt, err := u.Query("select * from friend_info  where uid=? ", friendUid)
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

//1
func (u *UserRelated) getLocalFriendList22() ([]friendInfo, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	stmt, err := u.Query("select * from friend_info")
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

//1
func (u *UserRelated) getLocalFriendApplication() ([]applyUserInfo, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	stmt, err := u.Query("select * from friend_request order by create_time desc")
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

//1
func (u *UserRelated) getLocalBlackList() ([]userInfo, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	stmt, err := u.Query("select * from black_list")
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

//1
func (u *UserRelated) getBlackUsInfoByUid(blackUid string) (*userInfo, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	stmt, err := u.Query("select * from black_list where uid=?", blackUid)
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

//
//func (u *UserRelated) updateLocalTransferGroupOwner(transfer *TransferGroupOwnerReq) error {
//	u.mRWMutex.Lock()
//	defer u.mRWMutex.Unlock()
//
//	stmt, err := u.Prepare("update group_member set administrator_level = ? where group_id = ? and uid = ?")
//	if err != nil {
//		sdkLog(err.Error())
//		return err
//	}
//	_, err = stmt.Exec(0, transfer.GroupID, transfer.OldOwner)
//	if err != nil {
//		sdkLog(err.Error())
//		return err
//	}
//
//	stmt, err = u.Prepare("update group_member set administrator_level = ? where group_id = ? and uid = ?")
//	if err != nil {
//		sdkLog(err.Error())
//		return err
//	}
//	_, err = stmt.Exec(1, transfer.GroupID, transfer.NewOwner)
//	if err != nil {
//		sdkLog(err.Error())
//		return err
//	}
//
//	return nil
//}
//
//func (u *UserRelated) insertLocalAcceptGroupApplication(addMem *groupMemberFullInfo) error {
//	u.mRWMutex.Lock()
//	defer u.mRWMutex.Unlock()
//	stmt, err := u.Prepare("insert into group_member(group_id,uid,nickname,user_group_face_url,administrator_level,join_time) values (?,?,?,?,?,?)")
//	if err != nil {
//		sdkLog(err.Error())
//		return err
//	}
//	_, err = stmt.Exec(addMem.GroupId, addMem.UserId, addMem.NickName, addMem.FaceUrl, addMem.Role, addMem.JoinTime)
//	if err != nil {
//		sdkLog(err.Error())
//		return err
//	}
//	return nil
//}
func (u *UserRelated) insertIntoLocalGroupInfo(info groupInfo) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("insert into group_info(group_id,name,introduction,notification,face_url,create_time,ex) values (?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.GroupId, info.GroupName, info.Introduction, info.Notification, info.FaceUrl, info.CreateTime, info.Ex)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserRelated) delLocalGroupInfo(groupId string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("delete from group_info where group_id=?")
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
func (u *UserRelated) replaceLocalGroupInfo(info groupInfo) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("replace into group_info(group_id,name,introduction,notification,face_url,create_time,ex) values (?,?,?,?,?,?,?)")
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

//func (u *UserRelated) updateLocalGroupInfo(info groupInfo) error {
//	u.mRWMutex.Lock()
//	defer u.mRWMutex.Unlock()
//	stmt, err := u.Prepare("update group_info set name=?,introduction=?,notification=?,face_url=? where group_id=?")
//	if err != nil {
//		return err
//	}
//	_, err = stmt.Exec(info.GroupName, info.Introduction, info.Notification, info.FaceUrl, info.GroupId)
//	if err != nil {
//		sdkLog(err.Error())
//		return err
//	}
//	return nil
//}

func (u *UserRelated) getLocalGroupsInfo() ([]groupInfo, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	stmt, err := u.Query("select * from group_info")

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

func (u *UserRelated) getLocalGroupsInfoByGroupID(groupID string) (*groupInfo, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	stmt, err := u.Query("select * from group_info where group_id=?", groupID)
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

//
//func (u *UserRelated) findLocalGroupOwnerByGroupId(groupId string) (uid string, err error) {
//	u.mRWMutex.RLock()
//	defer u.mRWMutex.RUnlock()
//	stmt, err := u.Query("select uid from group_member where group_id=? and administrator_level=?", groupId, 1)
//	if err != nil {
//		return "", err
//	}
//	for stmt.Next() {
//		err = stmt.Scan(&uid)
//		if err != nil {
//			sdkLog(err.Error())
//			continue
//		}
//	}
//
//	return uid, nil
//}
//
//func (u *UserRelated) getLocalGroupMemberNumByGroupId(groupId string) (num int, err error) {
//	u.mRWMutex.RLock()
//	defer u.mRWMutex.RUnlock()
//	stmt, err := u.Query("select count(*) from group_member where group_id=?", groupId)
//	if err != nil {
//		return 0, err
//	}
//	for stmt.Next() {
//		err = stmt.Scan(&num)
//		if err != nil {
//			sdkLog("getLocalGroupMemberNumByGroupId query failed, err", err.Error())
//			continue
//		}
//	}
//	return num, err
//}

//1
func (u *UserRelated) getLocalGroupMemberInfoByGroupIdUserId(groupId string, uid string) (*groupMemberFullInfo, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	stmt, err := u.Query("select * from group_member where group_id=? and uid=?", groupId, uid)
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

//1
func (u *UserRelated) getLocalGroupMemberList() ([]groupMemberFullInfo, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	stmt, err := u.Query("select * from group_member")

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

//1
func (u *UserRelated) getLocalGroupMemberListByGroupID(groupId string) ([]groupMemberFullInfo, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	stmt, err := u.Query("select * from group_member where group_id=?", groupId)
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

//1
func (u *UserRelated) insertIntoLocalGroupMember(info groupMemberFullInfo) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("insert into group_member(group_id, uid, nickname,user_group_face_url,administrator_level, join_time) values (?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.GroupId, info.UserId, info.NickName, info.FaceUrl, info.Role, info.JoinTime)
	if err != nil {
		return err
	}
	return nil
}

//1
func (u *UserRelated) delLocalGroupMember(info groupMemberFullInfo) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("delete from group_member where group_id=? and uid=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.GroupId, info.UserId)
	if err != nil {
		return err
	}
	return nil
}

//1
func (u *UserRelated) replaceLocalGroupMemberInfo(info groupMemberFullInfo) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("replace into group_member(group_id,uid,nickname,user_group_face_url,administrator_level, join_time) values (?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.GroupId, info.UserId, info.NickName, info.FaceUrl, info.Role, info.JoinTime)
	if err != nil {
		return err
	}
	return nil
}

//
//func (u *UserRelated) updateLocalGroupMemberInfo(info groupMemberFullInfo) error {
//	u.mRWMutex.Lock()
//	defer u.mRWMutex.Unlock()
//	stmt, err := u.Prepare("update group_member set nickname=?,user_group_face_url=? where group_id=? and uid=?")
//	if err != nil {
//		return err
//	}
//	_, err = stmt.Exec(info.NickName, info.FaceUrl, info.GroupId, info.UserId)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (u *UserRelated) insertIntoSelfApplyToGroupRequest(groupId, message string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("replace into self_apply_to_group_request(group_id,flag,req_msg,create_time) values (?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(groupId, 0, message, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRelated) insertMessageToLocalOrUpdateContent(message *MsgStruct) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("INSERT INTO chat_log(msg_id, send_id, is_read," +
		" seq,status, session_type, recv_id, content_type, sender_face_url,sender_nick_name,msg_from, content, remark,sender_platform_id, send_time,create_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)" +
		"ON CONFLICT(msg_id) DO UPDATE SET content = ?")
	if err != nil {
		sdkLog("failed, ", err.Error())
		return err
	}
	_, err = stmt.Exec(message.ClientMsgID, message.SendID,
		getIsRead(message.IsRead), message.Seq, message.Status, message.SessionType, message.RecvID, message.ContentType, message.SenderFaceURL, message.SenderNickName,
		message.MsgFrom, message.Content, message.Remark, message.SenderPlatformID, message.SendTime, message.CreateTime, message.Content)
	if err != nil {
		sdkLog("failed ", err.Error())
		return err
	}
	return nil
}

func (u *UserRelated) insertMessageToChatLog(message *MsgStruct) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("INSERT INTO chat_log(msg_id, send_id, is_read," +
		" seq,status, session_type, recv_id, content_type, sender_face_url,sender_nick_name,msg_from, content, remark,sender_platform_id, send_time,create_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)" +
		"ON CONFLICT(msg_id) DO UPDATE SET seq = ?")
	if err != nil {
		sdkLog("Prepare failed, ", err.Error())
		return err
	}
	_, err = stmt.Exec(message.ClientMsgID, message.SendID,
		getIsRead(message.IsRead), message.Seq, message.Status, message.SessionType, message.RecvID, message.ContentType, message.SenderFaceURL, message.SenderNickName,
		message.MsgFrom, message.Content, message.Remark, message.SenderPlatformID, message.SendTime, message.CreateTime, message.Seq)
	if err != nil {
		sdkLog("Exec failed, ", err.Error())
		return err
	}
	return nil
}
func (u *UserRelated) batchInsertMessageToChatLog(messages []*InsertMsg) (err error, errMsg []*MsgStruct) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	tx, err := u.db.Begin()
	if err != nil {
		sdkLog("start transaction err:", err.Error())
		return err, nil
	}
	stmt, err := u.Prepare("INSERT INTO chat_log(msg_id, send_id, is_read," +
		" seq,status, session_type, recv_id, content_type, sender_face_url,sender_nick_name,msg_from, content, remark,sender_platform_id, send_time,create_time,is_filter) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)" +
		"ON CONFLICT(msg_id) DO UPDATE SET seq = ?")
	if err != nil {
		sdkLog("Prepare failed, ", err.Error())
		return err, nil
	}
	defer stmt.Close()
	for _, message := range messages {
		_, err = stmt.Exec(message.ClientMsgID, message.SendID,
			getIsRead(message.IsRead), message.Seq, message.Status, message.SessionType, message.RecvID, message.ContentType, message.SenderFaceURL, message.SenderNickName,
			message.MsgFrom, message.Content, message.Remark, message.SenderPlatformID, message.SendTime, message.CreateTime, getIsFilter(message.isFilter), message.Seq)
		if err != nil {
			sdkLog("Exec failed, ", err.Error(), message)
			errMsg = append(errMsg, message.MsgStruct)
			continue
		}
	}
	err = tx.Commit()
	if err != nil {
		sdkLog("transaction commit failed, ", err.Error())
		return err, nil
	}
	return nil, errMsg
}
func (u *UserRelated) batchInsertErrorMessageToErrorChatLog(messages []*MsgStruct) (err error, errMsg []*MsgStruct) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	tx, err := u.db.Begin()
	if err != nil {
		sdkLog("start transaction err:", err.Error())
		return err, nil
	}
	stmt, err := u.Prepare("INSERT INTO error_chat_log(seq,msg_id, send_id, is_read," +
		" status, session_type, recv_id, content_type, sender_face_url,sender_nick_name,msg_from, content, remark,sender_platform_id, send_time,create_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		sdkLog("Prepare failed, ", err.Error())
		return err, nil
	}
	defer stmt.Close()
	for _, message := range messages {
		_, err = stmt.Exec(message.Seq, message.ClientMsgID, message.SendID,
			getIsRead(message.IsRead), message.Status, message.SessionType, message.RecvID, message.ContentType, message.SenderFaceURL, message.SenderNickName,
			message.MsgFrom, message.Content, message.Remark, message.SenderPlatformID, message.SendTime, message.CreateTime)
		if err != nil {
			sdkLog("Exec failed, ", err.Error())
			errMsg = append(errMsg, message)
			continue
		}
	}
	err = tx.Commit()
	if err != nil {
		sdkLog("transaction commit failed, ", err.Error())
		return err, nil
	}
	return nil, errMsg
}

//func (u *UserRelated) updateMessageSeq(message *MsgStruct) (err error) {
//	u.mRWMutex.Lock()
//	defer u.mRWMutex.Unlock()
//	stmt, err := u.Prepare("update chat_log set seq=?,status=? where msg_id=?")
//	if err != nil {
//		sdkLog("Prepare failed, ", err.Error())
//		return err
//	}
//	_, err = stmt.Exec(message.Seq, message.Status, message.ClientMsgID)
//	if err != nil {
//		sdkLog("Exec failed ", err.Error())
//		return err
//	}
//	return nil
//}
func (u *UserRelated) judgeMessageIfExists(msgID string) bool {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var count int
	rows, err := u.Query("select count(*) from chat_log where  msg_id=?", msgID)
	if err != nil {
		sdkLog("Query failed, ", err.Error())
		return false
	}
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			sdkLog("failed ", err.Error())
			continue
		}
	}
	if count == 1 {
		return true
	} else {
		return false
	}
}
func (u *UserRelated) judgeMessageIfExistsBySeq(seq int64) bool {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var count int
	rows, err := u.Query("select count(*) from chat_log where  seq=?", seq)
	if err != nil {
		sdkLog("Query failed, ", err.Error())
		return false
	}
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			sdkLog("failed ", err.Error())
			continue
		}
	}
	if count == 1 {
		return true
	} else {
		return false
	}
}
func (u *UserRelated) getOneMessage(msgID string) (m *MsgStruct, err error) {
	var isRead int
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	// query
	rows, err := u.Query("SELECT * FROM chat_log where msg_id = ?", msgID)
	if err != nil {
		sdkLog("getOneMessage failed", err.Error(), msgID)
		return nil, err
	}
	temp := new(MsgStruct)
	for rows.Next() {
		err = rows.Scan(&temp.ClientMsgID, &temp.SendID, &isRead,
			&temp.Seq, &temp.Status, &temp.SessionType, &temp.RecvID, &temp.ContentType, &temp.SenderFaceURL, &temp.SenderNickName,
			&temp.MsgFrom, &temp.Content, &temp.Remark, &temp.SenderPlatformID, &temp.SendTime, &temp.CreateTime)
		if err != nil {
			sdkLog("getOneMessage,failed", err.Error())
			continue
		}
	}
	if temp.ClientMsgID != "" {
		temp.IsRead = getIsReadB(isRead)
		return temp, nil
	} else {
		return nil, nil
	}
}

func (u *UserRelated) setSingleMessageHasRead(sendID string) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update chat_log set is_read=? where send_id=?And is_read=?AND session_type=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(HasRead, sendID, NotRead, SingleChatType)
	if err != nil {
		return err
	}
	return nil

}
func (u *UserRelated) setSingleMessageHasReadByMsgIDList(sendID string, msgIDList []string) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update chat_log set is_read=? where send_id=?And is_read=?AND session_type=?AND msg_id in(" + sqlStringHandle(msgIDList) + ")")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(HasRead, sendID, NotRead, SingleChatType)
	if err != nil {
		return err
	}
	return nil

}
func (u *UserRelated) setGroupMessageHasRead(groupID string) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update chat_log set is_read=? where recv_id=?And is_read=?AND session_type=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(HasRead, groupID, NotRead, GroupChatType)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserRelated) setMessageStatus(msgID string, status int) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update chat_log set status=? where msg_id=?")
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
func (u *UserRelated) setMessageStatusBySourceID(sourceID string, status, sessionType int) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update chat_log set status=? where (send_id=? or recv_id=?)AND session_type=?")
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
func (u *UserRelated) setMessageHasReadByMsgID(msgID string) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update chat_log set is_read=? where msg_id=?And is_read=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(HasRead, msgID, NotRead)
	if err != nil {
		return err
	}
	return nil

}
func (u *UserRelated) getHistoryMessage(sourceConversationID string, startTime int64, count int, sessionType int) (err error, list MsgFormats) {
	var isRead int
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	rows, err := u.Query("select * from chat_log WHERE (send_id = ? OR recv_id =? )AND (content_type<=? and content_type not in (?)or (content_type >=? and content_type <=?  and content_type not in(?,?)  ))AND status not in(?,?)AND session_type=?AND send_time<?  order by send_time DESC  LIMIT ? OFFSET 0 ",
		sourceConversationID, sourceConversationID, AcceptFriendApplicationTip, HasReadReceipt, GroupTipBegin, GroupTipEnd, SetGroupInfoTip, JoinGroupTip, MsgStatusHasDeleted, MsgStatusRevoked, sessionType, startTime, count)
	if err != nil {
		return err, nil
	}
	for rows.Next() {
		temp := new(MsgStruct)
		err = rows.Scan(&temp.ClientMsgID, &temp.SendID, &isRead, &temp.Seq, &temp.Status, &temp.SessionType,
			&temp.RecvID, &temp.ContentType, &temp.SenderFaceURL, &temp.SenderNickName, &temp.MsgFrom, &temp.Content, &temp.Remark, &temp.SenderPlatformID, &temp.SendTime, &temp.CreateTime)
		if err != nil {
			sdkLog("getHistoryMessage,err:", err.Error())
			continue
		} else {
			err = u.msgHandleByContentType(temp)
			if err != nil {
				sdkLog("getHistoryMessage,err:", err.Error())
				continue
			}
			temp.IsRead = getIsReadB(isRead)
			list = append(list, temp)
		}
	}
	return nil, list
}
func (u *UserRelated) deleteMessageByConversationModel(sourceConversationID string, maxSeq int64) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("delete from chat_log  where send_id=? or recv_id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(sourceConversationID, sourceConversationID)
	if err != nil {
		sdkLog("failed ", err.Error())
		return err
	}
	return nil
}
func (u *UserRelated) deleteMessageByMsgID(msgID string) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("delete from chat_log  where msg_id=?")
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
func (u *UserRelated) updateMessageTimeAndMsgIDStatus(ClientMsgID string, sendTime int64, status int) (err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("update chat_log set send_time=?, status=? where msg_id=? and seq=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(sendTime, status, ClientMsgID, 0)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserRelated) getMultipleMessageModel(messageIDList []string) (err error, list []*MsgStruct) {
	var isRead int
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	rows, err := u.Query("SELECT * FROM chat_log where msg_id in (" + sqlStringHandle(messageIDList) + ")")
	fmt.Println("SELECT * FROM conversation where conversation_id in (" + sqlStringHandle(messageIDList) + ")")
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	for rows.Next() {
		temp := new(MsgStruct)
		err = rows.Scan(&temp.ClientMsgID, &temp.SendID, &isRead, &temp.Seq, &temp.Status, &temp.SessionType,
			&temp.RecvID, &temp.ContentType, &temp.SenderFaceURL, &temp.SenderNickName, &temp.MsgFrom, &temp.Content, &temp.Remark, &temp.SenderPlatformID, &temp.SendTime, &temp.CreateTime)
		if err != nil {
			sdkLog("getMultipleMessageModel,err:", err.Error())
			continue
		} else {
			err = u.msgHandleByContentType(temp)
			if err != nil {
				sdkLog("getMultipleMessageModel,err:", err.Error())
				continue
			}
			temp.IsRead = getIsReadB(isRead)
			list = append(list, temp)
		}
	}
	return nil, list
}

func (u *UserRelated) getErrorChatLogSeq(startSeq int32) map[int32]interface{} {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()

	errSeq := make(map[int32]interface{}, 0)
	var seq int64
	rows, err := u.Query("SELECT seq FROM error_chat_log where seq>=? order by seq", startSeq)
	if err == nil {
		for rows.Next() {
			err = rows.Scan(&seq)
			if err != nil {
				sdkLog("Scan ,failed ", err.Error())
				continue
			} else {
				sdkLog("getErrorChatLogSeq", seq)
				errSeq[int32(seq)] = nil
			}
		}
	} else {
		sdkLog("Query failed ", err.Error())
	}
	LogEnd()
	return errSeq
}

func (u *UserRelated) getNormalChatLogSeq(startSeq int32) map[int32]interface{} {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()

	errSeq := make(map[int32]interface{}, 0)
	var seq int64
	rows, err := u.Query("SELECT seq FROM chat_log where seq>=? order by seq", startSeq)
	if err == nil {
		for rows.Next() {
			err = rows.Scan(&seq)
			if err != nil {
				sdkLog("Scan ,failed ", err.Error())
				continue
			} else {
				//				sdkLog("getNormalChatLogSeq", seq)
				errSeq[int32(seq)] = nil
			}
		}
	} else {
		sdkLog("Query failed ", err.Error())
	}
	LogEnd()
	return errSeq
}

func (u *UserRelated) getConsequentLocalMaxSeq() (seq int64, err error) {
	LogStart()
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()

	errSeq := make(map[int64]interface{}, 1)

	old := u.GetMinSeqSvr()
	var rSeq int64
	var rows *sql.Rows
	if old == 0 {
		rows, err = u.Query("SELECT seq FROM error_chat_log where seq>? order by seq", old)
		if err == nil {
			for rows.Next() {
				err = rows.Scan(&seq)
				if err != nil {
					sdkLog("Scan ,failed ", err.Error())
					continue
				} else {
					sdkLog("err seq in map: ", seq)
					errSeq[seq] = nil
				}
			}
		} else {
			sdkLog("Query failed ", err.Error())
		}

		rows, err = u.Query("SELECT seq FROM chat_log where seq>? order by seq", old)
		if err != nil {
			sdkLog("getLocalMaxSeqModel,Query  failed", err.Error(), old)
			LogFReturn(old, err)
			return old, err
		}
		var idx int64 = 0
		rSeq = old
		for rows.Next() {
			err = rows.Scan(&seq)
			if err != nil {
				sdkLog("getLocalMaxSeqModel ,failed ", err.Error())
				continue
			} else {
				idx++
				if seq == old+idx {
					rSeq = seq
				} else {
					_, ok := errSeq[old+idx]
					if ok {
						rSeq = old + idx
						sdkLog("seq in err map ", old+idx)
					} else {
						sdkLog("not consequent ", old, idx, seq)
						rows.Close()
						break
					}
				}
			}
		}
		LogSReturn(rSeq, nil)
		return rSeq, nil
	} else {
		rows, err = u.Query("SELECT seq FROM error_chat_log where seq>=? order by seq", old)
		if err == nil {
			for rows.Next() {
				err = rows.Scan(&seq)
				if err != nil {
					sdkLog("Scan ,failed ", err.Error())
					continue
				} else {
					errSeq[seq] = nil
				}
			}
		} else {
			sdkLog("Query failed ", err.Error())
		}

		rows, err = u.Query("SELECT seq FROM chat_log where seq>=? order by seq", old)
		if err != nil {
			sdkLog("getLocalMaxSeqModel,Query err:", err.Error(), old)
			LogFReturn(old, err)
			return old, err
		}
		var idx int64 = 0
		rSeq = old
		for rows.Next() {
			err = rows.Scan(&seq)
			if err != nil {
				sdkLog("getLocalMaxSeqModel ,err:", err.Error())
				continue
			} else {
				if seq == old+idx {
					rSeq = seq
					idx++
				} else {
					_, ok := errSeq[old+idx]
					if ok {
						rSeq = old + idx
						sdkLog("seq in err map ", old+idx)
						idx++
					} else {
						sdkLog("not consequent ", old, idx, seq)
						rows.Close()
						break
					}
				}
			}
		}
		LogSReturn(rSeq, nil)
		return rSeq, nil
	}
}

func (u *UserRelated) isExistsInErrChatLogBySeq(seq int64) bool {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var count int
	rows, err := u.Query("select count(*) from error_chat_log where  seq=?", seq)
	if err != nil {
		sdkLog("Query failed, ", err.Error())
		return false
	}
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			sdkLog("failed ", err.Error())
			continue
		}
	}
	if count == 1 {
		return true
	} else {
		return false
	}
}

//1
func (ur *UserRelated) getLoginUserInfoFromLocal() (userInfo, error) {
	ur.mRWMutex.RLock()
	defer ur.mRWMutex.RUnlock()
	var u userInfo
	rows, err := ur.Query("select * from user limit 1 ")
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
		sdkLog("db Query failed, ", err.Error())
		return u, err
	}
}

func (u *UserRelated) getOwnLocalGroupApplicationList(groupId string) (*groupApplicationResult, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()

	sql := "select id, group_id, from_user_id, to_user_id, flag, req_msg, handled_msg, create_time, from_user_nickname, " +
		" to_user_nickname, from_user_face_url, to_user_face_url, handled_user  from `group_request` "

	if len(groupId) > 0 {
		sql = fmt.Sprintf("%s where group_id = %s", sql, groupId)
	}

	rows, err := u.Query(sql)
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
			if glInfo.HandledUser == u.LoginUid {
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

	return &reply, nil
}

func (u *UserRelated) insertIntoRequestToGroupRequest(info GroupReqListInfo) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("insert into group_request(id, group_id, from_user_id, to_user_id,flag,req_msg,handled_msg, create_time,from_user_nickname,to_user_nickname,from_user_face_url,to_user_face_url,handled_user) values (?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.ID, info.GroupID, info.FromUserID, info.ToUserID, info.Flag, info.RequestMsg, info.HandledMsg, info.AddTime, info.FromUserNickname, info.ToUserNickname, info.FromUserFaceUrl, info.ToUserFaceUrl, info.HandledUser)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRelated) delRequestFromGroupRequest(info GroupReqListInfo) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("delete from group_request where id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRelated) replaceIntoRequestToGroupRequest(info GroupReqListInfo) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	stmt, err := u.Prepare("replace into group_request(id,group_id, from_user_id, to_user_id,flag,req_msg,handled_msg, create_time,from_user_nickname,to_user_nickname,from_user_face_url,to_user_face_url,handled_user) values (?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(info.ID, info.GroupID, info.FromUserID, info.ToUserID, info.Flag, info.RequestMsg, info.HandledMsg, info.AddTime, info.FromUserNickname, info.ToUserNickname, info.FromUserFaceUrl, info.ToUserFaceUrl, info.HandledUser)
	if err != nil {
		return err
	}
	return nil
}
