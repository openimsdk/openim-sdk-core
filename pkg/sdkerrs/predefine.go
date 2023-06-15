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

package sdkerrs

import "github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

var (
	ErrArgs           = errs.NewCodeError(ArgsError, "ArgsError")
	ErrCtxDeadline    = errs.NewCodeError(CtxDeadlineExceededError, "CtxDeadlineExceededError")
	ErrSdkInternal    = errs.NewCodeError(SdkInternalError, "SdkInternalError")
	ErrNetwork        = errs.NewCodeError(NetworkError, "NetworkError")
	ErrNetworkTimeOut = errs.NewCodeError(NetworkTimeoutError, "NetworkTimeoutError")
	ErrDuplicateKey   = errs.NewCodeError(DuplicateKeyError, "DuplicateKeyError")

	ErrGroupIDNotFound = errs.NewCodeError(GroupIDNotFoundError, "GroupIDNotFoundError")
	ErrUserIDNotFound  = errs.NewCodeError(UserIDNotFoundError, "UserIDNotFoundError")

	ErrRecordNotFound = errs.NewCodeError(RecordNotFoundError, "RecordNotFoundError")

	//消息相关
	ErrMsgDecodeBinaryWs        = errs.NewCodeError(MsgDecodeBinaryWsError, "MsgDecodeBinaryWsError")
	ErrMsgDeCompression         = errs.NewCodeError(MsgDeCompressionError, "MsgDeCompressionError")
	ErrMsgTypeNotSupport        = errs.NewCodeError(MsgTypeNotSupportError, "MsgTypeNotSupportError")
	ErrMsgRepeated              = errs.NewCodeError(MsgRepeatError, "only failed message can be repeatedly send")
	ErrMsgContentTypeNotSupport = errs.NewCodeError(MsgContentTypeNotSupportError, "contentType not support currently") // msg
	ErrMsgNotFound              = errs.NewCodeError(MsgNotFoundError, "MsgNotFoundError")                               // msg

	//会话相关
	ErrNotSupportOpt        = errs.NewCodeError(NotSupportOptError, "super group not support this opt")
	ErrNotResetGroupAtType  = errs.NewCodeError(NotResetGroupAtTypeError, "conversation don't need to reset")
	ErrNotFoundConversation = errs.NewCodeError(NotFoundConversation, "conversation not found")
	//群组相关

	ErrNotInGroup = errs.NewCodeError(NotInGroupError, "you not exist in this group")
	ErrGroupType  = errs.NewCodeError(GroupTypeErr, "group type error")

	ErrLoginOut = errs.NewCodeError(LoginOutError, "MsgLoginOutError")
)
