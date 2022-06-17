package db

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"sync"
	"time"
)

var UserDBMap map[string]*DataBase

var UserDBLock sync.RWMutex

func init() {
	UserDBMap = make(map[string]*DataBase, 0)
}

type DataBase struct {
	loginUserID string
	dbDir       string
	conn        *gorm.DB
	mRWMutex    sync.RWMutex
}

//func (d *DataBase) CloseDB() error {
//	d.mRWMutex.Lock()
//	defer d.mRWMutex.Unlock()
//	if d.conn != nil {
//
//		if err := d.conn.Close(); err != nil {
//			log.Error("", "GetSendingMessageList failed ", err.Error())
//			return err
//		}
//	}
//	return nil
//}

func NewDataBase(loginUserID string, dbDir string) (*DataBase, error) {
	UserDBLock.Lock()
	defer UserDBLock.Unlock()
	dataBase, ok := UserDBMap[loginUserID]
	if !ok {
		dataBase = &DataBase{loginUserID: loginUserID, dbDir: dbDir}
		err := dataBase.initDB()
		if err != nil {
			return dataBase, utils.Wrap(err, "initDB failed")
		}
		UserDBMap[loginUserID] = dataBase
		log.Info("", "open db", loginUserID)
	} else {
		log.Info("", "db in map", loginUserID)
	}
	dataBase.setChatLogFailedStatus()
	return dataBase, nil
}

func (d *DataBase) setChatLogFailedStatus() {
	msgList, err := d.GetSendingMessageList()
	if err != nil {
		log.Error("", "GetSendingMessageList failed ", err.Error())
		return
	}
	for _, v := range msgList {
		v.Status = constant.MsgStatusSendFailed
		err := d.UpdateMessage(v)
		if err != nil {
			log.Error("", "UpdateMessage failed ", err.Error(), v)
			continue
		}
	}
}

func (d *DataBase) initDB() error {
	if d.loginUserID == "" {
		return errors.New("no uid")
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()

	dbFileName := d.dbDir + "/OpenIM_" + constant.BigVersion + "_" + d.loginUserID + ".db"
	//db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
	//  Logger: logger.Default.LogMode(logger.Silent),
	//})
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	log.Info("open db:", dbFileName)
	if err != nil {
		return utils.Wrap(err, "open db failed")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return utils.Wrap(err, "get sql db failed")
	}
	sqlDB.SetConnMaxLifetime(time.Hour * 1)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(2)
	d.conn = db
	//db, err := sql.Open("sqlite3", SvrConf.DbDir+"OpenIM_"+uid+".db")
	//sdkLog("open db:", SvrConf.DbDir+"OpenIM_"+uid+".db")
	//if err != nil {
	//	sdkLog("failed open db:", SvrConf.DbDir+"OpenIM_"+uid+".db", err.Error())
	//	return err
	//}
	//u.db = db

	superGroup := &model_struct.LocalGroup{}
	localGroup := &model_struct.LocalGroup{}
	//err = db.Exec(`CREATE TABLE IF NOT EXISTS super_groups (
	//   group_id  varchar(64),
	//   name  text,
	//   notification  varchar(255),
	//   introduction  varchar(255),
	//   face_url  varchar(255),
	//   create_time  integer,
	//   status  integer,
	//   creator_user_id  varchar(64),
	//   group_type  integer,
	//   owner_user_id  varchar(64),
	//   member_count  integer,
	//   ex  varchar(1024),
	//   attached_info  varchar(1024),
	//   PRIMARY KEY ( group_id ))`).Error
	//if err != nil {
	//	log.Error("super_group","create super group failed",err.Error())
	//}

	db.AutoMigrate(&model_struct.LocalFriend{},
		&model_struct.LocalFriendRequest{},
		localGroup,
		&model_struct.LocalGroupMember{},
		&model_struct.LocalGroupRequest{},
		&model_struct.LocalErrChatLog{},
		&model_struct.LocalUser{},
		&model_struct.LocalBlack{},
		&model_struct.LocalSeqData{},
		&model_struct.LocalConversation{},
		&model_struct.LocalChatLog{},
		&model_struct.LocalAdminGroupRequest{},
		&model_struct.LocalDepartment{},
		&model_struct.LocalDepartmentMember{},
		&LocalWorkMomentsNotification{},
		&LocalWorkMomentsNotificationUnreadCount{},
	)
	db.Table(constant.SuperGroupTableName).AutoMigrate(superGroup)
	if !db.Migrator().HasTable(&model_struct.LocalFriend{}) {
		//log.NewInfo("CreateTable Friend")
		db.Migrator().CreateTable(&model_struct.LocalFriend{})
	}

	if !db.Migrator().HasTable(&model_struct.LocalFriendRequest{}) {
		//log.NewInfo("CreateTable FriendRequest")
		db.Migrator().CreateTable(&model_struct.LocalFriendRequest{})
	}

	if !db.Migrator().HasTable(localGroup) {
		//log.NewInfo("CreateTable Group")
		db.Migrator().CreateTable(localGroup)
	}
	if !db.Migrator().HasTable(&model_struct.LocalGroupMember{}) {
		//log.NewInfo("CreateTable GroupMember")
		db.Migrator().CreateTable(&model_struct.LocalGroupMember{})
	}

	if !db.Migrator().HasTable(&model_struct.LocalGroupRequest{}) {
		//log.NewInfo("CreateTable GroupRequest")
		db.Migrator().CreateTable(&model_struct.LocalGroupRequest{})
	}

	if !db.Migrator().HasTable(&model_struct.LocalUser{}) {
		//log.NewInfo("CreateTable User")
		db.Migrator().CreateTable(&model_struct.LocalUser{})
	}

	if !db.Migrator().HasTable(&model_struct.LocalBlack{}) {
		//log.NewInfo("CreateTable Black")
		db.Migrator().CreateTable(&model_struct.LocalBlack{})
	}

	if !db.Migrator().HasTable(&model_struct.LocalSeqData{}) {
		db.Migrator().CreateTable(&model_struct.LocalSeqData{})
	}
	if !db.Migrator().HasTable(&model_struct.LocalConversation{}) {
		db.Migrator().CreateTable(&model_struct.LocalConversation{})
	}
	if !db.Migrator().HasTable(&model_struct.LocalChatLog{}) {
		db.Migrator().CreateTable(&model_struct.LocalChatLog{})
	}
	if !db.Migrator().HasTable(&model_struct.LocalAdminGroupRequest{}) {
		db.Migrator().CreateTable(&model_struct.LocalAdminGroupRequest{})
	}
	if !db.Migrator().HasTable(&model_struct.LocalDepartment{}) {
		db.Migrator().CreateTable(&model_struct.LocalDepartment{})
	}
	if !db.Migrator().HasTable(&model_struct.LocalDepartmentMember{}) {
		db.Migrator().CreateTable(&model_struct.LocalDepartmentMember{})
	}
	if !db.Migrator().HasTable(&LocalWorkMomentsNotification{}) {
		db.Migrator().CreateTable(&LocalWorkMomentsNotification{})
	}
	if !db.Migrator().HasTable(&LocalWorkMomentsNotificationUnreadCount{}) {
		db.Migrator().CreateTable(&LocalWorkMomentsNotificationUnreadCount{})
	}
	log.NewInfo("init db", "startInitWorkMomentsNotificationUnreadCount ")
	if err := d.InitWorkMomentsNotificationUnreadCount(); err != nil {
		log.NewError("init InitWorkMomentsNotificationUnreadCount:", utils.GetSelfFuncName(), err.Error())
	}
	return nil
}
