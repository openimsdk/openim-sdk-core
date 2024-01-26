// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !js
// +build !js

package db

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"

	"github.com/OpenIMSDK/tools/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DataBase struct {
	loginUserID   string
	dbDir         string
	conn          *gorm.DB
	mRWMutex      sync.RWMutex
	groupMtx      sync.RWMutex
	friendMtx     sync.RWMutex
	userMtx       sync.RWMutex
	superGroupMtx sync.RWMutex
}

func (d *DataBase) GetMultipleMessageReactionExtension(ctx context.Context, msgIDList []string) (result []*model_struct.LocalChatLogReactionExtensions, err error) {
	//TODO implement me
	panic("implement me")
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
	dbConn, err := d.conn.WithContext(ctx).DB()
	if err != nil {
		return err
	} else {
		if dbConn != nil {
			err := dbConn.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewDataBase(ctx context.Context, loginUserID string, dbDir string, logLevel int) (*DataBase, error) {
	dataBase := &DataBase{loginUserID: loginUserID, dbDir: dbDir}
	err := dataBase.initDB(ctx, logLevel)
	if err != nil {
		return dataBase, utils.Wrap(err, "initDB failed "+dbDir)
	}
	return dataBase, nil
}

func (d *DataBase) initDB(ctx context.Context, logLevel int) error {
	var zLogLevel logger.LogLevel
	if d.loginUserID == "" {
		return errors.New("no uid")
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()

	path := d.dbDir + "/OpenIM_" + constant.BigVersion + "_" + d.loginUserID + ".db"
	dbFileName, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	log.ZInfo(ctx, "sqlite", "path", dbFileName)
	// slowThreshold := 500
	// sqlLogger := log.NewSqlLogger(logger.LogLevel(sdk_struct.SvrConf.LogLevel), true, time.Duration(slowThreshold)*time.Millisecond)
	if logLevel > 5 {
		zLogLevel = logger.Info
	} else {
		zLogLevel = logger.Silent
	}
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{Logger: log.NewSqlLogger(zLogLevel, false, time.Millisecond*200)})
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

	err = db.AutoMigrate(
		&model_struct.LocalFriend{},
		&model_struct.LocalFriendRequest{},
		&model_struct.LocalGroup{},
		&model_struct.LocalGroupMember{},
		&model_struct.LocalGroupRequest{},
		&model_struct.LocalErrChatLog{},
		&model_struct.LocalUser{},
		&model_struct.LocalBlack{},
		&model_struct.LocalConversation{},
		&model_struct.NotificationSeqs{},
		&model_struct.LocalChatLog{},
		&model_struct.LocalAdminGroupRequest{},
		&model_struct.LocalWorkMomentsNotification{},
		&model_struct.LocalWorkMomentsNotificationUnreadCount{},
		&model_struct.TempCacheLocalChatLog{},
		&model_struct.LocalChatLogReactionExtensions{},
		&model_struct.LocalUpload{},
		&model_struct.LocalStranger{},
		&model_struct.LocalSendingMessages{},
		&model_struct.LocalUserCommand{},
	)
	if err != nil {
		return err
	}
	//if err := db.Table(constant.SuperGroupTableName).AutoMigrate(superGroup); err != nil {
	//	return err
	//}

	return nil
}

func (d *DataBase) versionDataFix(ctx context.Context) {
	//todo some model auto migrate data conversion
	//conversationIDs, err := d.FindAllConversationConversationID(ctx)
	//if err != nil {
	//	log.ZError(ctx, "FindAllConversationConversationID err", err)
	//}
	//for _, conversationID := range conversationIDs {
	//	d.conn.WithContext(ctx).Table(utils.GetTableName(conversationID)).AutoMigrate(&model_struct.LocalChatLog{})
	//}
}
