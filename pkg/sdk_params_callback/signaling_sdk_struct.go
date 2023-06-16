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

import "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"

type InviteCallback *sdkws.SignalInviteReply

type InviteInGroupCallback *sdkws.SignalInviteInGroupReply

type CancelCallback *sdkws.SignalCancelReply

type RejectCallback *sdkws.SignalRejectReply

type AcceptCallback *sdkws.SignalAcceptReply

type HungUpCallback *sdkws.SignalHungUpReply

type GetRoomByGroupIDCallback *sdkws.SignalGetRoomByGroupIDReply

type GetTokenByRoomID *sdkws.SignalGetTokenByRoomIDReply
