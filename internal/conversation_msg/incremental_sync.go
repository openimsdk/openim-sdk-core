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
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/syncer"
	pbConversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
)

// fullSyncReconcileBatchSize 非空库全量对账时，每批从服务端拉取并比对的会话数。
// 100/批的内存占用约 1-2MB，万级会话也不会触发 iOS jetsam。
const fullSyncReconcileBatchSize = 100

func (c *Conversation) IncrSyncConversations(ctx context.Context) error {
	conversationSyncer := syncer.VersionSynchronizer[*model_struct.LocalConversation, *pbConversation.GetIncrementalConversationResp]{
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
			return c.getIncrementalConversationFromServer(ctx, version.Version, version.VersionID)
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
			return c.conversationSyncer.Sync(ctx, server, local, nil, true)
		},
		FullSyncer: func(ctx context.Context) error {
			conversationIDList, err := c.db.GetAllConversationIDList(ctx)
			if err != nil {
				return err
			}
			if len(conversationIDList) == 0 {
				return c.conversationSyncer.FullSync(ctx, c.loginUserID)
			}
			// 非空库对账（中断恢复场景）。原实现一次性拉全部会话
			// （GetAllConversations 协议无分页），万级会话账号会把 Go 堆撑到
			// GB 级，iOS 触发 jetsam 杀进程；且同步版本号只在整轮成功后落库，
			// 中断重启后必然再次落入本分支 → 死循环。
			// 改为：拉服务端全部会话 ID（纯 ID 列表，轻量）→ 按
			// fullSyncReconcileBatchSize 分批「服务端拉取 + 本地比对」，
			// 每批处理完即释放，内存峰值恒定；本地已有数据只做增量比对，
			// 中断重启后无需重拉全量 → 天然断点续传。
			// 删除语义与原实现一致（skipDeletion=true，删除由增量通道处理）。
			resp, err := c.getAllConversationIDsFromServer(ctx)
			if err != nil {
				return err
			}
			serverIDs := resp.ConversationIDs
			localIDSet := make(map[string]struct{}, len(conversationIDList))
			for _, id := range conversationIDList {
				localIDSet[id] = struct{}{}
			}
			var missingNum int
			for _, id := range serverIDs {
				if _, ok := localIDSet[id]; !ok {
					missingNum++
				}
			}
			log.ZDebug(ctx, "full sync reconcile by ids", "serverNum", len(serverIDs),
				"localNum", len(conversationIDList), "missingNum", missingNum)

			reconcile := func(ids []string) error {
				convResp, err := c.getConversationsByIDsFromServer(ctx, ids)
				if err != nil {
					return err
				}
				server := make([]*model_struct.LocalConversation, 0, len(convResp.Conversations))
				for _, conv := range convResp.Conversations {
					server = append(server, ServerConversationToLocal(conv))
				}
				local := make([]*model_struct.LocalConversation, 0, len(ids))
				for _, id := range ids {
					lc, err := c.db.GetConversation(ctx, id)
					if err != nil {
						// 本地缺失：跳过，交给 Sync 按缺失走 insert
						if errs.Unwrap(err) == errs.ErrRecordNotFound {
							continue
						}
						return err
					}
					local = append(local, lc)
				}
				return c.conversationSyncer.Sync(ctx, server, local, nil, true)
			}

			for start := 0; start < len(serverIDs); start += fullSyncReconcileBatchSize {
				end := start + fullSyncReconcileBatchSize
				if end > len(serverIDs) {
					end = len(serverIDs)
				}
				batchStart := time.Now()
				if err := reconcile(serverIDs[start:end]); err != nil {
					return err
				}
				log.ZDebug(ctx, "full sync reconcile batch done",
					"index", start/fullSyncReconcileBatchSize,
					"size", end-start,
					"duration", time.Since(batchStart).String())
			}
			return nil
		},
		FullID: func(ctx context.Context) ([]string, error) {
			resp, err := c.getAllConversationIDsFromServer(ctx)
			if err != nil {
				return nil, err
			}
			return resp.ConversationIDs, nil
		},
	}

	return conversationSyncer.IncrementalSync()
}

func (c *Conversation) conversationTableName() string {
	return model_struct.LocalConversation{}.TableName()
}

func (c *Conversation) IncrSyncConversationsWithLock(ctx context.Context) error {
	c.conversationSyncMutex.Lock()
	defer c.conversationSyncMutex.Unlock()
	return c.IncrSyncConversations(ctx)
}
