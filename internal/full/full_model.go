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

package full

import (
	"context"
	"open_im_sdk/pkg/db/model_struct"
)

func (u *Full) GetGroupInfoByGroupID(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
	g2, err := u.group.GetGroupInfoFromLocal2Svr(ctx, groupID)
	return g2, err
}

func (u *Full) GetGroupsInfo(ctx context.Context, groupIDs ...string) (map[string]*model_struct.LocalGroup, error) {
	return u.group.GetGroupsInfoFromLocal2Svr(ctx, groupIDs...)
}
