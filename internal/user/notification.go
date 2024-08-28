package user

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

// DoNotification handles incoming notifications for the user.
func (u *User) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	log.ZDebug(ctx, "user notification", "msg", msg)
	go func() {
		switch msg.ContentType {
		case constant.UserInfoUpdatedNotification:
			u.userInfoUpdatedNotification(ctx, msg)
		case constant.UserStatusChangeNotification:
			//u.userStatusChangeNotification(ctx, msg)
		case constant.UserCommandAddNotification:
			u.userCommandAddNotification(ctx, msg)
		case constant.UserCommandDeleteNotification:
			u.userCommandDeleteNotification(ctx, msg)
		case constant.UserCommandUpdateNotification:
			u.userCommandUpdateNotification(ctx, msg)
		default:
			// log.Error(operationID, "type failed ", msg.ClientMsgID, msg.ServerMsgID, msg.ContentType)
		}
	}()
}

// userInfoUpdatedNotification handles notifications about updated user information.
func (u *User) userInfoUpdatedNotification(ctx context.Context, msg *sdkws.MsgData) {
	log.ZDebug(ctx, "userInfoUpdatedNotification", "msg", msg)
	tips := sdkws.UserInfoUpdatedTips{}
	if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
		log.ZError(ctx, "comm.UnmarshalTips failed", err, "msg", msg.Content)
		return
	}

	if tips.UserID == u.loginUserID {
		err := u.SyncLoginUserInfo(ctx)
		if err != nil {
			log.ZWarn(ctx, "SyncLoginUserInfo", err)
			return
		}
	} else {
		log.ZDebug(ctx, "detail.UserID != u.loginUserID, do nothing", "detail.UserID", tips.UserID, "u.loginUserID", u.loginUserID)
	}
}

// userCommandAddNotification handle notification when user add favorite
func (u *User) userCommandAddNotification(ctx context.Context, msg *sdkws.MsgData) {
	log.ZDebug(ctx, "userCommandAddNotification", "msg", msg)
	tip := sdkws.UserCommandAddTips{}
	if tip.ToUserID == u.loginUserID {
		err := u.SyncAllCommand(ctx)
		if err != nil {
			log.ZWarn(ctx, "userCommandAddNotification", err)
			return
		}
	} else {
		log.ZDebug(ctx, "ToUserID != u.loginUserID, do nothing", "detail.UserID", tip.ToUserID, "u.loginUserID", u.loginUserID)
	}
}

// userCommandDeleteNotification handle notification when user delete favorite
func (u *User) userCommandDeleteNotification(ctx context.Context, msg *sdkws.MsgData) {
	log.ZDebug(ctx, "userCommandAddNotification", "msg", msg)
	tip := sdkws.UserCommandDeleteTips{}
	if tip.ToUserID == u.loginUserID {
		err := u.SyncAllCommand(ctx)
		if err != nil {
			log.ZWarn(ctx, "SyncAllCommand", err)
			return
		}
	} else {
		log.ZDebug(ctx, "ToUserID != u.loginUserID, do nothing", "detail.UserID", tip.ToUserID, "u.loginUserID", u.loginUserID)
	}
}

// userCommandUpdateNotification handle notification when user update favorite
func (u *User) userCommandUpdateNotification(ctx context.Context, msg *sdkws.MsgData) {
	log.ZDebug(ctx, "userCommandAddNotification", "msg", msg)
	tip := sdkws.UserCommandUpdateTips{}
	if tip.ToUserID == u.loginUserID {
		err := u.SyncAllCommand(ctx)
		if err != nil {
			log.ZWarn(ctx, "SyncAllCommand", err)
			return
		}
	} else {
		log.ZDebug(ctx, "ToUserID != u.loginUserID, do nothing", "detail.UserID", tip.ToUserID, "u.loginUserID", u.loginUserID)
	}
}
