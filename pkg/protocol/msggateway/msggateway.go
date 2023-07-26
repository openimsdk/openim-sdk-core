// Copyright © 2023 OpenIM. All rights reserved.
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

package msggateway

import "errors"

func (x *OnlinePushMsgReq) Check() error {
	if x.MsgData == nil {
		return errors.New("MsgData is empty")
	}
	if err := x.MsgData.Check(); err != nil {
		return err
	}
	if x.PushToUserID == "" {
		return errors.New("PushToUserID is empty")
	}
	return nil
}

func (x *OnlineBatchPushOneMsgReq) Check() error {
	if x.MsgData == nil {
		return errors.New("MsgData is empty")
	}
	if err := x.MsgData.Check(); err != nil {
		return err
	}
	if x.PushToUserIDs == nil {
		return errors.New("PushToUserIDs is empty")
	}
	return nil
}

func (x *GetUsersOnlineStatusReq) Check() error {
	if x.UserIDs == nil {
		return errors.New("UserIDs is empty")
	}
	return nil
}

func (x *KickUserOfflineReq) Check() error {
	if x.PlatformID < 1 || x.PlatformID > 9 {
		return errors.New("PlatformID is invalid")
	}
	if x.KickUserIDList == nil {
		return errors.New("KickUserIDList is empty")
	}
	return nil
}

func (x *MultiTerminalLoginCheckReq) Check() error {
	if x.PlatformID < 1 || x.PlatformID > 9 {
		return errors.New("PlatformID is invalid")
	}
	if x.UserID == "" {
		return errors.New("UserID is empty")
	}
	if x.Token == "" {
		return errors.New("token is empty")
	}
	return nil
}
