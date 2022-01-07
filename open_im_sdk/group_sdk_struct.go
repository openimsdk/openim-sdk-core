package open_im_sdk

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
	GroupInfo
}
