package open_im_sdk

import "errors"

func (u *UserRelated) _getBlackList() ([]Black, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var blackList []Black
	return blackList, wrap(u.imdb.Find(&blackList).Error, "_getBlackList failed")
}

func (u *UserRelated) _getBlackInfoByBlockUserID(blockUserID string) (*Black, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var black Black
	return &black, wrap(u.imdb.Where("owner_user_id = ? AND block_user_id = ? ",
		u.loginUserID, blockUserID).Find(&black).Error, "_getBlackInfoByBlockUserID failed")
}

func (u *UserRelated) _insertBlack(black *Black) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return wrap(u.imdb.Create(black).Error, "_insertBlack failed")
}

func (u *UserRelated) _updateBlack(black *Black) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(black)
	if t.RowsAffected == 0 {
		return wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return wrap(t.Error, "_updateBlack failed")
}

func (u *UserRelated) _delBlack(blockUserID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	black := Black{OwnerUserID: u.loginUserID, BlockUserID: blockUserID}
	return wrap(u.imdb.Delete(&black).Error, "_delBlack failed")
}
