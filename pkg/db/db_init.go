package db

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"sync"
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
	}
	log.Info("", "db in map", loginUserID)
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

	d.conn = db
	//db, err := sql.Open("sqlite3", SvrConf.DbDir+"OpenIM_"+uid+".db")
	//sdkLog("open db:", SvrConf.DbDir+"OpenIM_"+uid+".db")
	//if err != nil {
	//	sdkLog("failed open db:", SvrConf.DbDir+"OpenIM_"+uid+".db", err.Error())
	//	return err
	//}
	//u.db = db

	db.AutoMigrate(&LocalFriend{},
		&LocalFriendRequest{},
		&LocalGroup{},
		&LocalGroupMember{},
		&LocalGroupRequest{},
		&LocalErrChatLog{},
		&LocalUser{},
		&LocalBlack{},
		&LocalSeqData{},
		&LocalConversation{},
		&LocalChatLog{},
		&LocalAdminGroupRequest{},
		&LocalDepartment{},
		&LocalDepartmentMember{},
		&LocalWorkMomentsNotification{},
		&LocalWorkMomentsNotificationUnreadCount{},
	)
	if !db.Migrator().HasTable(&LocalFriend{}) {
		//log.NewInfo("CreateTable Friend")
		db.Migrator().CreateTable(&LocalFriend{})
	}

	if !db.Migrator().HasTable(&LocalFriendRequest{}) {
		//log.NewInfo("CreateTable FriendRequest")
		db.Migrator().CreateTable(&LocalFriendRequest{})
	}

	if !db.Migrator().HasTable(&LocalGroup{}) {
		//log.NewInfo("CreateTable Group")
		db.Migrator().CreateTable(&LocalGroup{})
	}

	if !db.Migrator().HasTable(&LocalGroupMember{}) {
		//log.NewInfo("CreateTable GroupMember")
		db.Migrator().CreateTable(&LocalGroupMember{})
	}

	if !db.Migrator().HasTable(&LocalGroupRequest{}) {
		//log.NewInfo("CreateTable GroupRequest")
		db.Migrator().CreateTable(&LocalGroupRequest{})
	}

	if !db.Migrator().HasTable(&LocalUser{}) {
		//log.NewInfo("CreateTable User")
		db.Migrator().CreateTable(&LocalUser{})
	}

	if !db.Migrator().HasTable(&LocalBlack{}) {
		//log.NewInfo("CreateTable Black")
		db.Migrator().CreateTable(&LocalBlack{})
	}

	if !db.Migrator().HasTable(&LocalSeqData{}) {
		db.Migrator().CreateTable(&LocalSeqData{})
	}
	if !db.Migrator().HasTable(&LocalConversation{}) {
		db.Migrator().CreateTable(&LocalConversation{})
	}
	if !db.Migrator().HasTable(&LocalChatLog{}) {
		db.Migrator().CreateTable(&LocalChatLog{})
	}
	if !db.Migrator().HasTable(&LocalAdminGroupRequest{}) {
		db.Migrator().CreateTable(&LocalAdminGroupRequest{})
	}
	if !db.Migrator().HasTable(&LocalDepartment{}) {
		db.Migrator().CreateTable(&LocalDepartment{})
	}
	if !db.Migrator().HasTable(&LocalDepartmentMember{}) {
		db.Migrator().CreateTable(&LocalDepartmentMember{})
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
