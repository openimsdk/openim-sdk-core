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

// Common error codes
const (
	NetworkError             = 10000
	NetworkTimeoutError      = 10001
	ArgsError                = 10002 // Invalid input parameters
	CtxDeadlineExceededError = 10003 // Context deadline exceeded

	ResourceLoadNotCompleteError = 10004 // Resource initialization incomplete
	UnknownCode                  = 10005 // Unrecognized code
	SdkInternalError             = 10006 // SDK internal error

	NoUpdateError = 10007 // No updates available

	UserIDNotFoundError = 10100 // UserID not found or not registered
	LoginOutError       = 10101 // User has logged out
	LoginRepeatError    = 10102 // User logged in repeatedly

	// Message-related errors
	FileNotFoundError             = 10200 // Record not found
	MsgDeCompressionError         = 10201 // Message decompression failed
	MsgDecodeBinaryWsError        = 10202 // Message decoding failed
	MsgBinaryTypeNotSupportError  = 10203 // Message type not supported
	MsgRepeatError                = 10204 // Message repeated
	MsgContentTypeNotSupportError = 10205 // Message content type not supported
	MsgHasNoSeqError              = 10206 // Message does not have a sequence number
	MsgHasDeletedError            = 10207 // Message has been deleted

	// Conversation-related errors
	NotSupportOptError  = 10301 // Operation not supported
	NotSupportTypeError = 10302 // Type not supported
	UnreadCountError    = 10303 // Unread count is zero

	// Group-related errors
	GroupIDNotFoundError = 10400 // GroupID not found
	GroupTypeErr         = 10401 // Invalid group type
)
