package open_im_sdk

import (
	"database/sql"
	"errors"
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
	table := `CREATE TABLE if not exists users(
        user_id char(64) PRIMARY KEY NOT NULL , 
		name varchar(64) DEFAULT NULL , 
		face_url varchar(100) DEFAULT NULL , 
		gender int(11) DEFAULT NULL , 
		phone_number varchar(32) DEFAULT NULL , 
		birth INTEGER DEFAULT NULL , 
		email varchar(64) DEFAULT NULL , 
		create_time INTEGER DEFAULT NULL , 
		app_manger_level int DEFAULT NULL , 
		ex varchar(1024) DEFAULT NULL, 
         )`
	_, err = db.Exec(table)
	if err != nil {
		sdkLog(err.Error())
		return err
	}

	table = `create table if not exists  black_list (
       owner_user_id CHAR (64) NOT NULL,
       block_user_id CHAR (64) NOT NULL ,
	   create_time INTEGER DEFAULT NULL , 
       add_source int DEFAULT NULL , 
       operator_user_id CHAR(64) DEFAULT NULL,
       ex varchar(1024) DEFAULT NULL,
       PRIMARY KEY (owner_user_id,block_user_id)
       )`
	_, err = db.Exec(table)
	if err != nil {
		sdkLog(err.Error())
		return err
	}

	table = `create table if not exists friend_requests (
    	from_user_id CHAR (64) NOT NULL,
    	to_user_id  CHAR(64) NoT NULL ,
     	handle_result int DEFAULT NULL , 
     	req_msg varchar(255) DEFAULT NULL , 
	    create_time INTEGER DEFAULT NULL ,
	    handle_user_id  CHAR(64) NOT NULL,
	    handle_msg  VARCHAR(255) DEFAULT NULL,
	    handle_time INTEGER DEFAULT NULL ,
        ex varchar(1024) DEFAULT NULL,
        PRIMARY KEY (from_user_id,to_user_id)
      )`
	_, err = db.Exec(table)
	if err != nil {
		sdkLog(err.Error())
		return err
	}

	table = ` CREATE TABLE IF NOT EXISTS friends(
     owner_user_id CHAR (64) NOT NULL,
     friend_user_id CHAR (64) NOT NULL ,
     name varchar(64) DEFAULT NULL , 
	 face_url varchar(100) DEFAULT NULL ,
     remark varchar(255) DEFAULT NULL,
     gender int DEFAULT NULL , 
   	 phone_number varchar(32) DEFAULT NULL , 
	 birth INTEGER DEFAULT NULL , 
	 email varchar(64) DEFAULT NULL ,
	 create_time INTEGER DEFAULT NULL ,
	 add_source int DEFAULT NULL ,
	 operator_user_id CHAR(64) DEFAULT NULL,
  	 ex varchar(1024) DEFAULT NULL,
  	 PRIMARY KEY (owner_user_id,friend_user_id)
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

	table = `create table if not exists  groups (
    	group_id char(64) NOT NULL,
		name varchar(64) DEFAULT NULL , 
    	introduction varchar(255) DEFAULT NULL,
    	notification varchar(255) DEFAULT NULL,
    	face_url varchar(100) DEFAULT NULL,
    	group_type int DEFAULT NULL,
    	status int DEFAULT NULL,
    	creator_user_id char(64) DEFAULT NULL,
    	create_time INTEGER DEFAULT NULL,
    	ex varchar(1024) DEFAULT NULL,
    	PRIMARY KEY (group_id)
	)`

	_, err = db.Exec(table)
	if err != nil {
		sdkLog(err.Error())
		return err
	}

	table = `create table if not exists  group_members (
   group_id char(64) NOT NULL,
   user_id char(64) NOT NULL,
   nickname varchar(64) DEFAULT NULL,
   user_group_face_url varchar(64) DEFAULT NULL,
   role_level int DEFAULT NULL,
   join_time INTEGER DEFAULT NULL,
   join_source int DEFAULT NULL,
   operator_user_id char(64) NOT NULL,
   PRIMARY KEY (group_id,user_id)
)`

	_, err = db.Exec(table)
	if err != nil {
		sdkLog(err.Error())
		return err
	}

	table = `create table if not exists  group_requests (
   group_id char(64) NOT NULL,
   user_id char(64) NOT NULL,
req_msg varchar(255) DEFAULT NULL,
handle_msg  varchar(255) DEFAULT NULL,
req_time INTEGER DEFAULT NULL,
handle_user_id char(64) NOT NULL,
handle_time INTEGER DEFAULT NULL,
handle_result int DEFAULT NULL,
    	ex varchar(1024) DEFAULT NULL,
     PRIMARY KEY (group_id,user_id)
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
