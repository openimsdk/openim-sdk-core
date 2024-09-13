package user

import (
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/user"

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
func ServerCommandToLocalCommand(data *user.AllCommandInfoResp) *model_struct.LocalUserCommand {
	return &model_struct.LocalUserCommand{
		Type:       data.Type,
		CreateTime: data.CreateTime,
		Uuid:       data.Uuid,
		Value:      data.Value,
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
