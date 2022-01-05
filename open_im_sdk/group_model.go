package open_im_sdk

import "errors"

func (u *UserRelated) _insertGroup(groupInfo *LocalGroup) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return Wrap(u.imdb.Create(groupInfo).Error, "_insertGroup failed")
}
func (u *UserRelated) _deleteGroup(groupInfo *LocalGroup) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return Wrap(u.imdb.Delete(&groupInfo).Error, "_deleteGroup failed")
}
func (u *UserRelated) _updateGroup(groupInfo *LocalGroup) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(groupInfo)
	if t.RowsAffected == 0 {
		return wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return wrap(t.Error, "_updateGroup failed")
}
func (u *UserRelated) _getGroupList() ([]LocalGroup, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var groupList []LocalGroup
	return groupList, wrap(u.imdb.Find(&groupList).Error, "_getGroupList failed")
}
func (u *UserRelated) _getGroupInfoByGroupID(groupID string) (g *LocalGroup, err error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return g, wrap(u.imdb.Where("group_id = ?", groupID).Find(g).Error, "_getGroupList failed")
}
