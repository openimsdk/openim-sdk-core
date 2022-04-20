package sdk_params_callback

import "open_im_sdk/pkg/db"

type UserInDepartment struct {
	DepartmentInfo *db.LocalDepartment
	MemberInfo     *db.LocalDepartmentMember
}

type DepartmentAndUser struct {
	db.LocalDepartment
	db.LocalDepartmentMember
}

type GetSubDepartmentCallback []*db.LocalDepartment

type GetDepartmentMemberCallback []*db.LocalDepartmentMember

type GetUserInDepartmentCallback []*UserInDepartment

type GetDepartmentMemberAndSubDepartmentCallback struct {
	DepartmentList       []*db.LocalDepartment
	DepartmentMemberList []*db.LocalDepartmentMember
}
