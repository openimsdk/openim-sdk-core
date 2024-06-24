package group

import "github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"

type GetGroupMemberListV2Response struct {
	GroupMembersList []*model_struct.LocalGroupMember
	IsEnd            bool `json:"isEnd"`
}

type GetGroupListV2Response struct {
	GroupsList []*model_struct.LocalGroup
	IsEnd      bool `json:"isEnd"`
}
