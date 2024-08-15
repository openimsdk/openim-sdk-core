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

	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/protocol/sdkws"
)

type MessageController struct {
	db db_interface.DataBase
	ch chan common.Cmd2Value
}

func NewMessageController(db db_interface.DataBase, ch chan common.Cmd2Value) *MessageController {
	return &MessageController{db: db, ch: ch}
}

func (c *Conversation) PullMessageBySeqs(ctx context.Context, seqs []*sdkws.SeqRange) (*sdkws.PullMessageBySeqsResp, error) {
	return util.CallApi[sdkws.PullMessageBySeqsResp](ctx, constant.PullUserMsgBySeqRouter, sdkws.PullMessageBySeqsReq{UserID: c.loginUserID, SeqRanges: seqs})
}
func (m *MessageController) SearchMessageByContentTypeAndKeyword(ctx context.Context, contentType []int, keywordList []string,
	keywordListMatchType int, startTime, endTime int64) (result []*model_struct.LocalChatLog, err error) {
	var list []*model_struct.LocalChatLog
	conversationIDList, err := m.db.GetAllConversationIDList(ctx)
	for _, v := range conversationIDList {
		sList, err := m.db.SearchMessageByContentTypeAndKeyword(ctx, contentType, v, keywordList, keywordListMatchType, startTime, endTime)
		if err != nil {
			// TODO: log.Error(operationID, "search message in group err", err.Error(), v)
			continue
		}
		list = append(list, sList...)
	}

	return list, nil
}
