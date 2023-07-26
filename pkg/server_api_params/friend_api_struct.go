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

import "github.com/OpenIMSDK/protocol/sdkws"

type ParamsCommFriend struct {
	OperationID string `json:"operationID" binding:"required"`
	ToUserID    string `json:"toUserID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}

type AddBlacklistReq struct {
	ParamsCommFriend
}
type AddBlacklistResp struct {
	CommResp
}

type ImportFriendReq struct {
	FriendUserIDList []string `json:"friendUserIDList" binding:"required"`
	OperationID      string   `json:"operationID" binding:"required"`
	FromUserID       string   `json:"fromUserID" binding:"required"`
}
type UserIDResult struct {
	UserID string `json:"userID""`
	Result int32  `json:"result"`
}
type ImportFriendResp struct {
	CommResp
	UserIDResultList []UserIDResult `json:"data"`
}

type AddFriendReq struct {
	ParamsCommFriend
	ReqMsg string `json:"reqMsg"`
}
type AddFriendResp struct {
	CommResp
}

type AddFriendResponseReq struct {
	ParamsCommFriend
	Flag      int32  `json:"flag" binding:"required,oneof=-1 0 1"`
	HandleMsg string `json:"handleMsg"`
}
type AddFriendResponseResp struct {
	CommResp
}

type DeleteFriendReq struct {
	ParamsCommFriend
}
type DeleteFriendResp struct {
	CommResp
}

type GetBlackListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetBlackListResp struct {
	CommResp
	BlackUserInfoList []*sdkws.PublicUserInfo
	Data              []map[string]interface{} `json:"data"`
}

//type PublicUserInfo struct {
//	UserID   string `json:"userID"`
//	Nickname string `json:"nickname"`
//	FaceUrl  string `json:"faceUrl"`
//	Gender   int32  `json:"gender"`
//}

type SetFriendRemarkReq struct {
	ParamsCommFriend
	Remark string `json:"remark" binding:"required"`
}
type SetFriendRemarkResp struct {
	CommResp
}

type RemoveBlackListReq struct {
	ParamsCommFriend
}
type RemoveBlackListResp struct {
	CommResp
}

type IsFriendReq struct {
	ParamsCommFriend
}
type Response struct {
	Friend bool `json:"isFriend"`
}
type IsFriendResp struct {
	CommResp
	Response Response `json:"data"`
}

type GetFriendsInfoReq struct {
	ParamsCommFriend
}
type GetFriendsInfoResp struct {
	CommResp
	FriendInfoList []*sdkws.FriendInfo
	Data           []map[string]interface{} `json:"data"`
}

type GetFriendListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetFriendListResp struct {
	CommResp
	FriendInfoList []*sdkws.FriendInfo
	Data           []map[string]interface{} `json:"data"`
}

type GetFriendApplyListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetFriendApplyListResp struct {
	CommResp
	FriendRequestList []*sdkws.FriendRequest
	Data              []map[string]interface{} `json:"data"`
}

type GetSelfFriendApplyListReq struct {
	OperationID string `json:"operationID" binding:"required"`
	FromUserID  string `json:"fromUserID" binding:"required"`
}
type GetSelfFriendApplyListResp struct {
	CommResp
	FriendRequestList []*sdkws.FriendRequest
	Data              []map[string]interface{} `json:"data"`
}
