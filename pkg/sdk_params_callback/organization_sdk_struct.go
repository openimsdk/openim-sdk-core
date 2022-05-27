package sdk_params_callback

import "open_im_sdk/pkg/db"

type UserInDepartment struct {
	DepartmentInfo *db.LocalDepartment       `json:"department"`
	MemberInfo     *db.LocalDepartmentMember `json:"member"`
}

type DepartmentAndUser struct {
	db.LocalDepartment
	db.LocalDepartmentMember
}

type GetSubDepartmentCallback []*db.LocalDepartment

type GetDepartmentMemberCallback []*db.LocalDepartmentMember

type GetUserInDepartmentCallback []*UserInDepartment

type GetDepartmentMemberAndSubDepartmentCallback struct {
	DepartmentList       []*db.LocalDepartment       `json:"departmentList"`
	DepartmentMemberList []*db.LocalDepartmentMember `json:"departmentMemberList"`
}
