package cache

import (
	"errors"
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/user"
	"open_im_sdk/pkg/utils"
	"sync"
)

type UserInfo struct {
	Nickname string
	faceURL  string
}
type Cache struct {
	user    *user.User
	friend  *friend.Friend
	userMap sync.Map
}

func NewCache(user *user.User, friend *friend.Friend) *Cache {
	return &Cache{user: user, friend: friend}
}

func (c *Cache) Update(userID, faceURL, nickname string) {
	c.userMap.Store(userID, UserInfo{faceURL: faceURL, Nickname: nickname})
}

func (c *Cache) GetUserNameAndFaceURL(userID string, operationID string) (faceURL, name string, err error) {
	//find in cache
	user, ok := c.userMap.Load(userID)
	if ok {
		faceURL = user.(UserInfo).faceURL
		name = user.(UserInfo).Nickname
		return faceURL, name, nil
	}

	//get from local db
	friendInfo, err := c.friend.Db().GetFriendInfoByFriendUserID(userID)
	if err == nil {
		faceURL = friendInfo.FaceURL
		if friendInfo.Remark != "" {
			name = friendInfo.Remark
		} else {
			name = friendInfo.Nickname
		}
		return faceURL, name, nil
	}

	if operationID == "" {
		operationID = utils.OperationIDGenerator()
	}

	userInfos, err := c.user.GetUsersInfoFromCacheSvr([]string{userID}, operationID)
	if err != nil {
		return "", "", err
	}
	for _, v := range userInfos {
		faceURL = v.FaceURL
		name = v.Nickname
		c.userMap.Store(userID, UserInfo{faceURL: faceURL, Nickname: name})
		return v.FaceURL, v.Nickname, nil
	}
	return "", "", errors.New("no user ")
}
