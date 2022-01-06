package open_im_sdk

import "errors"

func (u *UserRelated) _getBlackList() ([]*LocalBlack, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var blackList []LocalBlack

	err := u.imdb.Find(&blackList).Error
	var transfer []*LocalBlack
	for _, v := range blackList {
		transfer = append(transfer, &v)
	}
	return transfer, err

}
func (u *UserRelated) _getBlackListUid() (blackListUid []string, err error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	return blackListUid, wrap(u.imdb.Select("block_user_id").Find(&blackListUid).Error, "_getBlackList failed")
}

func (u *UserRelated) _getBlackInfoByBlockUserID(blockUserID string) (*LocalBlack, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var black LocalBlack
	return &black, wrap(u.imdb.Where("owner_user_id = ? AND block_user_id = ? ",
		u.loginUserID, blockUserID).Find(&black).Error, "_getBlackInfoByBlockUserID failed")
}

func (u *UserRelated) _getBlackInfoList(blockUserIDList []string) ([]LocalBlack, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var black []LocalBlack
	return black, wrap(u.imdb.Where("block_user_id IN ? ", blockUserIDList).Find(&black).Error, "_getBlackInfoList failed")
}

func (u *UserRelated) _insertBlack(black *LocalBlack) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return wrap(u.imdb.Create(black).Error, "_insertBlack failed")
}

func (u *UserRelated) _updateBlack(black *LocalBlack) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(black)
	if t.RowsAffected == 0 {
		return wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return wrap(t.Error, "_updateBlack failed")
}

func (u *UserRelated) _deleteBlack(blockUserID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	black := LocalBlack{OwnerUserID: u.loginUserID, BlockUserID: blockUserID}
	return wrap(u.imdb.Delete(&black).Error, "_delBlack failed")
}
