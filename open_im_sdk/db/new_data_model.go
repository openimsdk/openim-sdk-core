package db

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/open_im_sdk/utils"
)

func (u *open_im_sdk.UserRelated) initDB() error {
	if u.loginUserID == "" {
		return errors.New("no uid")
	}
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()

	db, err := gorm.Open(sqlite.Open(open_im_sdk.SvrConf.DbDir+"OpenIM_"+u.loginUserID+".db"), &gorm.Config{})
	utils.sdkLog("open db:", open_im_sdk.SvrConf.DbDir+"OpenIM_"+u.loginUserID+".db")
	if err != nil {
		panic("failed to connect database" + err.Error())
		return err
	}
	u.validate = validator.New()
	u.imdb = db
	//db, err := sql.Open("sqlite3", SvrConf.DbDir+"OpenIM_"+uid+".db")
	//sdkLog("open db:", SvrConf.DbDir+"OpenIM_"+uid+".db")
	//if err != nil {
	//	sdkLog("failed open db:", SvrConf.DbDir+"OpenIM_"+uid+".db", err.Error())
	//	return err
	//}
	//u.db = db

	db.AutoMigrate(&open_im_sdk.LocalFriend{},
		&open_im_sdk.LocalFriendRequest{},
		&open_im_sdk.LocalGroup{},
		&open_im_sdk.LocalGroupMember{},
		&open_im_sdk.LocalGroupRequest{},
		&open_im_sdk.LocalUser{},
		&open_im_sdk.LocalBlack{}, &open_im_sdk.LocalSeqData{})
	if !db.Migrator().HasTable(&open_im_sdk.LocalFriend{}) {
		//log.NewInfo("CreateTable Friend")
		db.Migrator().CreateTable(&open_im_sdk.LocalFriend{})
	}

	if !db.Migrator().HasTable(&open_im_sdk.LocalFriendRequest{}) {
		//log.NewInfo("CreateTable FriendRequest")
		db.Migrator().CreateTable(&open_im_sdk.LocalFriendRequest{})
	}

	if !db.Migrator().HasTable(&open_im_sdk.LocalGroup{}) {
		//log.NewInfo("CreateTable Group")
		db.Migrator().CreateTable(&open_im_sdk.LocalGroup{})
	}

	if !db.Migrator().HasTable(&open_im_sdk.LocalGroupMember{}) {
		//log.NewInfo("CreateTable GroupMember")
		db.Migrator().CreateTable(&open_im_sdk.LocalGroupMember{})
	}

	if !db.Migrator().HasTable(&open_im_sdk.LocalGroupRequest{}) {
		//log.NewInfo("CreateTable GroupRequest")
		db.Migrator().CreateTable(&open_im_sdk.LocalGroupRequest{})
	}

	if !db.Migrator().HasTable(&open_im_sdk.LocalUser{}) {
		//log.NewInfo("CreateTable User")
		db.Migrator().CreateTable(&open_im_sdk.LocalUser{})
	}

	if !db.Migrator().HasTable(&open_im_sdk.LocalBlack{}) {
		//log.NewInfo("CreateTable Black")
		db.Migrator().CreateTable(&open_im_sdk.LocalBlack{})
	}

	if !db.Migrator().HasTable(&open_im_sdk.LocalSeqData{}) {
		db.Migrator().CreateTable(&open_im_sdk.LocalSeqData{})
	}
	return nil
}
