package user

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

// DoNotification handles incoming notifications for the user.
func (u *User) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	log.ZDebug(ctx, "user notification", "msg", msg)
	go func() {
		if err := u.doNotification(ctx, msg); err != nil {
			log.ZError(ctx, "DoUserNotification failed", err)
		}
	}()
}

func (u *User) doNotification(ctx context.Context, msg *sdkws.MsgData) error {
	switch msg.ContentType {
	case constant.UserInfoUpdatedNotification:
		return u.userInfoUpdatedNotification(ctx, msg)
	default:
		return errs.New("unknown content type", "contentType", msg.ContentType, "clientMsgID", msg.ClientMsgID, "serverMsgID", msg.ServerMsgID).Wrap()
	}
}

// userInfoUpdatedNotification handles notifications about updated user information.
func (u *User) userInfoUpdatedNotification(ctx context.Context, msg *sdkws.MsgData) error {
	tips := sdkws.UserInfoUpdatedTips{}
	if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
		return err
	}

	if tips.UserID == u.loginUserID {
		err := u.SyncLoginUserInfo(ctx)
		if err != nil {
			return err
		}
	} else {
		log.ZDebug(ctx, "detail.UserID != u.loginUserID, do nothing", "detail.UserID", tips.UserID, "u.loginUserID", u.loginUserID)
	}
	return nil
}
