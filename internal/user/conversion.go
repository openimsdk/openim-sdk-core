package user

import (
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/sdkws"
)

func ServerUserToLocalUser(user *sdkws.UserInfo) *model_struct.LocalUser {
	return &model_struct.LocalUser{
		UserID:     user.UserID,
		Nickname:   user.Nickname,
		FaceURL:    user.FaceURL,
		CreateTime: user.CreateTime,
		Ex:         user.Ex,
		//AppMangerLevel:   user.AppMangerLevel,
		GlobalRecvMsgOpt: user.GlobalRecvMsgOpt,
		//AttachedInfo: user.AttachedInfo,
	}
}
func LocalUserToPublicUser(user *model_struct.LocalUser) *sdk_struct.PublicUser {
	return &sdk_struct.PublicUser{
		UserID:     user.UserID,
		Nickname:   user.Nickname,
		FaceURL:    user.FaceURL,
		Ex:         user.Ex,
		CreateTime: user.CreateTime,
	}
}
