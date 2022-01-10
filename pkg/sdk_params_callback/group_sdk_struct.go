package sdk_params_callback

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/server_api_params"
)

type CreateGroupBaseInfoParam struct {
	GroupName   string                `json:"groupName"`
	GroupType   int32                 `json:"groupType"`
}


type CreateGroupMemberRoleParam []*server_api_params.GroupAddMemberInfo
type CreateGroupCallback struct {
	server_api_params.GroupInfo
}



//param groupID reqMsg
const JoinGroupCallback = constant.SuccessCallbackDefault


//type QuitGroupParam // groupID
const QuitGroupCallback  = constant.SuccessCallbackDefault

//type GetJoinedGroupListParam null
type GetJoinedGroupListCallback []*db.LocalGroup


type GetGroupsInfoParam []string
type GetGroupsInfoCallback []*db.LocalGroup


type SetGroupInfoParam struct {
	GroupName    string `json:"groupName"`
	Notification string `json:"notification"`
	Introduction string `json:"introduction"`
	FaceUrl      string `json:"faceUrl"`
	Ex string `json:"ex"`
}
const SetGroupInfoCallback = constant.SuccessCallbackDefault

//type GetGroupMemberListParam groupID ...
type GetGroupMemberListCallback struct{
	MemberList [] *db.LocalGroupMember `json:"data"`
	NextSeq int32                 `json:"nextSeq"`
}

type GetGroupMembersInfoParam []string
type GetGroupMembersInfoCallback []*db.LocalGroupMember


type KickGroupMemberParam []string
type KickGroupMemberCallback []*server_api_params.UserIDResult


//type TransferGroupOwnerParam
const TransferGroupOwnerCallback = constant.SuccessCallbackDefault

type InviteUserToGroupParam []string
type InviteUserToGroupCallback []*server_api_params.UserIDResult

//type GetGroupApplicationListParam
type GetGroupApplicationListCallback []*db.LocalGroupRequest

//type AcceptGroupApplicationParam
const AcceptGroupApplicationCallback = constant.SuccessCallbackDefault

//type RefuseGroupApplicationParam
const RefuseGroupApplicationCallback = constant.SuccessCallbackDefault




