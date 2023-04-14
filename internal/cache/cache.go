package cache

import (
	"context"
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/user"
	"open_im_sdk/pkg/db/model_struct"
	"sync"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
)

type UserInfo struct {
	Nickname string
	FaceURL  string
}
type Cache struct {
	user            *user.User
	friend          *friend.Friend
	userMap         sync.Map
	conversationMap sync.Map
}

func NewCache(user *user.User, friend *friend.Friend) *Cache {
	return &Cache{user: user, friend: friend}
}

func (c *Cache) Update(userID, faceURL, nickname string) {
	c.userMap.Store(userID, UserInfo{FaceURL: faceURL, Nickname: nickname})
}
func (c *Cache) UpdateConversation(conversation model_struct.LocalConversation) {
	c.conversationMap.Store(conversation.ConversationID, conversation)
}
func (c *Cache) UpdateConversations(conversations []*model_struct.LocalConversation) {
	for _, conversation := range conversations {
		c.conversationMap.Store(conversation.ConversationID, *conversation)
	}
}
func (c *Cache) GetAllConversations() (conversations []*model_struct.LocalConversation) {
	c.conversationMap.Range(func(key, value interface{}) bool {
		temp := value.(model_struct.LocalConversation)
		conversations = append(conversations, &temp)
		return true
	})
	return conversations
}
func (c *Cache) GetAllHasUnreadMessageConversations() (conversations []*model_struct.LocalConversation) {
	c.conversationMap.Range(func(key, value interface{}) bool {
		temp := value.(model_struct.LocalConversation)
		if temp.UnreadCount > 0 {
			conversations = append(conversations, &temp)
		}
		return true
	})
	return conversations
}

func (c *Cache) GetConversation(conversationID string) model_struct.LocalConversation {
	var result model_struct.LocalConversation
	conversation, ok := c.conversationMap.Load(conversationID)
	if ok {
		result = conversation.(model_struct.LocalConversation)
	}
	return result
}

func (c *Cache) GetUserNameAndFaceURL(ctx context.Context, userID string) (faceURL, name string, err error) {
	//find in cache
	if value, ok := c.userMap.Load(userID); ok {
		info := value.(UserInfo)
		return info.FaceURL, info.Nickname, nil
	}
	//get from local db
	friendInfo, err := c.friend.Db().GetFriendInfoByFriendUserID(ctx, userID)
	if err == nil {
		faceURL = friendInfo.FaceURL
		if friendInfo.Remark != "" {
			name = friendInfo.Remark
		} else {
			name = friendInfo.Nickname
		}
		return faceURL, name, nil
	}
	//get from server db
	users, err := c.user.GetServerUserInfo(ctx, []string{userID})
	if err != nil {
		return "", "", err
	}
	if len(users) == 0 {
		return "", "", errs.ErrUserIDNotFound.Wrap(userID)
	}
	c.userMap.Store(userID, UserInfo{FaceURL: users[0].FaceURL, Nickname: users[0].Nickname})
	return users[0].FaceURL, users[0].Nickname, nil
}
