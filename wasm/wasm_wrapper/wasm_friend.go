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

//go:build js && wasm
// +build js,wasm

package wasm_wrapper

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/event_listener"
	"syscall/js"
)

// ------------------------------------group---------------------------
type WrapperFriend struct {
	*WrapperCommon
}

func NewWrapperFriend(wrapperCommon *WrapperCommon) *WrapperFriend {
	return &WrapperFriend{WrapperCommon: wrapperCommon}
}

func (w *WrapperFriend) GetSpecifiedFriendsInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetSpecifiedFriendsInfo, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) GetFriendList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetFriendList, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) SearchFriends(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SearchFriends, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) CheckFriend(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.CheckFriend, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) AddFriend(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.AddFriend, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) SetFriendRemark(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SetFriendRemark, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) DeleteFriend(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.DeleteFriend, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) GetFriendApplicationListAsRecipient(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetFriendApplicationListAsRecipient, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) GetFriendApplicationListAsApplicant(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetFriendApplicationListAsApplicant, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) AcceptFriendApplication(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.AcceptFriendApplication, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) RefuseFriendApplication(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.RefuseFriendApplication, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) GetBlackList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetBlackList, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) RemoveBlack(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.RemoveBlack, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) AddBlack(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.AddBlack, callback, &args).AsyncCallWithCallback()
}
