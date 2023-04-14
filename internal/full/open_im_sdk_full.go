package full

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"open_im_sdk/pkg/db/model_struct"
	api "open_im_sdk/pkg/server_api_params"
)

func (u *Full) GetUsersInfo(ctx context.Context, userIDList []string) ([]*api.FullUserInfo, error) {
	friendList, err := u.db.GetFriendInfoList(ctx, userIDList)
	if err != nil {
		return nil, err
	}
	blackList, err := u.db.GetBlackInfoList(ctx, userIDList)
	if err != nil {
		return nil, err
	}
	users, err := u.user.GetServerUserInfo(ctx, userIDList)
	if err != nil {
		return nil, err
	}
	friendMap := make(map[string]*model_struct.LocalFriend)
	for i, f := range friendList {
		friendMap[f.FriendUserID] = friendList[i]
	}
	blackMap := make(map[string]*model_struct.LocalBlack)
	for i, b := range blackList {
		blackMap[b.BlockUserID] = blackList[i]
	}
	userMap := make(map[string]*sdkws.UserInfo)
	for i, info := range users {
		userMap[info.UserID] = users[i]
	}
	res := make([]*api.FullUserInfo, 0, len(users))
	for _, userID := range userIDList {
		info, ok := userMap[userID]
		if !ok {
			continue
		}
		res = append(res, &api.FullUserInfo{
			PublicInfo: info,
			FriendInfo: friendMap[userID],
			BlackInfo:  blackMap[userID],
		})
	}
	return res, nil
}
