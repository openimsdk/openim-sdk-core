package conversation_msg

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

func (c *Conversation) parseStreamMsgTips(content []byte) (*sdkws.StreamMsgTips, error) {
	var notificationElem sdkws.NotificationElem
	if err := json.Unmarshal(content, &notificationElem); err != nil {
		return nil, err
	}
	var tips sdkws.StreamMsgTips
	if err := json.Unmarshal([]byte(notificationElem.Detail), &tips); err != nil {
		return nil, err
	}
	return &tips, nil
}

func (c *Conversation) doStreamMsgNotification(ctx context.Context, msg *sdkws.MsgData) error {
	tips, err := c.parseStreamMsgTips(msg.Content)
	if err != nil {
		return err
	}
	dbMsg, err := c.db.GetMessage(ctx, tips.ConversationID, tips.ClientMsgID)
	if err != nil {
		log.ZWarn(ctx, "get db stream msg failed", err, "tips", tips)
		return err
	}
	if dbMsg.ContentType != constant.Stream {
		return errors.New("content type is not stream")
	}
	var streamElem sdk_struct.StreamElem
	if err := json.Unmarshal([]byte(dbMsg.Content), &streamElem); err != nil {
		return err
	}
	if streamElem.End {
		log.ZWarn(ctx, "db stream msg is end", nil)
		return nil
	}
	if len(streamElem.Packets) < int(tips.StartIndex) {
		log.ZWarn(ctx, "db stream msg packets is not enough", nil, "streamElem", streamElem, "tips", tips)
		c.asyncStreamMsg(ctx, tips.ConversationID, tips.ClientMsgID)
		return nil
	}
	streamElem.Packets = streamElem.Packets[:tips.StartIndex]
	for _, packet := range tips.Packets {
		streamElem.Packets = append(streamElem.Packets, packet)
	}
	streamElem.End = tips.End
	data := utils.StructToJsonString(streamElem)
	if data == dbMsg.Content {
		log.ZDebug(ctx, "stream msg unchanged")
		return nil
	}
	dbMsg.Content = string(data)
	return c.setStreamMsg(ctx, tips.ConversationID, dbMsg)
}

func (c *Conversation) setStreamMsg(ctx context.Context, conversationID string, msg *model_struct.LocalChatLog) error {
	if err := c.db.UpdateMessage(ctx, conversationID, msg); err != nil {
		return err
	}
	//_, res := c.LocalChatLog2MsgStruct(ctx, []*model_struct.LocalChatLog{msg})
	//if len(res) == 0 {
	//	log.ZWarn(ctx, "LocalChatLog2MsgStruct failed", nil, "msg", msg)
	//	return nil
	//}
	//data := utils.StructToJsonString(res[0])
	//log.ZDebug(ctx, "setStreamMsg", "data", data)
	//c.msgListener().OnMsgEdited(data)
	return c.updateConversationLastMsg(ctx, conversationID, msg)
}

func (c *Conversation) updateConversationLastMsg(ctx context.Context, conversationID string, msg *model_struct.LocalChatLog) error {
	oc, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	if oc.LatestMsg != "" {
		var conversationMsg model_struct.LocalChatLog
		if err := json.Unmarshal([]byte(oc.LatestMsg), &conversationMsg); err != nil {
			return err
		}
		if conversationMsg.SendTime >= msg.SendTime && conversationMsg.ClientMsgID != msg.ClientMsgID {
			return nil
		}
	}
	//_, res := c.LocalChatLog2MsgStruct(ctx, []*model_struct.LocalChatLog{msg})
	//if len(res) == 0 {
	//	log.ZWarn(ctx, "LocalChatLog2MsgStruct failed", nil, "msg", msg)
	//	return nil
	//}
	oc.LatestMsgSendTime = msg.SendTime
	//oc.LatestMsg = utils.StructToJsonString(res[0])
	if err := c.db.UpdateConversation(ctx, oc); err != nil {
		return err
	}
	conversationData := utils.StructToJsonString([]*model_struct.LocalConversation{oc})
	log.ZDebug(ctx, "setStreamMsg conversation changed", "conversationData", conversationData)
	c.ConversationListener().OnConversationChanged(conversationData)
	return nil
}

func (c *Conversation) syncStreamMsg(ctx context.Context, conversationID string, clientMsgID string) error {
	c.streamMsgMutex.Lock()
	defer c.streamMsgMutex.Unlock()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	dbMsg, err := c.db.GetMessage(ctx, conversationID, clientMsgID)
	if err != nil {
		return err
	}
	if dbMsg.ContentType != constant.Stream {
		return errors.New("content type is not stream")
	}
	var streamElem sdk_struct.StreamElem
	if err := json.Unmarshal([]byte(dbMsg.Content), &streamElem); err != nil {
		return errs.WrapMsg(err, "unmarshal stream msg failed")
	}
	if streamElem.End {
		return nil
	}
	resp, err := c.getStreamMsg(ctx, clientMsgID)
	if err != nil {
		return err
	}
	streamElem.Packets = resp.Packets
	streamElem.End = resp.End
	content := utils.StructToJsonString(&streamElem)
	if dbMsg.Content == content {
		log.ZDebug(ctx, "stream msg unchanged", "conversationID", conversationID, "clientMsgID", clientMsgID, "streamElem", streamElem)
		return nil
	}
	dbMsg.Content = content
	return c.setStreamMsg(ctx, conversationID, dbMsg)
}

func (c *Conversation) asyncStreamMsg(ctx context.Context, conversationID string, clientMsgID string) {
	ctx = context.WithoutCancel(ctx)
	go func() {
		if err := c.syncStreamMsg(ctx, conversationID, clientMsgID); err != nil {
			log.ZError(ctx, "syncStreamMsg failed", err, "conversationID", conversationID, "clientMsgID", clientMsgID)
		}
	}()
}

func (c *Conversation) streamMsgReplace(ctx context.Context, conversationID string, msgs []*sdk_struct.MsgStruct) {
	for _, msg := range msgs {
		if msg.ContentType != constant.Stream {
			continue
		}
		var tips sdkws.StreamMsgTips
		if err := json.Unmarshal([]byte(msg.Content), &tips); err != nil {
			log.ZError(ctx, "unmarshal stream msg tips failed", err, "msg", msg)
			continue
		}
		if tips.End {
			continue
		}
		c.asyncStreamMsg(ctx, conversationID, msg.ClientMsgID)
	}
}
