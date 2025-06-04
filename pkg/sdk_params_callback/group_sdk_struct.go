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

type SearchGroupsParam struct {
	KeywordList       []string `json:"keywordList"`
	IsSearchGroupID   bool     `json:"isSearchGroupID"`
	IsSearchGroupName bool     `json:"isSearchGroupName"`
}

type SearchGroupMembersParam struct {
	GroupID                string   `json:"groupID"`
	KeywordList            []string `json:"keywordList"`
	IsSearchUserID         bool     `json:"isSearchUserID"`
	IsSearchMemberNickname bool     `json:"isSearchMemberNickname"`
	Offset                 int      `json:"offset"`
	Count                  int      `json:"count"`
	PageNumber             int      `json:"pageNumber"`
}

type GetGroupApplicationListAsRecipientReq struct {
	GroupIDs      []string `json:"groupIDs"`
	HandleResults []int32  `json:"handleResults"`
	Offset        int32    `json:"offset"`
	Count         int32    `json:"count"`
}

type GetGroupApplicationListAsApplicantReq struct {
	GroupIDs      []string `json:"groupIDs"`
	HandleResults []int32  `json:"handleResults"`
	Offset        int32    `json:"offset"`
	Count         int32    `json:"count"`
}

type GetGroupApplicationUnhandledCountReq struct {
	Time int64 `json:"time"`
}
