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

type GetParentDepartmentListCallback []*db.LocalDepartment

type GetDepartmentMemberCallback []*db.LocalDepartmentMember

type GetUserInDepartmentCallback []*UserInDepartment

type ParentDepartmentCallback struct {
	Name         string `json:"name"`
	DepartmentID string `json:"departmentID"`
}

type GetDepartmentMemberAndSubDepartmentCallback struct {
	DepartmentList       []*db.LocalDepartment       `json:"departmentList"`
	DepartmentMemberList []*db.LocalDepartmentMember `json:"departmentMemberList"`
	ParentDepartmentList []ParentDepartmentCallback  `json:"parentDepartmentList"`
}

type GetDepartmentInfoCallback *db.LocalDepartment

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
	DepartmentList       []*db.LocalDepartment `json:"departmentList"`
	DepartmentMemberList []*struct {
		*db.SearchDepartmentMemberResult
		ParentDepartmentList []*ParentDepartmentCallback `json:"parentDepartmentList"`
	} `json:"departmentMemberList"`
}
