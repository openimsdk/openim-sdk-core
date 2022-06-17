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

type GetParentDepartmentListCallback []*model_struct.LocalDepartment

type GetSubDepartmentCallback []*model_struct.LocalDepartment

type GetDepartmentMemberCallback []*model_struct.LocalDepartmentMember

type GetUserInDepartmentCallback []*UserInDepartment

type ParentDepartmentCallback struct {
	Name         string `json:"name"`
	DepartmentID string `json:"departmentID"`
}

type GetDepartmentMemberAndSubDepartmentCallback struct {
	DepartmentList       []*model_struct.LocalDepartment       `json:"departmentList"`
	DepartmentMemberList []*model_struct.LocalDepartmentMember `json:"departmentMemberList"`
	ParentDepartmentList []ParentDepartmentCallback            `json:"parentDepartmentList"`
}

type GetDepartmentInfoCallback *model_struct.LocalDepartment

type SearchOrganizationParams struct {
	KeyWord                 string `json:"keyWord"`
	IsSearchUserName        bool   `json:"isSearchUserName"`
	IsSearchUserEnglishName bool   `json:"isSearchEnglishName"`
	IsSearchPosition        bool   `json:"isSearchPosition"`
	IsSearchUserID          bool   `json:"isSearchUserID"`
	IsSearchMobile          bool   `json:"isSearchMobile"`
	IsSearchEmail           bool   `json:"isSearchEmail"`
	IsSearchTelephone       bool   `json:"isSearchTelephone"`
}

type SearchOrganizationCallback struct {
	DepartmentList       []*model_struct.LocalDepartment `json:"departmentList"`
	DepartmentMemberList []*struct {
		*model_struct.SearchDepartmentMemberResult
		ParentDepartmentList []*ParentDepartmentCallback `json:"parentDepartmentList"`
	} `json:"departmentMemberList"`
}
