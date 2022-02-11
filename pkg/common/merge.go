package common

import (
	"open_im_sdk/pkg/db"
	api "open_im_sdk/pkg/server_api_params"
)

func MergeResult(base []*db.LocalFriend, add []*db.LocalBlack) []api.FullUserInfo {
	m := make(map[string]api.FullUserInfo, 0)
	for _, v := range base {
		node := api.FullUserInfo{}
		node.FriendInfo = v
		m[v.FriendUserID] = node
	}

	for _, v := range add {
		if t, ok := m[v.BlockUserID]; ok {
			t.BlackInfo = v
		}
	}

	r := make([]api.FullUserInfo, 0)
	for _, v := range m {
		r = append(r, v)
	}
	return r
}
