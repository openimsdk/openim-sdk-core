package conversation_msg

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	pconstant "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	_ int = iota
	stateCodeSuccess
	stateCodeEnd
)

const (
	inputStatesSendTime   = time.Second * 10                            // input status sending interval time
	inputStatesTimeout    = inputStatesSendTime + inputStatesSendTime/2 // input status timeout
	inputStatesMsgTimeout = inputStatesSendTime / 2                     // message sending timeout
)

func newTyping(c *Conversation) *typing {
	e := &typing{
		conv:  c,
		send:  cache.New(inputStatesSendTime, inputStatesTimeout),
		state: cache.New(inputStatesTimeout, inputStatesTimeout),
	}
	e.platformIDs = make([]int32, 0, len(pconstant.PlatformID2Name))
	e.platformIDSet = make(map[int32]struct{})
	for id := range pconstant.PlatformID2Name {
		e.platformIDSet[int32(id)] = struct{}{}
		e.platformIDs = append(e.platformIDs, int32(id))
	}
	datautil.Sort(e.platformIDs, true)
	e.state.OnEvicted(func(key string, val interface{}) {
		var data inputStatesKey
		if err := json.Unmarshal([]byte(key), &data); err != nil {
			return
		}
		e.changes(data.ConversationID, data.UserID)
	})
	return e
}

type typing struct {
	send  *cache.Cache
	state *cache.Cache

	conv *Conversation

	platformIDs   []int32
	platformIDSet map[int32]struct{}
}

func (e *typing) ChangeInputStates(ctx context.Context, conversationID string, focus bool) error {
	if conversationID == "" {
		return errs.ErrArgs.WrapMsg("conversationID can't be empty")
	}
	conversation, err := e.conv.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	key := conversation.ConversationID
	if focus {
		if val, ok := e.send.Get(key); ok {
			if val.(int) == stateCodeSuccess {
				log.ZDebug(ctx, "typing stateCodeSuccess", "conversationID", conversationID, "focus", focus)
				return nil
			}
		}
		e.send.SetDefault(key, stateCodeSuccess)
	} else {
		if val, ok := e.send.Get(key); ok {
			if val.(int) == stateCodeEnd {
				log.ZDebug(ctx, "typing stateCodeEnd", "conversationID", conversationID, "focus", focus)
				return nil
			}
			e.send.SetDefault(key, stateCodeEnd)
		} else {
			log.ZDebug(ctx, "typing send not found", "conversationID", conversationID, "focus", focus)
			return nil
		}
	}
	ctx, cancel := context.WithTimeout(ctx, inputStatesMsgTimeout)
	defer cancel()
	if err := e.sendMsg(ctx, conversation, focus); err != nil {
		e.send.Delete(key)
		return err
	}
	return nil
}

func (e *typing) sendMsg(ctx context.Context, conversation *model_struct.LocalConversation, focus bool) error {
	s := sdk_struct.MsgStruct{}
	err := e.conv.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Typing)
	if err != nil {
		return err
	}
	s.RecvID = conversation.UserID
	s.GroupID = conversation.GroupID
	s.SessionType = conversation.ConversationType
	var typingElem sdk_struct.TypingElem
	if focus {
		typingElem.MsgTips = "yes"
	} else {
		typingElem.MsgTips = "no"
	}
	s.Content = utils.StructToJsonString(typingElem)
	options := make(map[string]bool, 6)
	utils.SetSwitchFromOptions(options, constant.IsHistory, false)
	utils.SetSwitchFromOptions(options, constant.IsPersistent, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderSync, false)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	var wsMsgData sdkws.MsgData
	copier.Copy(&wsMsgData, s)
	wsMsgData.Content = []byte(s.Content)
	wsMsgData.CreateTime = s.CreateTime
	wsMsgData.Options = options
	var sendMsgResp sdkws.UserSendMsgResp
	err = e.conv.LongConnMgr.SendReqWaitResp(ctx, &wsMsgData, constant.SendMsg, &sendMsgResp)
	if err != nil {
		log.ZError(ctx, "typing msg to server failed", err, "message", s)
		return err
	}
	return nil
}

type inputStatesKey struct {
	ConversationID string `json:"cid,omitempty"`
	UserID         string `json:"uid,omitempty"`
	PlatformID     int32  `json:"pid,omitempty"`
}

func (e *typing) getStateKey(conversationID string, userID string, platformID int32) string {
	data, err := json.Marshal(inputStatesKey{ConversationID: conversationID, UserID: userID, PlatformID: platformID})
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (e *typing) onNewMsg(ctx context.Context, msg *sdkws.MsgData) {
	var enteringElem sdk_struct.TypingElem
	if err := json.Unmarshal(msg.Content, &enteringElem); err != nil {
		log.ZError(ctx, "typing onNewMsg Unmarshal failed", err, "message", msg)
		return
	}
	if msg.SendID == e.conv.loginUserID {
		return
	}
	if _, ok := e.platformIDSet[msg.SenderPlatformID]; !ok {
		return
	}
	now := time.Now().UnixMilli()
	expirationTimestamp := now + int64(inputStatesSendTime/time.Millisecond)
	var sourceID string
	if msg.GroupID == "" {
		sourceID = msg.SendID
	} else {
		sourceID = msg.GroupID
	}
	conversationID := e.conv.getConversationIDBySessionType(sourceID, int(msg.SessionType))
	key := e.getStateKey(conversationID, msg.SendID, msg.SenderPlatformID)
	if enteringElem.MsgTips == "yes" {
		d := time.Duration(expirationTimestamp-now) * time.Millisecond
		if v, t, ok := e.state.GetWithExpiration(key); ok {
			if t.UnixMilli() >= expirationTimestamp {
				return
			}
			e.state.Set(key, v, d)
		} else {
			e.state.Set(key, struct{}{}, d)
			e.changes(conversationID, msg.SendID)
		}
	} else {
		if _, ok := e.state.Get(key); ok {
			e.state.Delete(key)
		}
	}
}

type InputStatesChangedData struct {
	ConversationID string  `json:"conversationID"`
	UserID         string  `json:"userID"`
	PlatformIDs    []int32 `json:"platformIDs"`
}

func (e *typing) changes(conversationID string, userID string) {
	data := InputStatesChangedData{ConversationID: conversationID, UserID: userID, PlatformIDs: e.GetInputStates(conversationID, userID)}
	e.conv.ConversationListener().OnConversationUserInputStatusChanged(utils.StructToJsonString(data))
}

func (e *typing) GetInputStates(conversationID string, userID string) []int32 {
	platformIDs := make([]int32, 0, 1)
	for _, platformID := range e.platformIDs {
		key := e.getStateKey(conversationID, userID, platformID)
		if _, ok := e.state.Get(key); ok {
			platformIDs = append(platformIDs, platformID)
		}
	}
	return platformIDs
}
