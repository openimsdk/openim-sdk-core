package user

import (
	"context"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

func (u *User) GetUsersInfo(ctx context.Context, userIDs []string) ([]*model_struct.LocalUser, error) {
	return u.GetUsersInfoFromSvr(ctx, userIDs)
}

func (u *User) GetSelfUserInfo(ctx context.Context) (*model_struct.LocalUser, error) {
	return u.getSelfUserInfo(ctx)
}

func (u *User) SetSelfInfo(ctx context.Context, userInfo *sdkws.UserInfo) error {
	return u.updateSelfUserInfo(ctx, userInfo)
}

func (u *User) UpdateMsgSenderInfo(ctx context.Context, nickname, faceURL string) (err error) {
	if nickname != "" {
		if err = u.DataBase.UpdateMsgSenderNickname(ctx, u.loginUserID, nickname, constant.SingleChatType); err != nil {
			return err
		}
	}
	if faceURL != "" {
		if err = u.DataBase.UpdateMsgSenderFaceURL(ctx, u.loginUserID, faceURL, constant.SingleChatType); err != nil {
			return err
		}
	}
	return nil
}
