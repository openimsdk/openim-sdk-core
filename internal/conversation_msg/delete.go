package conversation_msg

import (
	"context"
	"errors"
	pbMsg "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"open_im_sdk/internal/util"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/jinzhu/copier"
	"go.starlark.net/lib/proto"
)

/*
first:
+ delele server message
+ delete local message
*/

type DeleteInterface interface {
	// DeleteMessageFromLocal deletes a message from local storage.
	DeleteMessageFromLocal(ctx context.Context, s *sdk_struct.MsgStruct) error
	// DeleteConversationFromLocalAndSvr deletes a conversation from both local and server storage.
	DeleteConversationFromLocalAndSvr(ctx context.Context, conversationID string) error
	// DeleteMessageFromLocalAndSvr deletes a message from both local and server storage.
	DeleteMessageFromLocalAndSvr(ctx context.Context, s *sdk_struct.MsgStruct) error
	// DeleteAllMsgFromLocalAndSvr deletes all messages from both local and server storage.
	DeleteAllMsgFromLocalAndSvr(ctx context.Context) error
	// DeleteAllMsgFromLocal deletes all messages from local storage.
	DeleteAllMsgFromLocal(ctx context.Context) error
	// DeleteMessageFromLocalStorage deletes a message from local storage.
	DeleteMessageFromLocalStorage(ctx context.Context, message *sdk_struct.MsgStruct) error
	// ClearC2CHistoryMessage clears all messages in a C2C conversation.
	ClearC2CHistoryMessage(ctx context.Context, userID string) error
	// ClearGroupHistoryMessage clears all messages in a group conversation.
	ClearGroupHistoryMessage(ctx context.Context, groupID string) error
	// ClearC2CHistoryMessageFromLocalAndSvr clears all messages in a C2C conversation from both local and server storage.
	ClearC2CHistoryMessageFromLocalAndSvr(ctx context.Context, userID string) error
	// ClearGroupHistoryMessageFromLocalAndSvr clears all messages in a group conversation from both local and server storage.
	ClearGroupHistoryMessageFromLocalAndSvr(ctx context.Context, groupID string) error
}

// Delete all messages
func (c *Conversation) DeleteAllMessage(ctx context.Context) error {
	err := c.deleteAllMsgFromLocal(ctx)
	if err != nil {
		return err
	}
	err = c.clearMessageFromSvr(ctx)
	if err != nil {
		return err
	}
	return nil
}

