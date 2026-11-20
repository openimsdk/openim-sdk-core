package user

import (
	"context"
	"strings"
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
)

type profileInsertTestDB struct {
	db_interface.DataBase
	user *model_struct.LocalUser
}

func (d *profileInsertTestDB) GetLoginUser(context.Context, string) (*model_struct.LocalUser, error) {
	return d.user, nil
}

func (d *profileInsertTestDB) InsertLoginUser(_ context.Context, user *model_struct.LocalUser) error {
	d.user = user
	return nil
}

type profileInsertTestListener struct {
	selfInfo string
}

func (l *profileInsertTestListener) OnSelfInfoUpdated(userInfo string) {
	l.selfInfo = userInfo
}

func (*profileInsertTestListener) OnUserStatusChanged(string) {}

func TestProfileInsertInvalidatesCacheAndNotifies(t *testing.T) {
	ctx := context.Background()
	old := &model_struct.LocalUser{UserID: "self", Ex: "old"}
	updated := &model_struct.LocalUser{UserID: "self", Ex: "new"}
	database := &profileInsertTestDB{}
	listener := &profileInsertTestListener{}
	u := &User{
		DataBase:    database,
		loginUserID: "self",
		listener: func() open_im_sdk_callback.OnUserListener {
			return listener
		},
	}
	u.initSyncer()
	u.UserCache().Store("self", old)

	if err := u.userSyncer.Sync(ctx, []*model_struct.LocalUser{updated}, nil, nil); err != nil {
		t.Fatal(err)
	}
	readBack, err := u.GetSelfUserInfo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if readBack.Ex != "new" {
		t.Fatalf("read back stale profile: %q", readBack.Ex)
	}
	if !strings.Contains(listener.selfInfo, "new") {
		t.Fatalf("insert did not emit updated profile: %q", listener.selfInfo)
	}
}
