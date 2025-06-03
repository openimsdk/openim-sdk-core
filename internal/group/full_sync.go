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

package group

import (
	"context"

	"github.com/openimsdk/protocol/sdkws"
)

func (g *Group) GetServerJoinGroup(ctx context.Context) ([]*sdkws.GroupInfo, error) {
	return g.getServerJoinGroup(ctx)
}

func (g *Group) SyncAllJoinedGroupsAndMembers(ctx context.Context) error {
	if err := g.IncrSyncJoinGroup(ctx); err != nil {
		return err
	}
	return g.IncrSyncJoinGroupMember(ctx)
}
