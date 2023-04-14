package full

import (
	"context"
	"open_im_sdk/pkg/db/model_struct"
)

func (u *Full) GetGroupInfoByGroupID(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
	g1, err := u.SuperGroup.GetGroupInfoFromLocal2Svr(ctx, groupID)
	if err == nil {
		return g1, nil
	}
	g2, err := u.group.GetGroupInfoFromLocal2Svr(ctx, groupID)
	return g2, err
}
