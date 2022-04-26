package sdk_params_callback

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/server_api_params"
)

type CreateGroupBaseInfoParam struct {
	GroupType int32 `json:"groupType"`
	SetGroupInfoParam
}

type CreateGroupMemberRoleParam []*server_api_params.GroupAddMemberInfo
type CreateGroupCallback server_api_params.GroupInfo

//param groupID reqMsg
const JoinGroupCallback = constant.SuccessCallbackDefault

//type QuitGroupParam // groupID
const QuitGroupCallback = constant.SuccessCallbackDefault

const DismissGroupCallback = constant.SuccessCallbackDefault

const GroupMuteChangeCallback = constant.SuccessCallbackDefault

const GroupMemberMuteChangeCallback = constant.SuccessCallbackDefault

const SetGroupMemberNicknameCallback = constant.SuccessCallbackDefault

//type GetJoinedGroupListParam null
type GetJoinedGroupListCallback []*db.LocalGroup

type GetGroupsInfoParam []string
type GetGroupsInfoCallback []*db.LocalGroup
type SearchGroupsParam struct {
	KeywordList       []string `json:"keywordList"`
	IsSearchGroupID   bool     `json:"isSearchGroupID"`
	IsSearchGroupName bool     `json:"isSearchGroupName"`
}
type SearchGroupsCallback []*db.LocalGroup

type SetGroupInfoParam struct {
	GroupName    string `json:"groupName"`
	Notification string `json:"notification"`
	Introduction string `json:"introduction"`
	FaceURL      string `json:"faceURL"`
	Ex           string `json:"ex"`
}

const SetGroupInfoCallback = constant.SuccessCallbackDefault

//type GetGroupMemberListParam groupID ...
type GetGroupMemberListCallback []*db.LocalGroupMember

type GetGroupMembersInfoParam []string
type GetGroupMembersInfoCallback []*db.LocalGroupMember

type KickGroupMemberParam []string
type KickGroupMemberCallback []*server_api_params.UserIDResult

//type TransferGroupOwnerParam
const TransferGroupOwnerCallback = constant.SuccessCallbackDefault

type InviteUserToGroupParam []string
type InviteUserToGroupCallback []*server_api_params.UserIDResult

//type GetGroupApplicationListParam
type GetGroupApplicationListCallback []*db.LocalAdminGroupRequest

type GetSendGroupApplicationListCallback []*db.LocalGroupRequest

//type AcceptGroupApplicationParam
const AcceptGroupApplicationCallback = constant.SuccessCallbackDefault

//type RefuseGroupApplicationParam
const RefuseGroupApplicationCallback = constant.SuccessCallbackDefault
