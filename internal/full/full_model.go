package full

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (u *Full) GetGroupInfoByGroupID(groupID string) (*model_struct.LocalGroup, error) {
	t, err := u.db.GetGroupType(groupID)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	if t == constant.NormalGroup {
		return u.db.GetGroupInfoByGroupID(groupID)
	}
	return u.db.GetSuperGroupInfoByGroupID(groupID)
}
