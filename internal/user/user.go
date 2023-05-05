// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package user

import (
	"context"
	"fmt"
	comm "open_im_sdk/internal/common"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/sdkerrs"
	"open_im_sdk/pkg/syncer"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	authPb "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/auth"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	userPb "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"

	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"
)

// User is a struct that represents a user in the system.
type User struct {
	db_interface.DataBase
	loginUserID    string
	listener       open_im_sdk_callback.OnUserListener
	loginTime      int64
	userSyncer     *syncer.Syncer[*model_struct.LocalUser, string]
	conversationCh chan common.Cmd2Value
}

// LoginTime gets the login time of the user.
func (u *User) LoginTime() int64 {
	return u.loginTime
}

// SetLoginTime sets the login time of the user.
func (u *User) SetLoginTime(loginTime int64) {
	u.loginTime = loginTime
}

// SetListener sets the user's listener.
func (u *User) SetListener(listener open_im_sdk_callback.OnUserListener) {
	u.listener = listener
}

// NewUser creates a new User object.
func NewUser(dataBase db_interface.DataBase, loginUserID string, conversationCh chan common.Cmd2Value) *User {
	user := &User{DataBase: dataBase, loginUserID: loginUserID, conversationCh: conversationCh}
	user.initSyncer()
	return user
}

func (u *User) initSyncer() {
	u.userSyncer = syncer.New(
		func(ctx context.Context, value *model_struct.LocalUser) error {
			return u.InsertLoginUser(ctx, value)
		},
		func(ctx context.Context, value *model_struct.LocalUser) error {
			return fmt.Errorf("not support delete user %s", value.UserID)
		},
		func(ctx context.Context, serverUser, localUser *model_struct.LocalUser) error {
			return u.DataBase.UpdateLoginUser(context.Background(), serverUser)
		},
		func(user *model_struct.LocalUser) string {
			return user.UserID
		},
		nil,
		nil,
	)
}

// DoNotification handles incoming notifications for the user.
func (u *User) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg)
	if u.listener == nil {
		log.Error(operationID, "listener == nil")
		return
	}
	if msg.SendTime < u.loginTime {
		log.ZWarn(ctx, "ignore notification ", nil, "msg", *msg)
		return
	}
	go func() {
		switch msg.ContentType {
		case constant.UserInfoUpdatedNotification:
			u.userInfoUpdatedNotification(msg, operationID)
		default:
			log.Error(operationID, "type failed ", msg.ClientMsgID, msg.ServerMsgID, msg.ContentType)
		}
	}()
}

// userInfoUpdatedNotification handles notifications about updated user information.
func (u *User) userInfoUpdatedNotification(msg *sdkws.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	var detail sdkws.UserInfoUpdatedTips
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "comm.UnmarshalTips failed ", err.Error(), msg.Content)
		return
	}
	if detail.UserID == u.loginUserID {
		log.Info(operationID, "detail.UserID == u.loginUserID, SyncLoginUserInfo", detail.UserID)
		u.SyncLoginUserInfo(context.Background())
	} else {
		log.Debug(operationID, "detail.UserID != u.loginUserID, do nothing", detail.UserID, u.loginUserID)
	}
}

// GetUsersInfoFromSvr retrieves user information from the server.
func (u *User) GetUsersInfoFromSvr(ctx context.Context, userIDs []string) ([]*model_struct.LocalUser, error) {
	resp, err := util.CallApi[userPb.GetDesignateUsersResp](ctx, constant.GetUsersInfoRouter, userPb.GetDesignateUsersReq{UserIDs: userIDs})
	if err != nil {
		return nil, sdkerrs.Warp(err, "GetUsersInfoFromSvr failed")
	}
	return util.Batch(ServerUserToLocalUser, resp.UsersInfo), nil
}

// GetUsersInfoFromSvrNoCallback retrieves user information from the server.
func (u *User) GetSingleUserFromSvr(ctx context.Context, userID string) (*model_struct.LocalUser, error) {
	users, err := u.GetUsersInfoFromSvr(ctx, []string{userID})
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return users[0], nil
	}
	return nil, sdkerrs.ErrRecordNotFound.Wrap(fmt.Sprintf("getSelfUserInfo failed, userID: %s not exist", userID))
}

// getSelfUserInfo retrieves the user's information.
func (u *User) getSelfUserInfo(ctx context.Context) (*model_struct.LocalUser, error) {
	return u.GetLoginUser(ctx, u.loginUserID)
}

// updateSelfUserInfo updates the user's information.
func (u *User) updateSelfUserInfo(ctx context.Context, userInfo *sdkws.UserInfo) error {
	userInfo.UserID = u.loginUserID
	if err := util.ApiPost(ctx, constant.UpdateSelfUserInfoRouter, userPb.UpdateUserInfoReq{UserInfo: userInfo}, nil); err != nil {
		return err
	}
	_ = u.SyncLoginUserInfo(ctx)
	return nil
}

// ParseTokenFromSvr parses a token from the server.
func (u *User) ParseTokenFromSvr(ctx context.Context) (int64, error) {
	resp, err := util.CallApi[authPb.ParseTokenResp](ctx, constant.ParseTokenRouter, authPb.ParseTokenReq{})
	return resp.ExpireTimeSeconds, err
}

// GetServerUserInfo retrieves user information from the server.
func (u *User) GetServerUserInfo(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error) {
	resp, err := util.CallApi[userPb.GetDesignateUsersResp](ctx, constant.GetUsersInfoRouter, &userPb.GetDesignateUsersReq{UserIDs: userIDs})
	if err != nil {
		return nil, err
	}
	return resp.UsersInfo, nil
}
