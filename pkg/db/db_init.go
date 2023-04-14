//go:build !js
// +build !js

package db

import (
	"context"
	"errors"
	"fmt"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//"github.com/glebarez/sqlite"

var UserDBMap map[string]*DataBase

var UserDBLock sync.RWMutex

func init() {
	UserDBMap = make(map[string]*DataBase, 0)
}

type DataBase struct {
	loginUserID   string
	dbDir         string
	conn          *gorm.DB
	mRWMutex      sync.RWMutex
	groupMtx      sync.RWMutex
	friendMtx     sync.RWMutex
	departmentMtx sync.RWMutex
	userMtx       sync.RWMutex
	superGroupMtx sync.RWMutex
}

func (d *DataBase) InitSuperLocalErrChatLog(ctx context.Context, groupID string) {
	panic("implement me")
}

func (d *DataBase) InitSuperLocalChatLog(ctx context.Context, groupID string) {
	panic("implement me")
}

func (d *DataBase) SetChatLogFailedStatus(ctx context.Context) {
	panic("implement me")
}

func (d *DataBase) InitDB(ctx context.Context, userID string, dataDir string) error {
	panic("implement me")
}

func (d *DataBase) Close(ctx context.Context) error {
	UserDBLock.Lock()
	dbConn, err := d.conn.WithContext(ctx).DB()
	if err != nil {
		log.Error("", "get db conn failed ", err.Error())
	} else {
		if dbConn != nil {
			log.Info("", "close db finished")
			err := dbConn.Close()
			if err != nil {
				log.Error("", "close db failed ", err.Error())
			}
		}
	}

	log.NewInfo("", "CloseDB ok, delete db map ", d.loginUserID)
	delete(UserDBMap, d.loginUserID)
	UserDBLock.Unlock()
	return nil
}

func NewDataBase(ctx context.Context, loginUserID string, dbDir string) (*DataBase, error) {
	UserDBLock.Lock()
	defer UserDBLock.Unlock()

	dataBase, ok := UserDBMap[loginUserID]
	if !ok {
		dataBase = &DataBase{loginUserID: loginUserID, dbDir: dbDir}
		err := dataBase.initDB(ctx)
		if err != nil {
			return dataBase, utils.Wrap(err, "initDB failed "+dbDir)
		}
		UserDBMap[loginUserID] = dataBase
		//log.Info(operationID, "open db", loginUserID)
	} else {
		//log.Info(operationID, "db in map", loginUserID)
	}
	dataBase.setChatLogFailedStatus(ctx)
	return dataBase, nil
}

func (d *DataBase) setChatLogFailedStatus(ctx context.Context) {
	msgList, err := d.GetSendingMessageList(ctx)
	if err != nil {
		log.Error("", "GetSendingMessageList failed ", err.Error())
		return
	}
	for _, v := range msgList {
		v.Status = constant.MsgStatusSendFailed
		err := d.UpdateMessage(ctx, v)
		if err != nil {
			log.Error("", "UpdateMessage failed ", err.Error(), v)
			continue
		}
	}
	groupIDList, err := d.GetReadDiffusionGroupIDList(ctx)
	if err != nil {
		log.Error("", "GetReadDiffusionGroupIDList failed ", err.Error())
		return
	}
	for _, v := range groupIDList {
		msgList, err := d.SuperGroupGetSendingMessageList(ctx, v)
		if err != nil {
			log.Error("", "GetSendingMessageList failed ", err.Error())
			return
		}
		if len(msgList) > 0 {
			for _, v := range msgList {
				v.Status = constant.MsgStatusSendFailed
				err := d.SuperGroupUpdateMessage(ctx, v)
				if err != nil {
					log.Error("", "UpdateMessage failed ", err.Error(), v)
					continue
				}
			}
		}

	}

}

func (d *DataBase) initDB(ctx context.Context) error {
	if d.loginUserID == "" {
		return errors.New("no uid")
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()

	//cxn := "memdb1?mode=memory&cache=shared"
	dbFileName := d.dbDir + "/OpenIM_" + constant.BigVersion + "_" + d.loginUserID + ".db"

	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	log.Info("open db:", dbFileName)
	if err != nil {
		return utils.Wrap(err, "open db failed "+dbFileName)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return utils.Wrap(err, "get sql db failed")
	}
	sqlDB.SetConnMaxLifetime(time.Hour * 1)
	sqlDB.SetMaxOpenConns(3)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)
	d.conn = db

	superGroup := &model_struct.LocalGroup{}
	localGroup := &model_struct.LocalGroup{}

	db.AutoMigrate(&model_struct.LocalFriend{},
		&model_struct.LocalFriendRequest{},
		localGroup,
		&model_struct.LocalGroupMember{},
		&model_struct.LocalGroupRequest{},
		&model_struct.LocalErrChatLog{},
		&model_struct.LocalUser{},
		&model_struct.LocalBlack{},
		&model_struct.LocalConversationUnreadMessage{},
		//&model_struct.LocalSeqData{},
		//&model_struct.LocalSeq{},
		&model_struct.LocalConversation{},
		&model_struct.LocalChatLog{},
		&model_struct.LocalAdminGroupRequest{},
		&model_struct.LocalWorkMomentsNotification{},
		&model_struct.LocalWorkMomentsNotificationUnreadCount{},
		&model_struct.TempCacheLocalChatLog{},
		&model_struct.LocalChatLogReactionExtensions{},
	)
	db.Table(constant.SuperGroupTableName).AutoMigrate(superGroup)
	groupIDList, err := d.GetJoinedSuperGroupIDList(ctx)
	if err != nil {
		log.Error("auto migrate super db err:", err.Error())
	}
	wkGroupIDList, err2 := d.GetJoinedWorkingGroupIDList(ctx)
	if err2 != nil {
		log.Error("auto migrate working group  db err:", err2)

	}
	groupIDList = append(groupIDList, wkGroupIDList...)
	for _, v := range groupIDList {
		d.conn.WithContext(ctx).Table(utils.GetSuperGroupTableName(v)).AutoMigrate(&model_struct.LocalChatLog{})
		var count int64
		_ = db.Raw(fmt.Sprintf("SELECT COUNT(*) AS count FROM sqlite_master WHERE type = 'index' AND name ='%s' AND tbl_name = '%s'",
			"index_seq_"+v, utils.GetSuperGroupTableName(v))).Row().Scan(&count)
		if count == 0 {
			result := db.Exec(fmt.Sprintf("CREATE INDEX %s ON %s (seq)", "index_seq_"+v, utils.GetSuperGroupTableName(v)))
			if result.Error != nil {
				log.Error("create super group index failed:", result.Error.Error(), v)
			}
		}
		var count2 int64
		_ = db.Raw(fmt.Sprintf("SELECT COUNT(*) AS count FROM sqlite_master WHERE type = 'index' AND name ='%s' AND tbl_name = '%s'",
			"index_send_time_"+v, utils.GetSuperGroupTableName(v))).Row().Scan(&count)
		if count2 == 0 {
			result := db.Exec(fmt.Sprintf("CREATE INDEX %s ON %s (send_time)", "index_send_time_"+v, utils.GetSuperGroupTableName(v)))
			if result.Error != nil {
				log.Error("create super group index failed:", result.Error.Error(), v)
			}
		}

	}
	if !db.Migrator().HasTable(&model_struct.LocalFriend{}) {
		db.Migrator().CreateTable(&model_struct.LocalFriend{})
	}

	if !db.Migrator().HasTable(&model_struct.LocalFriendRequest{}) {
		db.Migrator().CreateTable(&model_struct.LocalFriendRequest{})
	}
	if !db.Migrator().HasTable(&model_struct.LocalConversationUnreadMessage{}) {
		db.Migrator().CreateTable(&model_struct.LocalConversationUnreadMessage{})
	}
	if !db.Migrator().HasTable(localGroup) {
		db.Migrator().CreateTable(localGroup)
	}
	if !db.Migrator().HasTable(&model_struct.LocalGroupMember{}) {
		db.Migrator().CreateTable(&model_struct.LocalGroupMember{})
	}

	if !db.Migrator().HasTable(&model_struct.LocalGroupRequest{}) {
		db.Migrator().CreateTable(&model_struct.LocalGroupRequest{})
	}

	if !db.Migrator().HasTable(&model_struct.LocalUser{}) {
		db.Migrator().CreateTable(&model_struct.LocalUser{})
	}

	if !db.Migrator().HasTable(&model_struct.LocalBlack{}) {
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
	if !db.Migrator().HasTable(&model_struct.LocalWorkMomentsNotification{}) {
		db.Migrator().CreateTable(&model_struct.LocalWorkMomentsNotification{})
	}
	if !db.Migrator().HasTable(&model_struct.LocalWorkMomentsNotificationUnreadCount{}) {
		db.Migrator().CreateTable(&model_struct.LocalWorkMomentsNotificationUnreadCount{})
	}
	if !db.Migrator().HasTable(&model_struct.LocalChatLogReactionExtensions{}) {
		db.Migrator().CreateTable(&model_struct.LocalChatLogReactionExtensions{})
	}
	log.NewInfo("init db", "startInitWorkMomentsNotificationUnreadCount ")
	if err := d.InitWorkMomentsNotificationUnreadCount(ctx); err != nil {
		log.NewError("init InitWorkMomentsNotificationUnreadCount:", utils.GetSelfFuncName(), err.Error())
	}
	return nil
}
