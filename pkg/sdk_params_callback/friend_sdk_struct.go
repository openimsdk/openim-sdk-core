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

package sdk_params_callback

import (
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/protocol/wrapperspb"
)

type ProcessFriendApplicationParams struct {
	ToUserID  string `json:"toUserID" validate:"required"`
	HandleMsg string `json:"handleMsg"`
}

type SearchFriendsParam struct {
	KeywordList      []string `json:"keywordList"`
	IsSearchUserID   bool     `json:"isSearchUserID"`
	IsSearchNickname bool     `json:"isSearchNickname"`
	IsSearchRemark   bool     `json:"isSearchRemark"`
}

type SearchFriendItem struct {
	model_struct.LocalFriend
	Relationship int `json:"relationship"`
}

type SetFriendRemarkParams struct {
	ToUserID string `json:"toUserID" validate:"required"`
	Remark   string `json:"remark" validate:"required"`
}
type SetFriendPinParams struct {
	ToUserIDs []string              `json:"toUserIDs" validate:"required"`
	IsPinned  *wrapperspb.BoolValue `json:"isPinned" validate:"required"`
}

type GetFriendApplicationListAsRecipientReq struct {
	HandleResults []int32 `json:"handleResults"`
	Offset        int32   `json:"offset"`
	Count         int32   `json:"count"`
}

type GetFriendApplicationListAsApplicantReq struct {
	Offset int32 `json:"offset"`
	Count  int32 `json:"count"`
}

type GetSelfUnhandledApplyCountReq struct {
	Time int64 `json:"time"`
}
