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

package conversation_msg

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/internal/incrversion"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	pbConversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/tools/utils/datautil"
)

func (c *Conversation) IncrSyncConversations(ctx context.Context) error {
	conversationSyncer := incrversion.VersionSynchronizer[*model_struct.LocalConversation, *pbConversation.GetIncrementalConversationResp]{
		Ctx:       ctx,
		DB:        c.db,
		TableName: c.conversationTableName(),
		EntityID:  c.loginUserID,
		Key: func(localConversation *model_struct.LocalConversation) string {
			return localConversation.ConversationID
		},
		Local: func() ([]*model_struct.LocalConversation, error) {
			return c.db.GetAllConversations(ctx)
		},
		Server: func(version *model_struct.LocalVersionSync) (*pbConversation.GetIncrementalConversationResp, error) {
			return util.CallApi[pbConversation.GetIncrementalConversationResp](ctx, constant.GetIncrementalFriends, &pbConversation.GetIncrementalConversationReq{
				UserID:    c.loginUserID,
				Version:   version.Version,
				VersionID: version.VersionID,
			})
		},
		Full: func(resp *pbConversation.GetIncrementalConversationResp) bool {
			return resp.Full
		},
		Version: func(resp *pbConversation.GetIncrementalConversationResp) (string, uint64) {
			return resp.VersionID, resp.Version
		},
		Delete: func(resp *pbConversation.GetIncrementalConversationResp) []string {
			return resp.Delete
		},
		Update: func(resp *pbConversation.GetIncrementalConversationResp) []*model_struct.LocalConversation {
			return datautil.Batch(ServerConversationToLocal, resp.Update)
		},
		Insert: func(resp *pbConversation.GetIncrementalConversationResp) []*model_struct.LocalConversation {
			return datautil.Batch(ServerConversationToLocal, resp.Insert)
		},
		Syncer: func(server, local []*model_struct.LocalConversation) error {
			return c.conversationSyncer.Sync(ctx, server, local, nil)
		},
		FullSyncer: func(ctx context.Context) error {
			return c.conversationSyncer.FullSync(ctx, c.loginUserID)
		},
		FullID: func(ctx context.Context) ([]string, error) {
			resp, err := util.CallApi[pbConversation.GetFullOwnerConversationIDsResp](ctx, constant.GetFullConversationIDs, &pbConversation.GetFullOwnerConversationIDsReq{
				UserID: c.loginUserID,
			})
			if err != nil {
				return nil, err
			}
			return resp.ConversationIDs, nil
		},
	}

	return conversationSyncer.Sync()
}

func (c *Conversation) conversationTableName() string {
	return model_struct.LocalConversation{}.TableName()
}
