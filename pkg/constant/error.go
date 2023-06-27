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

package constant

// key = errCode, string = errMsg
//type ErrInfo struct {
//	ErrCode int32
//	ErrMsg  string
//}
//
//var (
//	OK = ErrInfo{0, ""}
//
//	ErrParseToken = ErrInfo{200, ParseTokenMsg.Error()}
//
//	ErrTencentCredential = ErrInfo{400, ThirdPartyMsg.Error()}
//	ErrInBlackList       = ErrInfo{ErrCode: 600, ErrMsg: InBlackList.Error()}
//	ErrNotFriend         = ErrInfo{ErrCode: 601, ErrMsg: NotFriend.Error()}
//
//	ErrTokenExpired     = ErrInfo{701, TokenExpiredMsg.Error()}
//	ErrTokenInvalid     = ErrInfo{702, TokenInvalidMsg.Error()}
//	ErrTokenMalformed   = ErrInfo{703, TokenMalformedMsg.Error()}
//	ErrTokenNotValidYet = ErrInfo{704, TokenNotValidYetMsg.Error()}
//	ErrTokenUnknown     = ErrInfo{705, TokenUnknownMsg.Error()}
//	ErrTokenKicked      = ErrInfo{706, TokenUserKickedMsg.Error()}
//
//	ErrAccess       = ErrInfo{ErrCode: 801, ErrMsg: AccessMsg.Error()}
//	ErrDB           = ErrInfo{ErrCode: 802, ErrMsg: DBMsg.Error()}
//	ErrArgs         = ErrInfo{ErrCode: 803, ErrMsg: ArgsMsg.Error()}
//	ErrApi          = ErrInfo{ErrCode: 804, ErrMsg: ApiMsg.Error()}
//	ErrData         = ErrInfo{ErrCode: 805, ErrMsg: DataMsg.Error()}
//	ErrLogin        = ErrInfo{ErrCode: 806, ErrMsg: LoginMsg.Error()}
//	ErrConfig       = ErrInfo{ErrCode: 807, ErrMsg: ConfigMsg.Error()}
//	ErrThirdParty   = ErrInfo{ErrCode: 808, ErrMsg: ThirdPartyMsg.Error()}
//	ErrServerReturn = ErrInfo{ErrCode: 809, ErrMsg: ServerReturn.Error()}
//
//	ErrWsRecvConnDiff          = ErrInfo{ErrCode: 901, ErrMsg: WsRecvConnDiff.Error()}
//	ErrWsRecvConnSame          = ErrInfo{ErrCode: 902, ErrMsg: WsRecvConnSame.Error()}
//	ErrWsRecvCode              = ErrInfo{ErrCode: 903, ErrMsg: WsRecvCode.Error()}
//	ErrWsSendTimeout           = ErrInfo{ErrCode: 904, ErrMsg: WsSendTimeout.Error()}
//	ErrResourceLoadNotComplete = ErrInfo{ErrCode: 905, ErrMsg: ResourceLoadNotComplete.Error()}
//	ErrNotSupportFunction      = ErrInfo{ErrCode: 906, ErrMsg: NotSupportFunction.Error()}
//)
//
//var (
//	ParseTokenMsg       = errors.New("parse token failed")
//	TokenExpiredMsg     = errors.New("token is timed out, please log in again")
//	TokenInvalidMsg     = errors.New("token has been invalidated")
//	TokenNotValidYetMsg = errors.New("token not active yet")
//	TokenMalformedMsg   = errors.New("that's not even a token")
//	TokenUnknownMsg     = errors.New("couldn't handle this token")
//	TokenUserKickedMsg  = errors.New("user has been kicked")
//
//	AccessMsg = errors.New("no permission")
//	DBMsg     = errors.New("db failed")
//	ArgsMsg   = errors.New("args failed")
//	ApiMsg    = errors.New("api failed")
//	DataMsg   = errors.New("data failed ")
//	LoginMsg  = errors.New("you can only login once")
//	ConfigMsg = errors.New("config failed")
//
//	ThirdPartyMsg = errors.New("third party error")
//	ServerReturn  = errors.New("server return data err")
//
//	WsRecvConnDiff          = errors.New("recv timeout, conn diff")
//	WsRecvConnSame          = errors.New("recv timeout, conn same")
//	WsRecvCode              = errors.New("recv code err")
//	WsSendTimeout           = errors.New("send timeout")
//	ResourceLoadNotComplete = errors.New("resource loading is not complete")
//	NotSupportFunction      = errors.New("unsupported function")
//
//	NotFriend   = errors.New("not friend")
//	InBlackList = errors.New("in blackList")
//)
//
//funcation (e *ErrInfo) Error() string {
//	return e.ErrMsg
//}
//
//const (
//	StatusErrTokenExpired     = 701
//	StatusErrTokenInvalid     = 702
//	StatusErrTokenMalformed   = 703
//	StatusErrTokenNotValidYet = 704
//	StatusErrTokenUnknown     = 705
//	StatusErrTokenKicked      = 706
//)
//
//var statusText = map[int]*ErrInfo{
//	StatusErrTokenExpired:     &ErrTokenExpired,
//	StatusErrTokenInvalid:     &ErrTokenInvalid,
//	StatusErrTokenMalformed:   &ErrTokenMalformed,
//	StatusErrTokenNotValidYet: &ErrTokenNotValidYet,
//	StatusErrTokenUnknown:     &ErrTokenUnknown,
//	StatusErrTokenKicked:      &ErrTokenKicked,
//}
//
//funcation StatusText(code int) *ErrInfo {
//	return statusText[code]
//}
