package conversation_msg

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

type MessageController struct {
	db db_interface.DataBase
}

func NewMessageController(db db_interface.DataBase) *MessageController {
	return &MessageController{db: db}
}
func (m *MessageController) BatchUpdateMessageList(ctx context.Context, updateMsg map[string][]*model_struct.LocalChatLog) error {
	if updateMsg == nil {
		return nil
	}
	for conversationID, messages := range updateMsg {
		for _, v := range messages {
			v1 := new(model_struct.LocalChatLog)
			v1.ClientMsgID = v.ClientMsgID
			v1.Seq = v.Seq
			v1.Status = v.Status
			v1.RecvID = v.RecvID
			v1.SessionType = v.SessionType
			v1.ServerMsgID = v.ServerMsgID
			err := m.db.UpdateMessage(ctx, conversationID, v1)
			if err != nil {
				return utils.Wrap(err, "BatchUpdateMessageList failed")
			}
		}

	}
	return nil
}

func (m *MessageController) BatchInsertMessageList(ctx context.Context, insertMsg map[string][]*model_struct.LocalChatLog) error {
	if insertMsg == nil {
		return nil
	}
	for conversationID, messages := range insertMsg {
		err := m.db.BatchInsertMessageList(ctx, conversationID, messages)
		if err != nil {
			log.ZError(ctx, "insert GetMessage detail err:", err, "conversationID", conversationID, "messages", messages)
			for _, v := range messages {
				e := m.db.InsertMessage(ctx, conversationID, v)
				if e != nil {
					log.ZError(ctx, "InsertMessage err", err, "conversationID", conversationID, "message", v)
				}
			}
		}

	}
	return nil
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
