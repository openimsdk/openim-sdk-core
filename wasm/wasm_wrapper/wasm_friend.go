package wasm_wrapper

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/event_listener"
	"syscall/js"
)

//------------------------------------group---------------------------
type WrapperFriend struct {
	*WrapperCommon
}

func NewWrapperFriend(wrapperCommon *WrapperCommon) *WrapperFriend {
	return &WrapperFriend{WrapperCommon: wrapperCommon}
}

func (w *WrapperFriend) GetDesignatedFriendsInfo(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetDesignatedFriendsInfo, callback, &args).AsyncCallWithCallback()
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

func (w *WrapperFriend) GetRecvFriendApplicationList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetRecvFriendApplicationList, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperFriend) GetSendFriendApplicationList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetSendFriendApplicationList, callback, &args).AsyncCallWithCallback()
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
