package full

import (
	"open_im_sdk/pkg/db/model_struct"
)

func (u *Full) GetGroupInfoByGroupID(groupID string) (*model_struct.LocalGroup, error) {
	g1, err := u.SuperGroup.GetGroupInfoFromLocal2Svr(groupID)
	if err == nil {
		return g1, nil
	}
	g2, err := u.group.GetGroupInfoFromLocal2Svr(groupID)
	return g2, err
}
