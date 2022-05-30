package sdk_params_callback

import (
	"open_im_sdk/pkg/db/model_struct"
)

type UserInDepartment struct {
	DepartmentInfo *model_struct.LocalDepartment       `json:"department"`
	MemberInfo     *model_struct.LocalDepartmentMember `json:"member"`
}

type DepartmentAndUser struct {
	model_struct.LocalDepartment
	model_struct.LocalDepartmentMember
}

type GetSubDepartmentCallback []*model_struct.LocalDepartment

type GetDepartmentMemberCallback []*model_struct.LocalDepartmentMember

type GetUserInDepartmentCallback []*UserInDepartment

type GetDepartmentMemberAndSubDepartmentCallback struct {
	DepartmentList       []*model_struct.LocalDepartment       `json:"departmentList"`
	DepartmentMemberList []*model_struct.LocalDepartmentMember `json:"departmentMemberList"`
}
