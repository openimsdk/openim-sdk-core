// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

//go:build !js
// +build !js

package db

import (
	"context"
	"errors"
	"fmt"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"sync"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
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
		log.ZError(ctx, "GetSendingMessageList failed", err)
		return
	}
	for _, v := range msgList {
		v.Status = constant.MsgStatusSendFailed
		err := d.UpdateMessage(ctx, v)
		if err != nil {
			log.ZError(ctx, "UpdateMessage failed", err, "msg", v)
			continue
		}
	}
	groupIDList, err := d.GetReadDiffusionGroupIDList(ctx)
	if err != nil {
		log.ZError(ctx, "GetReadDiffusionGroupIDList failed", err)
		return
	}
	for _, v := range groupIDList {
		msgList, err := d.SuperGroupGetSendingMessageList(ctx, v)
		if err != nil {
			log.ZError(ctx, "GetSendingMessageList failed", err)
			return
		}
		if len(msgList) > 0 {
			for _, v := range msgList {
				v.Status = constant.MsgStatusSendFailed
				err := d.SuperGroupUpdateMessage(ctx, v)
				if err != nil {
					log.ZError(ctx, "UpdateMessage failed", err, "msg", v)
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

	dbFileName := d.dbDir + "/OpenIM_" + constant.BigVersion + "_" + d.loginUserID + ".db"
	// slowThreshold := 500
	// sqlLogger := log.NewSqlLogger(logger.LogLevel(sdk_struct.SvrConf.LogLevel), true, time.Duration(slowThreshold)*time.Millisecond)
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return utils.Wrap(err, "open db failed "+dbFileName)
	}
	log.ZDebug(ctx, "open db success", "db", db, "dbFileName", dbFileName)
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
		log.ZError(ctx, "auto migrate super db err:", err)
	}
	wkGroupIDList, err := d.GetJoinedWorkingGroupIDList(ctx)
	if err != nil {
		log.ZError(ctx, "auto migrate working group  db err:", err)
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
	if err := d.InitWorkMomentsNotificationUnreadCount(ctx); err != nil {
		log.ZError(ctx, "init InitWorkMomentsNotificationUnreadCount failed", err)
	}
	return nil
}
