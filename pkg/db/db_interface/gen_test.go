package db_interface

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"open_im_sdk/pkg/db/model_struct"
	"os"
	"strings"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	data, err := os.ReadFile("databse.go")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n")
	var arr []string
	for _, line := range lines {
		if strings.Index(line, "(") > 0 && strings.Index(line, ")") > 0 {
			line = strings.Replace(line, "(", "(ctx context.Context, ", 1)
		}
		arr = append(arr, line)
	}
	fmt.Println(strings.Join(arr, "\n"))
}

func TestSync(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&model_struct.LocalUser{}); err != nil {
		panic(err)
	}

	changes := []model_struct.LocalUser{
		{
			UserID:           "1234567",
			Nickname:         "123",
			FaceURL:          "456",
			Gender:           1,
			PhoneNumber:      "111111111111",
			Birth:            12345,
			Email:            "123@qq.com",
			CreateTime:       uint32(time.Now().Unix()),
			AppMangerLevel:   1,
			Ex:               "",
			AttachedInfo:     "",
			GlobalRecvMsgOpt: 1,
			BirthTime:        time.Now(),
		},
		{
			UserID:           "",
			Nickname:         "123",
			FaceURL:          "456",
			Gender:           1,
			PhoneNumber:      "123124554",
			Birth:            12345,
			Email:            "123@qq.com",
			CreateTime:       uint32(time.Now().Unix()),
			AppMangerLevel:   1,
			Ex:               "",
			AttachedInfo:     "",
			GlobalRecvMsgOpt: 1,
			BirthTime:        time.Now(),
		},
	}

	db = db.Debug()

	changeStates, deleteStates, err := SyncDB(db, changes, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(changeStates)
	fmt.Println(deleteStates)
}
