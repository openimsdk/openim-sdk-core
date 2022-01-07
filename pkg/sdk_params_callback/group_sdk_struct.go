package sdk_params_callback

import (
	"open_im_sdk/pkg/server_api_params"
)

type CreateGroupBaseInfoParam struct {
	GroupName    string `json:"groupName"`
	Notification string `json:"notification"`
	Introduction string `json:"introduction"`
	FaceUrl      string `json:"faceUrl"`
	Ex           string `json:"ex"`
}
type UserRole struct {
	UserID    string `json:"userID"`
	RoleLevel int32  `json:"RoleLevel"`
}
type CreateGroupMemberRoleParam []UserRole
type CreateGroupCallback struct {
	server_api_params.GroupInfo
}
