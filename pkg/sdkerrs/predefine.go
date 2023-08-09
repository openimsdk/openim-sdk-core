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

import "github.com/OpenIMSDK/tools/errs"

var (
	ErrArgs           = errs.NewCodeError(ArgsError, "ArgsError")
	ErrCtxDeadline    = errs.NewCodeError(CtxDeadlineExceededError, "CtxDeadlineExceededError")
	ErrSdkInternal    = errs.NewCodeError(SdkInternalError, "SdkInternalError")
	ErrNetwork        = errs.NewCodeError(NetworkError, "NetworkError")
	ErrNetworkTimeOut = errs.NewCodeError(NetworkTimeoutError, "NetworkTimeoutError")

	ErrGroupIDNotFound = errs.NewCodeError(GroupIDNotFoundError, "GroupIDNotFoundError")
	ErrUserIDNotFound  = errs.NewCodeError(UserIDNotFoundError, "UserIDNotFoundError")

	ErrResourceLoad = errs.NewCodeError(ResourceLoadNotCompleteError, "ResourceLoadNotCompleteError")

	//消息相关
	ErrFileNotFound             = errs.NewCodeError(FileNotFoundError, "RecordNotFoundError")
	ErrMsgDecodeBinaryWs        = errs.NewCodeError(MsgDecodeBinaryWsError, "MsgDecodeBinaryWsError")
	ErrMsgDeCompression         = errs.NewCodeError(MsgDeCompressionError, "MsgDeCompressionError")
	ErrMsgBinaryTypeNotSupport  = errs.NewCodeError(MsgBinaryTypeNotSupportError, "MsgTypeNotSupportError")
	ErrMsgRepeated              = errs.NewCodeError(MsgRepeatError, "only failed message can be repeatedly send")
	ErrMsgContentTypeNotSupport = errs.NewCodeError(MsgContentTypeNotSupportError, "contentType not support currently") // msg 	// msg
	ErrMsgHasNoSeq              = errs.NewCodeError(MsgHasNoSeqError, "msg has no seq")                                 // msg 	// msg

	//会话相关
	ErrNotSupportOpt = errs.NewCodeError(NotSupportOptError, "super group not support this opt")
	//群组相关

	ErrGroupType = errs.NewCodeError(GroupTypeErr, "group type error")

	ErrLoginOut    = errs.NewCodeError(LoginOutError, "LoginOutError")
	ErrLoginRepeat = errs.NewCodeError(LoginRepeatError, "LoginRepeatError")
)
