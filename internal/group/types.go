package group

import "github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"

type GroupMemberListWithIsEnd struct {
	GroupMemberList []*model_struct.LocalGroupMember
	IsEnd           bool `json:"isEnd"`
}

type GroupListWithIsEnd struct {
	GroupList []*model_struct.LocalGroup
	IsEnd     bool `json:"isEnd"`
}
