package open_im_sdk

import "errors"

func (u *UserRelated) _getLoginUser() (*User, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var user User
	return &user, wrap(u.imdb.First(&user).Error, "_getLoginUserInfo failed")
}

func (u *UserRelated) _updateLoginUser(user *User) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(user)
	if t.RowsAffected == 0 {
		return wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return wrap(t.Error, "_updateLoginUser failed")
}

func (u *UserRelated) _insertLoginUser(user *User) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return wrap(u.imdb.Create(user).Error, "_insertLoginUser failed")
}
