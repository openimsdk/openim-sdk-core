package db

import (
	"errors"
	"open_im_sdk/pkg/utils"
)

func (u *open_im_sdk.UserRelated) _getBlackList() ([]*LocalBlack, error) {
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
func (u *open_im_sdk.UserRelated) _getBlackListUid() (blackListUid []string, err error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	return blackListUid, utils.wrap(u.imdb.Model(&LocalBlack{}).Select("block_user_id").Find(&blackListUid).Error, "_getBlackList failed")
}

func (u *open_im_sdk.UserRelated) _getBlackInfoByBlockUserID(blockUserID string) (*LocalBlack, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var black LocalBlack
	return &black, utils.wrap(u.imdb.Where("owner_user_id = ? AND block_user_id = ? ",
		u.loginUserID, blockUserID).Find(&black).Error, "_getBlackInfoByBlockUserID failed")
}

func (u *open_im_sdk.UserRelated) _getBlackInfoList(blockUserIDList []string) ([]LocalBlack, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var black []LocalBlack
	return black, utils.wrap(u.imdb.Where("block_user_id IN ? ", blockUserIDList).Find(&black).Error, "_getBlackInfoList failed")
}

func (u *open_im_sdk.UserRelated) _insertBlack(black *LocalBlack) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.wrap(u.imdb.Create(black).Error, "_insertBlack failed")
}

func (u *open_im_sdk.UserRelated) _updateBlack(black *LocalBlack) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(black)
	if t.RowsAffected == 0 {
		return utils.wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.wrap(t.Error, "_updateBlack failed")
}

func (u *open_im_sdk.UserRelated) _deleteBlack(blockUserID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.wrap(u.imdb.Where("owner_user_id=? and block_user_id=?", u.loginUserID, blockUserID).Delete(&LocalBlack{}).Error, "_delBlack failed")
}
