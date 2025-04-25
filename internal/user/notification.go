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
	if err := u.doNotification(ctx, msg); err != nil {
		log.ZError(ctx, "DoUserNotification failed", err)
	}
}

func (u *User) doNotification(ctx context.Context, msg *sdkws.MsgData) error {
	switch msg.ContentType {
	case constant.UserInfoUpdatedNotification:
		return u.userInfoUpdatedNotification(ctx, msg)
	case constant.UserCommandAddNotification:
		return u.userCommandAddNotification(ctx, msg)
	case constant.UserCommandDeleteNotification:
		return u.userCommandDeleteNotification(ctx, msg)
	case constant.UserCommandUpdateNotification:
		return u.userCommandUpdateNotification(ctx, msg)
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

// userCommandAddNotification handle notification when user add favorite
func (u *User) userCommandAddNotification(ctx context.Context, msg *sdkws.MsgData) error {
	tip := sdkws.UserCommandAddTips{}
	if tip.ToUserID == u.loginUserID {
		err := u.SyncAllCommand(ctx)
		if err != nil {
			return err
		}
	} else {
		log.ZDebug(ctx, "ToUserID != u.loginUserID, do nothing", "detail.UserID", tip.ToUserID, "u.loginUserID", u.loginUserID)
	}
	return nil
}

// userCommandDeleteNotification handle notification when user delete favorite
func (u *User) userCommandDeleteNotification(ctx context.Context, msg *sdkws.MsgData) error {
	tip := sdkws.UserCommandDeleteTips{}
	if tip.ToUserID == u.loginUserID {
		err := u.SyncAllCommand(ctx)
		if err != nil {
			return err
		}
	} else {
		log.ZDebug(ctx, "ToUserID != u.loginUserID, do nothing", "detail.UserID", tip.ToUserID, "u.loginUserID", u.loginUserID)
	}
	return nil
}

// userCommandUpdateNotification handle notification when user update favorite
func (u *User) userCommandUpdateNotification(ctx context.Context, msg *sdkws.MsgData) error {
	tip := sdkws.UserCommandUpdateTips{}
	if tip.ToUserID == u.loginUserID {
		err := u.SyncAllCommand(ctx)
		if err != nil {
			return err
		}
	} else {
		log.ZDebug(ctx, "ToUserID != u.loginUserID, do nothing", "detail.UserID", tip.ToUserID, "u.loginUserID", u.loginUserID)
	}
	return nil
}
