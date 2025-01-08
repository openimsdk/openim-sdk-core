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

import "github.com/openimsdk/tools/errs"

var (
	// Common errors
	ErrArgs           = errs.NewCodeError(ArgsError, "Invalid input arguments")
	ErrCtxDeadline    = errs.NewCodeError(CtxDeadlineExceededError, "Context deadline exceeded")
	ErrSdkInternal    = errs.NewCodeError(SdkInternalError, "Internal SDK error")
	ErrNetwork        = errs.NewCodeError(NetworkError, "Network error")
	ErrNetworkTimeOut = errs.NewCodeError(NetworkTimeoutError, "Network timeout error")

	ErrGroupIDNotFound = errs.NewCodeError(GroupIDNotFoundError, "Group ID not found")
	ErrUserIDNotFound  = errs.NewCodeError(UserIDNotFoundError, "User ID not found")

	ErrResourceLoad = errs.NewCodeError(ResourceLoadNotCompleteError, "Resource initialization incomplete")

	// Message-related errors
	ErrFileNotFound             = errs.NewCodeError(FileNotFoundError, "File not found")
	ErrMsgDecodeBinaryWs        = errs.NewCodeError(MsgDecodeBinaryWsError, "Message binary WebSocket decoding failed")
	ErrMsgDeCompression         = errs.NewCodeError(MsgDeCompressionError, "Message decompression failed")
	ErrMsgBinaryTypeNotSupport  = errs.NewCodeError(MsgBinaryTypeNotSupportError, "Message type not supported")
	ErrMsgRepeated              = errs.NewCodeError(MsgRepeatError, "Only failed messages can be resent")
	ErrMsgContentTypeNotSupport = errs.NewCodeError(MsgContentTypeNotSupportError, "Message content type not supported")
	ErrMsgHasNoSeq              = errs.NewCodeError(MsgHasNoSeqError, "Message has no sequence number")
	ErrMsgHasDeleted            = errs.NewCodeError(MsgHasDeletedError, "Message has been deleted")

	// Conversation-related errors
	ErrNotSupportOpt  = errs.NewCodeError(NotSupportOptError, "Operation not supported for supergroup")
	ErrNotSupportType = errs.NewCodeError(NotSupportTypeError, "Only supergroup type supported")
	ErrUnreadCount    = errs.NewCodeError(UnreadCountError, "Unread count is zero")

	// Group-related errors
	ErrGroupType = errs.NewCodeError(GroupTypeErr, "Invalid group type")

	ErrLoginOut    = errs.NewCodeError(LoginOutError, "User has logged out")
	ErrLoginRepeat = errs.NewCodeError(LoginRepeatError, "User has logged in repeatedly")
)
