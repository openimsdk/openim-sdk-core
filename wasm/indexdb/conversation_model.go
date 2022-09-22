package indexdb

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/indexdb/temp_struct"
)

type LocalConversations struct {
}

func (i *LocalConversations) GetAllConversationList() (result []*model_struct.LocalConversation, err error) {
	cList, err := Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := cList.(string); ok {
			var temp []model_struct.LocalConversation
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}
func (i *LocalConversations) GetConversation(conversationID string) (*model_struct.LocalConversation, error) {
	c, err := Exec(conversationID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := c.(string); ok {
			result := model_struct.LocalConversation{}
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return &result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i *LocalConversations) GetHiddenConversationList() (result []*model_struct.LocalConversation, err error) {
	cList, err := Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := cList.(string); ok {
			var temp []model_struct.LocalConversation
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}
func (i *LocalConversations) GetAllConversationListToSync() (result []*model_struct.LocalConversation, err error) {
	cList, err := Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := cList.(string); ok {
			var temp []model_struct.LocalConversation
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}
func (i *LocalConversations) UpdateColumnsConversation(conversationID string, args map[string]interface{}) error {
	_, err := Exec(conversationID, args)
	return err
}
func (i IndexDB) GetConversationByUserID(userID string) (*model_struct.LocalConversation, error) {
	c, err := Exec(userID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := c.(string); ok {
			result := model_struct.LocalConversation{}
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return &result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i IndexDB) GetConversationListSplit(offset, count int) (result []*model_struct.LocalConversation, err error) {
	cList, err := Exec(offset, count)
	if err != nil {
		return nil, err
	} else {
		if v, ok := cList.(string); ok {
			var temp []model_struct.LocalConversation
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i IndexDB) BatchInsertConversationList(conversationList []*model_struct.LocalConversation) error {
	_, err := Exec(utils.StructToJsonString(conversationList))
	return err
}

func (i IndexDB) InsertConversation(conversationList *model_struct.LocalConversation) error {
	_, err := Exec(utils.StructToJsonString(conversationList))
	return err
}

func (i IndexDB) DeleteConversation(conversationID string) error {
	_, err := Exec(conversationID)
	return err
}

func (i IndexDB) UpdateConversation(c *model_struct.LocalConversation) error {
	if c.ConversationID == "" {
		return PrimaryKeyNull
	}
	tempLocalConversation := temp_struct.LocalConversation{
		ConversationType:      c.ConversationType,
		UserID:                c.UserID,
		GroupID:               c.GroupID,
		ShowName:              c.ShowName,
		FaceURL:               c.FaceURL,
		RecvMsgOpt:            c.RecvMsgOpt,
		UnreadCount:           c.UnreadCount,
		GroupAtType:           c.GroupAtType,
		LatestMsg:             c.LatestMsg,
		LatestMsgSendTime:     c.LatestMsgSendTime,
		DraftText:             c.DraftText,
		DraftTextTime:         c.DraftTextTime,
		IsPinned:              c.IsPinned,
		IsPrivateChat:         c.IsPrivateChat,
		IsNotInGroup:          c.IsNotInGroup,
		UpdateUnreadCountTime: c.UpdateUnreadCountTime,
		AttachedInfo:          c.AttachedInfo,
		Ex:                    c.Ex,
	}
	_, err := Exec(c.ConversationID, utils.StructToJsonString(tempLocalConversation))
	return err
}

func (i IndexDB) UpdateConversationForSync(c *model_struct.LocalConversation) error {
	if c.ConversationID == "" {
		return PrimaryKeyNull
	}
	tempLocalConversation := temp_struct.LocalConversation{
		RecvMsgOpt:            c.RecvMsgOpt,
		UnreadCount:           c.UnreadCount,
		GroupAtType:           c.GroupAtType,
		LatestMsg:             c.LatestMsg,
		LatestMsgSendTime:     c.LatestMsgSendTime,
		DraftText:             c.DraftText,
		DraftTextTime:         c.DraftTextTime,
		IsPinned:              c.IsPinned,
		IsPrivateChat:         c.IsPrivateChat,
		IsNotInGroup:          c.IsNotInGroup,
		UpdateUnreadCountTime: c.UpdateUnreadCountTime,
		AttachedInfo:          c.AttachedInfo,
		Ex:                    c.Ex,
	}
	_, err := Exec(c.ConversationID, utils.StructToJsonString(tempLocalConversation))
	return err
}

func (i IndexDB) BatchUpdateConversationList(conversationList []*model_struct.LocalConversation) error {
	panic("implement me")
}

func (i IndexDB) ConversationIfExists(conversationID string) (bool, error) {
	panic("implement me")
}

func (i IndexDB) ResetConversation(conversationID string) error {
	panic("implement me")
}

func (i IndexDB) ResetAllConversation() error {
	panic("implement me")
}

func (i IndexDB) ClearConversation(conversationID string) error {
	panic("implement me")
}

func (i IndexDB) CleaAllConversation() error {
	panic("implement me")
}

func (i IndexDB) SetConversationDraft(conversationID, draftText string) error {
	panic("implement me")
}

func (i IndexDB) RemoveConversationDraft(conversationID, draftText string) error {
	panic("implement me")
}

func (i IndexDB) UnPinConversation(conversationID string, isPinned int) error {
	panic("implement me")
}

func (i IndexDB) UpdateAllConversation(conversation *model_struct.LocalConversation) error {
	panic("implement me")
}

func (i IndexDB) IncrConversationUnreadCount(conversationID string) error {
	panic("implement me")
}

func (i IndexDB) GetTotalUnreadMsgCount() (totalUnreadCount int32, err error) {
	panic("implement me")
}

func (i IndexDB) SetMultipleConversationRecvMsgOpt(conversationIDList []string, opt int) (err error) {
	panic("implement me")
}

func (i IndexDB) GetMultipleConversation(conversationIDList []string) (result []*model_struct.LocalConversation, err error) {
	panic("implement me")
}
