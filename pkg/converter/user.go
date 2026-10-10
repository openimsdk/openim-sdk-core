package converter

import (
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/sdkws"
)

func ServerUserToLocal(info *sdkws.UserInfo) *model_struct.LocalUser {
	if info == nil {
		return nil
	}
	return &model_struct.LocalUser{
		UserID:           info.UserID,
		Nickname:         info.Nickname,
		FaceURL:          info.FaceURL,
		CreateTime:       info.CreateTime,
		Ex:               info.Ex,
		GlobalRecvMsgOpt: info.GlobalRecvMsgOpt,
	}
}

func LocalUserToPublic(user *model_struct.LocalUser) *sdk_struct.PublicUser {
	if user == nil {
		return nil
	}
	return &sdk_struct.PublicUser{
		UserID:     user.UserID,
		Nickname:   user.Nickname,
		FaceURL:    user.FaceURL,
		Ex:         user.Ex,
		CreateTime: user.CreateTime,
	}
}
