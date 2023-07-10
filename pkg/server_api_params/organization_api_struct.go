package server_api_params

type CreateDepartmentReq struct {
	*Department
	OperationID string `json:"operationID" binding:"required"`
}
type CreateDepartmentResp struct {
	CommResp
	Department *Department            `json:"-"`
	Data       map[string]interface{} `json:"data"`
}

type UpdateDepartmentReq struct {
	*Department
	DepartmentID string `json:"departmentID" binding:"required"`
	OperationID  string `json:"operationID" binding:"required"`
}
type UpdateDepartmentResp struct {
	CommResp
}

type GetSubDepartmentReq struct {
	OperationID  string `json:"operationID" binding:"required"`
	DepartmentID string `json:"departmentID" binding:"required"`
}
type GetSubDepartmentResp struct {
	CommResp
	DepartmentList []*Department            `json:"-"`
	Data           []map[string]interface{} `json:"data"`
}

type DeleteDepartmentReq struct {
	OperationID  string `json:"operationID" binding:"required"`
	DepartmentID string `json:"departmentID" binding:"required"`
}
type DeleteDepartmentResp struct {
	CommResp
}

type CreateOrganizationUserReq struct {
	OperationID string `json:"operationID" binding:"required"`
	*OrganizationUser
}
type CreateOrganizationUserResp struct {
	CommResp
}

type UpdateOrganizationUserReq struct {
	OperationID string `json:"operationID" binding:"required"`
	*OrganizationUser
}
type UpdateOrganizationUserResp struct {
	CommResp
}

type CreateDepartmentMemberReq struct {
	OperationID string `json:"operationID" binding:"required"`
	*DepartmentMember
}

type CreateDepartmentMemberResp struct {
	CommResp
}

type GetUserInDepartmentReq struct {
	UserID      string `json:"userID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}
type GetUserInDepartmentResp struct {
	CommResp
	UserInDepartment *UserInDepartment      `json:"-"`
	Data             map[string]interface{} `json:"data"`
}

type UpdateUserInDepartmentReq struct {
	OperationID string `json:"operationID" binding:"required"`
	*DepartmentMember
}
type UpdateUserInDepartmentResp struct {
	CommResp
}

type DeleteOrganizationUserReq struct {
	UserID      string `json:"userID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}
type DeleteOrganizationUserResp struct {
	CommResp
}

type GetDepartmentMemberReq struct {
	DepartmentID string `json:"departmentID" binding:"required"`
	OperationID  string `json:"operationID" binding:"required"`
}
type GetDepartmentMemberResp struct {
	CommResp
	UserInDepartmentList []*UserDepartmentMember  `json:"-"`
	Data                 []map[string]interface{} `json:"data"`
}

type DeleteUserInDepartmentReq struct {
	DepartmentID string `json:"departmentID" binding:"required"`
	UserID       string `json:"userID" binding:"required"`
	OperationID  string `json:"operationID" binding:"required"`
}
type DeleteUserInDepartmentResp struct {
	CommResp
}
