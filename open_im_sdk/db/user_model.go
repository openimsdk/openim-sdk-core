package db

import (
	"errors"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/open_im_sdk/utils"
)

func (u *open_im_sdk.UserRelated) _getLoginUser() (*open_im_sdk.LocalUser, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var user open_im_sdk.LocalUser
	return &user, utils.wrap(u.imdb.First(&user).Error, "_getLoginUserInfo failed")
}

func (u *open_im_sdk.UserRelated) _updateLoginUser(user *open_im_sdk.LocalUser) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(user)
	if t.RowsAffected == 0 {
		return utils.wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.wrap(t.Error, "_updateLoginUser failed")
}

func (u *open_im_sdk.UserRelated) _insertLoginUser(user *open_im_sdk.LocalUser) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.wrap(u.imdb.Create(user).Error, "_insertLoginUser failed")
}
