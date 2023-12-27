package conversation_msg

import (
	"context"
	"encoding/json"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	utils2 "github.com/OpenIMSDK/tools/utils"
	"github.com/jinzhu/copier"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	_ int = iota
	stateCodeSuccess
	stateCodeEnd
)

const (
	intervalsTime = time.Second * 10
)

func newEntering(c *Conversation) *entering {
	e := &entering{
		conv:  c,
		send:  cache.New(intervalsTime, intervalsTime+intervalsTime/2),
		state: cache.New(intervalsTime, intervalsTime),
	}
	e.platformIDs = make([]int32, 0, len(constant.PlatformID2Name))
	e.platformIDSet = make(map[int32]struct{})
	for id := range constant.PlatformID2Name {
		e.platformIDSet[int32(id)] = struct{}{}
		e.platformIDs = append(e.platformIDs, int32(id))
	}
	utils2.Sort(e.platformIDs, true)
	e.state.OnEvicted(func(key string, val interface{}) {
		var data inputStatesKey
		if err := json.Unmarshal([]byte(key), &data); err != nil {
			return
		}
		e.changes(data.UserID, data.GroupID)
	})
	return e
}

type entering struct {
	send  *cache.Cache
	state *cache.Cache

	conv *Conversation

	platformIDs   []int32
	platformIDSet map[int32]struct{}
}

func (e *entering) InputState(ctx context.Context, conversationID string, focus bool) error {
	if conversationID == "" {
		return errs.ErrArgs.Wrap("conversationID can't be empty")
	}
	conversation, err := e.conv.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	key := conversation.ConversationID
	if focus {
		if val, ok := e.send.Get(key); ok {
			if val.(int) == stateCodeSuccess {
				log.ZDebug(ctx, "entering stateCodeSuccess", "conversationID", conversationID, "focus", focus)
				return nil
			}
		}
		e.send.SetDefault(key, stateCodeSuccess)
	} else {
		if val, ok := e.send.Get(key); ok {
			if val.(int) == stateCodeEnd {
				log.ZDebug(ctx, "entering stateCodeEnd", "conversationID", conversationID, "focus", focus)
				return nil
			}
			e.send.SetDefault(key, stateCodeEnd)
		} else {
			log.ZDebug(ctx, "entering send not found", "conversationID", conversationID, "focus", focus)
			return nil
		}
	}
	ctx, cancel := context.WithTimeout(ctx, intervalsTime/2)
	defer cancel()
	if err := e.sendMsg(ctx, conversation, focus); err != nil {
		e.send.Delete(key)
		return err
	}
	return nil
}

func (e *entering) sendMsg(ctx context.Context, conversation *model_struct.LocalConversation, focus bool) error {
	s := sdk_struct.MsgStruct{}
	err := e.conv.initBasicInfo(ctx, &s, constant.UserMsgType, constant.Entering)
	if err != nil {
		return err
	}
	s.RecvID = conversation.UserID
	s.GroupID = conversation.GroupID
	s.SessionType = conversation.ConversationType
	enteringElem := sdk_struct.EnteringElem{
		Focus: focus,
	}
	s.Content = utils.StructToJsonString(enteringElem)
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
		log.ZError(ctx, "entering msg to server failed", err, "message", s)
		return err
	}
	return nil
}

type inputStatesKey struct {
	UserID     string `json:"uid,omitempty"`
	GroupID    string `json:"gid,omitempty"`
	PlatformID int32  `json:"pid,omitempty"`
}

func (e *entering) getStateKey(platformID int32, userID string, groupID string) string {
	data, err := json.Marshal(inputStatesKey{PlatformID: platformID, UserID: userID, GroupID: groupID})
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (e *entering) onNewMsg(ctx context.Context, msg *sdk_struct.MsgStruct) {
	if msg.EnteringElem == nil {
		return
	}
	if msg.SendID == e.conv.loginUserID {
		return
	}
	if _, ok := e.platformIDSet[msg.SenderPlatformID]; !ok {
		return
	}
	now := time.Now().UnixMilli()
	expirationTimestamp := msg.SendTime + int64(intervalsTime/time.Millisecond)
	if msg.SendTime > now || expirationTimestamp <= now {
		return
	}
	key := e.getStateKey(msg.SenderPlatformID, msg.SendID, msg.GroupID)
	if msg.EnteringElem.Focus {
		d := time.Duration(expirationTimestamp - now)
		if v, t, ok := e.state.GetWithExpiration(key); ok {
			if t.UnixMilli() >= expirationTimestamp {
				return
			}
			e.state.Set(key, v, d)
		} else {
			e.state.Set(key, struct{}{}, d)
		}
		e.changes(msg.SendID, msg.GroupID)
	} else {
		if _, ok := e.state.Get(key); ok {
			e.state.Delete(key)
		}
	}
}

type InputStatesChangedData struct {
	UserID      string  `json:"userID"`
	GroupID     string  `json:"groupID"`
	PlatformIDs []int32 `json:"platformIDs"`
}

func (e *entering) changes(userID string, groupID string) {
	data := utils.StructToJsonString(e.GetInputStatesInfo(userID, groupID))
	e.conv.userListener().OnUserInputStatusChanged(data)
}

func (e *entering) GetInputStatesInfo(userID string, groupID string) *InputStatesChangedData {
	data := InputStatesChangedData{UserID: userID, GroupID: groupID}
	for _, platformID := range e.platformIDs {
		key := e.getStateKey(platformID, userID, groupID)
		if _, ok := e.state.Get(key); ok {
			data.PlatformIDs = append(data.PlatformIDs, platformID)
		}
	}
	return &data
}
