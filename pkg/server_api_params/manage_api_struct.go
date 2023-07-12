// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server_api_params

type DeleteUsersReq struct {
	OperationID      string   `json:"operationID" binding:"required"`
	DeleteUserIDList []string `json:"deleteUserIDList" binding:"required"`
}
type DeleteUsersResp struct {
	CommResp
	FailedUserIDList []string `json:"data"`
}
type GetAllUsersUidReq struct {
	OperationID string `json:"operationID" binding:"required"`
}
type GetAllUsersUidResp struct {
	CommResp
	UserIDList []string `json:"data"`
}
type GetUsersOnlineStatusReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required,lte=200"`
}
type GetUsersOnlineStatusResp struct {
	CommResp
	SuccessResult []GetusersonlinestatusrespSuccessresult `json:"data"`
}
type AccountCheckReq struct {
	OperationID     string   `json:"operationID" binding:"required"`
	CheckUserIDList []string `json:"checkUserIDList" binding:"required,lte=100"`
}
type AccountCheckResp struct {
	CommResp
	Results []*AccountCheckResp_SingleUserStatus `json:"data"`
}
type AccountCheckResp_SingleUserStatus struct {
	UserID        string `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	AccountStatus string `protobuf:"bytes,2,opt,name=accountStatus" json:"accountStatus,omitempty"`
}

type GetusersonlinestatusrespSuccessdetail struct {
	Platform             string   `protobuf:"bytes,1,opt,name=platform" json:"platform,omitempty"`
	Status               string   `protobuf:"bytes,2,opt,name=status" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}
type GetusersonlinestatusrespSuccessresult struct {
	UserID               string                                   `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	Status               string                                   `protobuf:"bytes,2,opt,name=status" json:"status,omitempty"`
	DetailPlatformStatus []*GetusersonlinestatusrespSuccessdetail `protobuf:"bytes,3,rep,name=detailPlatformStatus" json:"detailPlatformStatus,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                 `json:"-"`
	XXX_unrecognized     []byte                                   `json:"-"`
	XXX_sizecache        int32                                    `json:"-"`
}
type AccountcheckrespSingleuserstatus struct {
	UserID               string   `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	AccountStatus        string   `protobuf:"bytes,2,opt,name=accountStatus" json:"accountStatus,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}
