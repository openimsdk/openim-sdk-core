package indexdb

import "open_im_sdk/pkg/db/model_struct"

type LocalFriend struct {
	loginUserID string
}

func (i *LocalFriend) InsertFriend(friend *model_struct.LocalFriend) error {
	panic("implement me")
}

func (i *LocalFriend) DeleteFriend(friendUserID string, ownerUserID string) error {
	panic("implement me")
}

func (i *LocalFriend) UpdateFriend(friend *model_struct.LocalFriend) error {
	panic("implement me")
}

func (i *LocalFriend) GetAllFriendList() ([]*model_struct.LocalFriend, error) {
	panic("implement me")
}

func (i *LocalFriend) SearchFriendList(keyword string, isSearchUserID, isSearchNickname, isSearchRemark bool) ([]*model_struct.LocalFriend, error) {
	panic("implement me")
}

func (i *LocalFriend) GetFriendInfoByFriendUserID(FriendUserID string) (*model_struct.LocalFriend, error) {
	panic("implement me")
}

func (i *LocalFriend) GetFriendInfoList(friendUserIDList []string) ([]*model_struct.LocalFriend, error) {
	panic("implement me")
}
