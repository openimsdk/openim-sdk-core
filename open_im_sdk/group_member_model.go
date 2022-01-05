package open_im_sdk

import "errors"

func (u *UserRelated) _getGroupMemberInfoByGroupIDUserID(groupID, userID string) (*GroupMember, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var groupMember GroupMember
	return &groupMember, wrap(u.imdb.Where("group_id = ? AND user_id = ?",
		groupID, userID).Error, "_getGroupMemberInfoByGroupIDUserID failed")
}

func (u *UserRelated) _getAllGroupMemberList() ([]GroupMember, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var groupMemberList []GroupMember
	return groupMemberList, wrap(u.imdb.Find(&groupMemberList).Error, "_getAllGroupMemberList failed")
}

func (u *UserRelated) _getGroupMemberListByGroupID(groupID string) ([]GroupMember, error) {
	u.mRWMutex.RLock()
	defer u.mRWMutex.RUnlock()
	var groupMemberList []GroupMember
	return groupMemberList, wrap(u.imdb.Where("group_id = ? ", groupID).Find(&groupMemberList).Error, "_getGroupMemberListByGroupID failed")
}

func (u *UserRelated) _insertGroupMember(groupMember *GroupMember) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return u.imdb.Create(groupMember).Error
}

func (u *UserRelated) _deleteGroupMember(groupID, userID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	groupMember := GroupMember{GroupID: groupID, UserID: userID}
	return u.imdb.Delete(&groupMember).Error
}

func (u *UserRelated) _updateGroupMember(groupMember *GroupMember) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(groupMember)
	if t.RowsAffected == 0 {
		return wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return t.Error
}
