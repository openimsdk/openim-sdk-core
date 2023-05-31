package conversation_msg

import (
	"context"
	"errors"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	pbMsg "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	utils2 "github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

func (c *Conversation) revokeMessage(ctx context.Context, msg *sdkws.MsgData) {
	var tips sdkws.RevokeMsgTips
	if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
		log.ZError(ctx, "unmarshal failed", err, "msg", msg)
		return
	}
	log.ZDebug(ctx, "revokeMessage", "tips", &tips)
	revokedMsg, err := c.db.GetMessageBySeq(ctx, tips.ConversationID, tips.Seq)
	if err != nil {
		log.ZError(ctx, "GetMessageBySeq failed", err, "tips", &tips)
		return
	}
	var revokerRole int32
	var revokerNickname string
	if tips.SesstionType == constant.SuperGroupChatType {
		groupMember, err := c.db.GetGroupMemberInfoByGroupIDUserID(ctx, msg.GroupID, tips.RevokerUserID)
		if err != nil {
			log.ZError(ctx, "GetGroupMemberInfoByGroupIDUserID failed", err, "tips", &tips)
		}
		revokerRole = groupMember.RoleLevel
		revokerNickname = groupMember.Nickname
	} else {
		_, userName, err := c.cache.GetUserNameAndFaceURL(ctx, tips.RevokerUserID)
		if err != nil {
			log.ZError(ctx, "GetUserNameAndFaceURL failed", err, "tips", &tips)
		}
		revokerNickname = userName
	}
	m := sdk_struct.MessageRevoked{
		RevokerID:                   tips.RevokerUserID,
		RevokerRole:                 revokerRole,
		ClientMsgID:                 revokedMsg.ClientMsgID,
		RevokerNickname:             revokerNickname,
		RevokeTime:                  tips.RevokeTime,
		SourceMessageSendTime:       revokedMsg.SendTime,
		SourceMessageSendID:         revokedMsg.SendID,
		SourceMessageSenderNickname: revokedMsg.SenderNickname,
		SessionType:                 int32(tips.SesstionType),
		Seq:                         tips.Seq,
		Ex:                          revokedMsg.Ex,
	}
	var n sdk_struct.NotificationElem
	n.Detail = utils.StructToJsonString(m)
	if err := c.db.UpdateMessageBySeq(ctx, tips.ConversationID, &model_struct.LocalChatLog{Seq: tips.Seq,
		Status: constant.MsgStatusRevoked, Content: utils.StructToJsonString(n), ContentType: constant.RevokeNotification}); err != nil {
		log.ZError(ctx, "UpdateMessageBySeq failed", err, "tips", &tips)
		return
	}
	c.msgListener.OnNewRecvMessageRevoked(utils.StructToJsonString(m))
	msgList, err := c.db.SearchAllMessageByContentType(ctx, constant.Quote)
	if err != nil {
		log.ZError(ctx, "SearchAllMessageByContentType failed", err, "tips", &tips)
		return
	}
	for _, v := range msgList {
		c.quoteMsgRevokeHandle(ctx, tips.ConversationID, v, m)
	}
}

func (c *Conversation) quoteMsgRevokeHandle(ctx context.Context, conversationID string, v *model_struct.LocalChatLog, revokedMsg sdk_struct.MessageRevoked) {
	s := sdk_struct.MsgStruct{}
	_ = utils.JsonStringToStruct(v.Content, &s.QuoteElem)

	if s.QuoteElem.QuoteMessage == nil {
		return
	}
	if s.QuoteElem.QuoteMessage.ClientMsgID != revokedMsg.ClientMsgID {
		return
	}
	s.QuoteElem.QuoteMessage.Content = utils.StructToJsonString(revokedMsg)
	s.QuoteElem.QuoteMessage.ContentType = constant.RevokeNotification
	v.Content = utils.StructToJsonString(s.QuoteElem)
	if err := c.db.UpdateMessageBySeq(ctx, conversationID, v); err != nil {
		log.ZError(ctx, "UpdateMessage failed", err, "v", v)
	}
}

func (c *Conversation) revokeOneMessage(ctx context.Context, req *sdk_struct.MsgStruct) error {
	var conversationID string
	switch req.SessionType {
	case constant.SingleChatType:
		conversationID = utils2.GetConversationIDBySessionType(int(req.SessionType), req.SendID, req.RecvID)
	case constant.SuperGroupChatType:
		conversationID = utils2.GetConversationIDBySessionType(int(req.SessionType), req.GroupID)
	}
	message, err := c.db.GetMessage(ctx, conversationID, req.ClientMsgID)
	if err != nil {
		return err
	}
	if message.Status != constant.MsgStatusSendSuccess {
		return errors.New("only send success message can be revoked")
	}
	switch req.SessionType {
	case constant.SingleChatType:
		if message.SendID != c.loginUserID {
			return errors.New("only send by yourself message can be revoked")
		}
	case constant.SuperGroupChatType:
		if message.SendID != c.loginUserID {
			groupAdmins, err := c.db.GetGroupMemberOwnerAndAdmin(ctx, req.GroupID)
			if err != nil {
				return err
			}
			var isAdmin bool
			for _, member := range groupAdmins {
				if member.UserID == c.loginUserID {
					isAdmin = true
					break
				}
			}
			if !isAdmin {
				return errors.New("only group admin can revoke message")
			}
		}
	}
	if err := util.ApiPost(ctx, constant.RevokeMsgRouter, pbMsg.RevokeMsgReq{ConversationID: conversationID, Seq: message.Seq, UserID: c.loginUserID}, &pbMsg.RevokeMsgResp{}); err != nil {
		return err
	}
	c.revokeMessage(ctx, nil)
	return nil
}
