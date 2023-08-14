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

package db

import (
	"context"
	"errors"
	"open_im_sdk/wasm/exec"
	"open_im_sdk/wasm/indexdb"
)

var ErrType = errors.New("from javascript data type err")

type IndexDB struct {
	*indexdb.LocalUsers
	*indexdb.LocalConversations
	*indexdb.LocalChatLogs
	*indexdb.LocalSuperGroupChatLogs
	*indexdb.LocalSuperGroup
	*indexdb.LocalConversationUnreadMessages
	*indexdb.LocalGroups
	*indexdb.LocalGroupMember
	*indexdb.LocalCacheMessage
	*indexdb.FriendRequest
	*indexdb.Black
	*indexdb.Friend
	*indexdb.LocalGroupRequest
	*indexdb.LocalChatLogReactionExtensions
	*indexdb.NotificationSeqs
	*indexdb.LocalUpload
	loginUserID string
}

func (i IndexDB) Close(ctx context.Context) error {
	_, err := exec.Exec()
	return err
}

func (i IndexDB) InitDB(ctx context.Context, userID string, dataDir string) error {
	_, err := exec.Exec(userID, dataDir)
	return err
}

func NewDataBase(ctx context.Context, loginUserID string, dbDir string, logLevel int) (*IndexDB, error) {
	i := &IndexDB{
		LocalUsers:                      indexdb.NewLocalUsers(),
		LocalConversations:              indexdb.NewLocalConversations(),
		LocalChatLogs:                   indexdb.NewLocalChatLogs(loginUserID),
		LocalSuperGroupChatLogs:         indexdb.NewLocalSuperGroupChatLogs(),
		LocalSuperGroup:                 indexdb.NewLocalSuperGroup(),
		LocalConversationUnreadMessages: indexdb.NewLocalConversationUnreadMessages(),
		LocalGroups:                     indexdb.NewLocalGroups(),
		LocalGroupMember:                indexdb.NewLocalGroupMember(),
		LocalCacheMessage:               indexdb.NewLocalCacheMessage(),
		FriendRequest:                   indexdb.NewFriendRequest(loginUserID),
		Black:                           indexdb.NewBlack(loginUserID),
		Friend:                          indexdb.NewFriend(loginUserID),
		LocalGroupRequest:               indexdb.NewLocalGroupRequest(),
		LocalChatLogReactionExtensions:  indexdb.NewLocalChatLogReactionExtensions(),
		NotificationSeqs:                indexdb.NewNotificationSeqs(),
		LocalUpload:                     indexdb.NewLocalUpload(),
		loginUserID:                     loginUserID,
	}
	err := i.InitDB(ctx, loginUserID, dbDir)
	if err != nil {
		return nil, err
	}
	return i, nil
}