// 删除单条消息
func (c *Conversation) DeleteMessageFromLocal(ctx context.Context, s *sdk_struct.MsgStruct) error {
	// 先获取到会话
	Conversation, err := c.db.GetAllConversations(ctx)
	if err != nil {
		return err
	}

	// 拿到数据库的表信息
	for _, v := range Conversation {
		if v.ConversationID == s.ClientMsgID {
			// 然后删除
			err := c.db.DeleteConversation(ctx, v.ConversationID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Delete the local and server
// Delete the local, do not change the server data
// To delete the server, you need to change the local message status to delete
func (c *Conversation) DeleteConversationFromLocalAndSvr(ctx context.Context, conversationID string) error {
	// Use conversationID to remove conversations and messages from the server first
	err := c.deleteConversationAndMsgFromSvr(ctx, conversationID)
	if err != nil {
		return err
	}
	return c.deleteConversation(ctx, conversationID)
}

// 删除单条消息，同时删除服务器上的消息
// 需要注意的是：如果是最新的一条消息，那么需要拆分出来，然后更新会话的最新消息
func (c *Conversation) DeleteMessageFromLocalAndSvr(ctx context.Context, s *sdk_struct.MsgStruct) error {
	err := c.deleteMessageFromSvr(ctx, s)
	if err != nil {
		return err
	}
	return c.deleteMessageFromLocalStorage(ctx, s)
}

// 删除所有消息，同时删除服务器上的消息
func (c *Conversation) DeleteAllMsgFromLocalAndSvr(ctx context.Context) error {
	return c.DeleteAllMsgFromLocalAndSvr(ctx)
}

// 删除所有本地消息
func (c *Conversation) DeleteAllMsgFromLocal(ctx context.Context) error {
	return c.deleteAllMsgFromLocal(ctx)
}

// 删除本地的单条消息
func (c *Conversation) DeleteMessageFromLocalStorage(ctx context.Context, message *sdk_struct.MsgStruct) error {
	return c.deleteMessageFromLocalStorage(ctx, message)
}

// 清除单个好友的历史消息记录
func (c *Conversation) ClearC2CHistoryMessage(ctx context.Context, userID string) error {
	return c.clearC2CHistoryMessage(ctx, userID)
}

// 清除单个群组的历史消息记录
func (c *Conversation) ClearGroupHistoryMessage(ctx context.Context, groupID string) error {
	return c.clearGroupHistoryMessage(ctx, groupID)

}

// 清除单个好友的历史消息记录，同时删除服务器上的消息
func (c *Conversation) ClearC2CHistoryMessageFromLocalAndSvr(ctx context.Context, userID string) error {
	conversationID := c.getConversationIDBySessionType(userID, constant.SingleChatType)
	err := c.deleteConversationAndMsgFromSvr(ctx, conversationID)
	if err != nil {
		return err
	}
	return c.clearC2CHistoryMessage(ctx, userID)

}

// 清除单个群组的历史消息记录，同时删除服务器上的消息
func (c *Conversation) ClearGroupHistoryMessageFromLocalAndSvr(ctx context.Context, groupID string) error {
	conversationID, _, err := c.getConversationTypeByGroupID(ctx, groupID)
	if err != nil {
		return err
	}
	err = c.deleteConversationAndMsgFromSvr(ctx, conversationID)
	if err != nil {
		return err
	}
	return c.clearGroupHistoryMessage(ctx, groupID)
}

// 删除单个会话
func (c *Conversation) DeleteConversation(ctx context.Context, conversationID string) error {
	lc, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	var sourceID string
	switch lc.ConversationType {
	case constant.SingleChatType, constant.NotificationChatType:
		sourceID = lc.UserID
	case constant.GroupChatType, constant.SuperGroupChatType:
		sourceID = lc.GroupID
	}
	if lc.ConversationType == constant.SuperGroupChatType {
		err = c.db.SuperGroupDeleteAllMessage(ctx, lc.GroupID)
		if err != nil {
			return err
		}
	} else {
		//Mark messages related to this conversation for deletion
		err = c.db.UpdateMessageStatusBySourceIDController(ctx, sourceID, constant.MsgStatusHasDeleted, lc.ConversationType)
		if err != nil {
			return err
		}
	}
	//Reset the conversation information, empty conversation
	err = c.db.ResetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{ConID: "", Action: constant.TotalUnreadMessageChanged, Args: ""}})
	return nil
}

func (c *Conversation) DeleteAllConversationFromLocal(ctx context.Context) error {
	err := c.db.ResetAllConversation(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Only remote information is deleted
func (c *Conversation) DeleteMessageFromRemote(ctx context.Context, s *sdk_struct.MsgStruct) error {
	// if s.ClientMsgID == "" {
	// 	return errors.New("clientMsgID is empty")
	// }
	// err := c.deleteMessageFromSvr(ctx, s)
	// if err != nil {
	// 	return err
	// }
	return nil

}

// To delete all messages, the server does it all, calls the interface,
// and the client receives a callback and deletes all messages locally.
func (c *Conversation) DeleteAllMessage(ctx context.Context) error {
	// err := c.clearMessageFromSvr(ctx)
	// if err != nil {
	// 	return err
	// }
	err := c.deleteAllMsgFromLocal(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Only the local information is deleted
// Local operation: Obtain all session ids, then obtain the table, and then clear the table
// Delete a message. If you delete the latest message, there is a latest message on the session.
// If you delete this message, there is no message on the session and the second message is displayed.
func (c *Conversation) DeleteMessageFromLocal(ctx context.Context, s *sdk_struct.MsgStruct) error {
	if s.ClientMsgID == "" {
		return errors.New("clientMsgID is empty")
	}
	err := c.deleteMessageFromLocalStorage(ctx, s)
	if err != nil {
		return err
	}
	return nil
}

func (c *Conversation) DeleteMessageFromLocalAndSvr(ctx context.Context, s *sdk_struct.MsgStruct) error {
	err := c.deleteMessageFromSvr(ctx, s)
	if err != nil {
		return err
	}
	return c.deleteMessageFromLocalStorage(ctx, s)
}

// Delete all messages from the server and local
func (c *Conversation) DeleteAllMsgFromLocalAndSvr(ctx context.Context) error {
	// err := c.clearMessageFromSvr(ctx)
	// if err != nil {
	// 	return err
	// }
	return c.DeleteAllMsgFromLocalAndSvr(ctx)
}

// Just delete the local, the server does not need to change
func (c *Conversation) DeleteAllMsgFromLocal(ctx context.Context) error {
	return c.deleteAllMsgFromLocal(ctx)
}

func (c *Conversation) DeleteMessageFromLocalStorage(ctx context.Context, message *sdk_struct.MsgStruct) error {
	return c.deleteMessageFromLocalStorage(ctx, message)
}

func (c *Conversation) ClearC2CHistoryMessage(ctx context.Context, userID string) error {
	return c.clearC2CHistoryMessage(ctx, userID)
}
func (c *Conversation) ClearGroupHistoryMessage(ctx context.Context, groupID string) error {
	return c.clearGroupHistoryMessage(ctx, groupID)

}
func (c *Conversation) ClearC2CHistoryMessageFromLocalAndSvr(ctx context.Context, userID string) error {
	conversationID := c.getConversationIDBySessionType(userID, constant.SingleChatType)
	err := c.deleteConversationAndMsgFromSvr(ctx, conversationID)
	if err != nil {
		return err
	}
	return c.clearC2CHistoryMessage(ctx, userID)

}

// fixme
func (c *Conversation) ClearGroupHistoryMessageFromLocalAndSvr(ctx context.Context, groupID string) error {
	conversationID, _, err := c.getConversationTypeByGroupID(ctx, groupID)
	if err != nil {
		return err
	}
	err = c.deleteConversationAndMsgFromSvr(ctx, conversationID)
	if err != nil {
		return err
	}
	return c.clearGroupHistoryMessage(ctx, groupID)
}

func (c *Conversation) clearGroupHistoryMessage(ctx context.Context, groupID string) error {
	_, sessionType, err := c.getConversationTypeByGroupID(ctx, groupID)
	if err != nil {
		return err
	}
	conversationID := c.getConversationIDBySessionType(groupID, int(sessionType))
	switch sessionType {
	case constant.SuperGroupChatType:
		err = c.db.SuperGroupDeleteAllMessage(ctx, groupID)
		if err != nil {
			return err
		}
	default:
		err = c.db.UpdateMessageStatusBySourceIDController(ctx, groupID, constant.MsgStatusHasDeleted, sessionType)
		if err != nil {
			return err
		}
	}

	err = c.db.ClearConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
	return nil

}

func (c *Conversation) clearC2CHistoryMessage(ctx context.Context, userID string) error {
	conversationID := c.getConversationIDBySessionType(userID, constant.SingleChatType)
	err := c.db.UpdateMessageStatusBySourceID(ctx, userID, constant.MsgStatusHasDeleted, constant.SingleChatType)
	if err != nil {
		return err
	}
	err = c.db.ClearConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
	return nil
}

// 删除服务器的话，先删除服务器信息，并且需要把本地的消息状态改成删除
func (c *Conversation) deleteMessageFromSvr(ctx context.Context, s *sdk_struct.MsgStruct) error {
	seq, err := c.db.GetMsgSeqByClientMsgIDController(ctx, s)
	if err != nil {
		return err
	}
	switch s.SessionType {
	case constant.SingleChatType, constant.GroupChatType:
		var apiReq pbMsg.DelMsgsReq
		// apiReq.Seqs = utils.Uint32ListConvert([]uint32{seq})
		// apiReq.UserID = c.loginUserID
		return util.ApiPost(ctx, constant.DeleteMsgRouter, &apiReq, nil)
	case constant.SuperGroupChatType:
		var apiReq pbMsg.DelSuperGroupMsgReq
		apiReq.UserID = c.loginUserID
		apiReq.GroupID = s.GroupID
		return util.ApiPost(ctx, constant.DeleteSuperGroupMsgRouter, &apiReq, nil)

	}
	return errors.New("session type error")

}

func (c *Conversation) clearMessageFromSvr(ctx context.Context) error {
	var apiReq pbMsg.ClearMsgReq
	apiReq.UserID = c.loginUserID
	err := util.ApiPost(ctx, constant.ClearMsgRouter, &apiReq, nil)
	if err != nil {
		return err
	}
	groupIDList, err := c.full.GetReadDiffusionGroupIDList(ctx)
	if err != nil {
		return err
	}
	var superGroupApiReq pbMsg.DelSuperGroupMsgReq
	superGroupApiReq.UserID = c.loginUserID
	for _, v := range groupIDList {
		superGroupApiReq.GroupID = v
		err := util.ApiPost(ctx, constant.DeleteSuperGroupMsgRouter, &superGroupApiReq, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// 判断删除的是否是最新的一条消息
func (c *Conversation) deleteMessageFromLocalStorage(ctx context.Context, s *sdk_struct.MsgStruct) error {
	var conversation model_struct.LocalConversation
	var latestMsg sdk_struct.MsgStruct
	var conversationID string
	var sourceID string
	chatLog := model_struct.LocalChatLog{ClientMsgID: s.ClientMsgID, Status: constant.MsgStatusHasDeleted, SessionType: s.SessionType}

	switch s.SessionType {
	case constant.GroupChatType:
		conversationID = c.getConversationIDBySessionType(s.GroupID, constant.GroupChatType)
		sourceID = s.GroupID
	case constant.SingleChatType:
		if s.SendID != c.loginUserID {
			conversationID = c.getConversationIDBySessionType(s.SendID, constant.SingleChatType)
			sourceID = s.SendID
		} else {
			conversationID = c.getConversationIDBySessionType(s.RecvID, constant.SingleChatType)
			sourceID = s.RecvID
		}
	case constant.SuperGroupChatType:
		conversationID = c.getConversationIDBySessionType(s.GroupID, constant.SuperGroupChatType)
		sourceID = s.GroupID
		chatLog.RecvID = s.GroupID
	}
	err := c.db.UpdateMessageController(ctx, &chatLog)
	if err != nil {
		return err
	}
	LocalConversation, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	err = utils.JsonStringToStruct(LocalConversation.LatestMsg, &latestMsg)
	if err != nil {
		return err
	}

	if s.ClientMsgID == latestMsg.ClientMsgID { //If the deleted message is the latest message of the conversation, update the latest message of the conversation
		list, err := c.db.GetMessageListNoTimeController(ctx, sourceID, int(s.SessionType), 1, false)
		if err != nil {
			return err
		}
		conversation.ConversationID = conversationID
		if list == nil {
			conversation.LatestMsg = ""
			conversation.LatestMsgSendTime = s.SendTime
		} else {
			copier.Copy(&latestMsg, list[0])
			err := c.msgConvert(&latestMsg)
			if err != nil {
				log.Error("", "Parsing data error:", err.Error(), latestMsg)
			}
			conversation.LatestMsg = utils.StructToJsonString(latestMsg)
			conversation.LatestMsgSendTime = latestMsg.SendTime
		}
		err = c.db.UpdateColumnsConversation(ctx, conversation.ConversationID, map[string]interface{}{"latest_msg_send_time": conversation.LatestMsgSendTime, "latest_msg": conversation.LatestMsg})
		if err != nil {
			log.Error("internal", "updateConversationLatestMsgModel err: ", err)
		} else {
			_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
		}
	}
	return nil
}

func (c *Conversation) delMsgBySeq(seqList []uint32) error {
	var SPLIT = 1000
	for i := 0; i < len(seqList)/SPLIT; i++ {
		if err := c.delMsgBySeqSplit(seqList[i*SPLIT : (i+1)*SPLIT]); err != nil {
			return utils.Wrap(err, "")
		}
	}
	return nil
}

func (c *Conversation) delMsgBySeqSplit(seqList []uint32) error {
	var req server_api_params.DelMsgListReq
	req.SeqList = seqList
	req.OperationID = utils.OperationIDGenerator()
	req.OpUserID = c.loginUserID
	req.UserID = c.loginUserID
	operationID := req.OperationID

	err := c.SendReqWaitResp(context.Background(), &req, constant.WsDelMsg, 30, c.loginUserID)
	if err != nil {
		return utils.Wrap(err, "SendReqWaitResp failed")
	}
	var delResp server_api_params.DelMsgListResp
	err = proto.Unmarshal(resp.Data, &delResp)
	if err != nil {
		log.Error(operationID, "Unmarshal failed ", err.Error())
		return utils.Wrap(err, "Unmarshal failed")
	}
	return nil
}

// old WS method
func (c *Conversation) deleteMessageFromSvr(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, operationID string) {
	seq, err := c.db.GetMsgSeqByClientMsgID(s.ClientMsgID)
	common.CheckDBErrCallback(callback, err, operationID)
	if seq == 0 {
		err = errors.New("seq == 0 ")
		common.CheckArgsErrCallback(callback, err, operationID)
	}
	seqList := []uint32{seq}
	err = c.delMsgBySeq(seqList)
	common.CheckArgsErrCallback(callback, err, operationID)
}

// To delete session information, delete the server first, and then invoke the interface.
// The client receives a callback to delete all local information.
func (c *Conversation) deleteConversationAndMsgFromSvr(ctx context.Context, conversationID string) error {
	//校验这个会话是否存在，防止客户端删除不存在的会话
	_, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	//由于它是删除一个会话，所以是该会话全量消息的删除，所以这里不需要传递seqList
	var apiReq pbMsg.ClearConversationsMsgReq
	apiReq.UserID = c.loginUserID
	apiReq.ConversationIDs = []string{conversationID}
	return util.ApiPost(ctx, constant.DeleteConversationMsgRouter, &apiReq, nil)
}

func (c *Conversation) deleteAllMsgFromLocal(ctx context.Context) error {
	//log.NewInfo(operationID, utils.GetSelfFuncName())
	err := c.db.DeleteAllMessage(ctx)
	if err != nil {
		return err
	}
	groupIDList, err := c.full.GetReadDiffusionGroupIDList(ctx)
	if err != nil {
		return err
	}
	for _, v := range groupIDList {
		err = c.db.SuperGroupDeleteAllMessage(ctx, v)
		if err != nil {
			//log.Error(operationID, "SuperGroupDeleteAllMessage err", err.Error())
			continue
		}
	}
	err = c.db.ClearAllConversation(ctx)
	if err != nil {
		return err
	}
	conversationList, err := c.db.GetAllConversationListDB(ctx)
	if err != nil {
		return err
	}
	var cidList []string
	for _, conversation := range conversationList {
		cidList = append(cidList, conversation.ConversationID)
	}
	_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.ConChange, Args: cidList}, c.GetCh())
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.TotalUnreadMessageChanged, ""}})
	return nil
}
