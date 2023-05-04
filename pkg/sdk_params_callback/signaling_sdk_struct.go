// Copyright © 2023 OpenIM SDK.
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

import "open_im_sdk/pkg/server_api_params"

type InviteCallback *server_api_params.SignalInviteReply

type InviteInGroupCallback *server_api_params.SignalInviteInGroupReply

type CancelCallback *server_api_params.SignalCancelReply

type RejectCallback *server_api_params.SignalRejectReply

type AcceptCallback *server_api_params.SignalAcceptReply

type HungUpCallback *server_api_params.SignalHungUpReply

type GetRoomByGroupIDCallback *server_api_params.SignalGetRoomByGroupIDReply

type GetTokenByRoomID *server_api_params.SignalGetTokenByRoomIDReply
