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
type WrapperGroup struct {
	*WrapperCommon
}

func NewWrapperGroup(wrapperCommon *WrapperCommon) *WrapperGroup {
	return &WrapperGroup{WrapperCommon: wrapperCommon}
}

func (w *WrapperGroup) CreateGroup(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.CreateGroup, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) JoinGroup(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.JoinGroup, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) QuitGroup(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.QuitGroup, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) DismissGroup(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.DismissGroup, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) ChangeGroupMute(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.ChangeGroupMute, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) ChangeGroupMemberMute(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.ChangeGroupMemberMute, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) SetGroupMemberRoleLevel(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SetGroupMemberRoleLevel, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetJoinedGroupList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetJoinedGroupList, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) SearchGroups(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SearchGroups, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) SetGroupInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SetGroupInfo, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) SetGroupVerification(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SetGroupVerification, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) SetGroupLookMemberInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SetGroupLookMemberInfo, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) SetGroupApplyMemberFriend(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SetGroupApplyMemberFriend, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetGroupMemberList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetGroupMemberList, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetGroupMemberOwnerAndAdmin(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetGroupMemberOwnerAndAdmin, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetGroupMemberListByJoinTimeFilter(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetGroupMemberListByJoinTimeFilter, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetSpecifiedGroupMembersInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetSpecifiedGroupMembersInfo, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) KickGroupMember(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.KickGroupMember, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) TransferGroupOwner(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.TransferGroupOwner, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) InviteUserToGroup(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.InviteUserToGroup, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetGroupApplicationListAsRecipient(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetGroupApplicationListAsRecipient, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetGroupApplicationListAsApplicant(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetGroupApplicationListAsApplicant, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) AcceptGroupApplication(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.AcceptGroupApplication, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) RefuseGroupApplication(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.RefuseGroupApplication, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) SetGroupMemberNickname(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SetGroupMemberNickname, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) SearchGroupMembers(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.SearchGroupMembers, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) IsJoinGroup(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.IsJoinGroup, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) GetSpecifiedGroupsInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetSpecifiedGroupsInfo, callback, &args).AsyncCallWithCallback()
}
