package cache

import (
	"context"
	"errors"
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/user"
	"open_im_sdk/pkg/db/model_struct"
	"sync"
)

type UserInfo struct {
	nickname string
	faceURL  string
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
	c.userMap.Store(userID, UserInfo{faceURL: faceURL, nickname: nickname})
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
	user, ok := c.userMap.Load(userID)
	if ok {
		faceURL = user.(UserInfo).faceURL
		name = user.(UserInfo).nickname
		return faceURL, name, nil
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
	userInfos, err := c.user.GetUsersInfoFromSvr(ctx, []string{userID})
	if err != nil {
		return "", "", err
	}
	for _, v := range userInfos {
		faceURL = v.FaceURL
		name = v.Nickname
		c.userMap.Store(userID, UserInfo{faceURL: faceURL, nickname: name})
		return v.FaceURL, v.Nickname, nil
	}
	return "", "", errors.New("user not exist" + userID)
}
