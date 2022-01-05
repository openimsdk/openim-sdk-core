package open_im_sdk

func (u *UserRelated) insertGroupData(groupInfo *Group) error {
	return Wrap(u.imdb.Create(groupInfo).Error, "insertIntoGroup failed")
}
func (u *UserRelated) deleteGroupData(groupInfo *Group) error {
	return Wrap(u.imdb.Delete(&groupInfo).Error, "deleteGroup failed")
}
func (u *UserRelated) updateGroupData() {

}
