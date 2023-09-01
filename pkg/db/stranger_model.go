package db

import (
	"context"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) GetStrangerInfo(ctx context.Context, userIDs []string) ([]*model_struct.LocalStranger, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendList []model_struct.LocalStranger
	err := utils.Wrap(d.conn.WithContext(ctx).Where("user_id IN ?", userIDs).Find(&friendList).Error, "GetFriendInfoListByFriendUserID failed")
	var transfer []*model_struct.LocalStranger
	for _, v := range friendList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, err
}

func (d *DataBase) SetStrangerInfo(ctx context.Context, localStrangerList []*model_struct.LocalStranger) error {

}
