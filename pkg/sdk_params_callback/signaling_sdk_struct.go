// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

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
