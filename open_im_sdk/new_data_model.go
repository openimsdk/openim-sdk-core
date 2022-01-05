package open_im_sdk

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (u *UserRelated) initDB() error {
	if u.loginUserID == "" {
		return errors.New("no uid")
	}
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()

	db, err := gorm.Open(sqlite.Open(SvrConf.DbDir+"OpenIM_"+u.loginUserID+".db"), &gorm.Config{})
	sdkLog("open db:", SvrConf.DbDir+"OpenIM_"+u.loginUserID+".db")
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

	db.AutoMigrate(&LocalFriend{},
		&LocalFriendRequest{},
		&LocalGroup{},
		&LocalGroupMember{},
		&LocalGroupRequest{},
		&LocalUser{},
		&LocalBlack{}, &LocalSeqData{})
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
	return nil
}
