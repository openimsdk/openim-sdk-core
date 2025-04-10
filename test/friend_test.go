package test

//
//import (
//	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
//	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
//	"github.com/openimsdk/protocol/wrapperspb"
//	"testing"
//	"time"
//
//	friend2 "github.com/openimsdk/protocol/relation"
//)
//
//func Test_GetSpecifiedFriendsInfo(t *testing.T) {
//	info, err := open_im_sdk.UserForSDK.Friend().GetSpecifiedFriendsInfo(ctx, []string{"45644221123"})
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("GetDesignatedFriendsInfo success", ctx.Value("operationID"))
//	for _, userInfo := range info {
//		t.Log(userInfo)
//	}
//}
//
//func Test_AddFriend(t *testing.T) {
//	err := open_im_sdk.UserForSDK.Friend().AddFriend(ctx, &friend2.ApplyToAddFriendReq{
//		ToUserID: "2",
//		ReqMsg:   "test add",
//		Ex:       "add",
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("AddFriend success", ctx.Value("operationID"))
//}
//
//
//func Test_AcceptFriendApplication(t *testing.T) {
//	req := &sdk_params_callback.ProcessFriendApplicationParams{ToUserID: "1", HandleMsg: "test accept"}
//	err := open_im_sdk.UserForSDK.Friend().AcceptFriendApplication(ctx, req)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("AcceptFriendApplication success", ctx.Value("operationID"))
//	time.Sleep(time.Second * 30)
//}
//
//func Test_RefuseFriendApplication(t *testing.T) {
//	req := &sdk_params_callback.ProcessFriendApplicationParams{ToUserID: "6754269405", HandleMsg: "test refuse"}
//	err := open_im_sdk.UserForSDK.Friend().RefuseFriendApplication(ctx, req)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("RefuseFriendApplication success", ctx.Value("operationID"))
//	time.Sleep(time.Second * 30)
//}
//
//func Test_CheckFriend(t *testing.T) {
//	res, err := open_im_sdk.UserForSDK.Friend().CheckFriend(ctx, []string{"863454357", "45644221123"})
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("CheckFriend success", ctx.Value("operationID"))
//	for _, re := range res {
//		t.Log(re)
//	}
//}
//func Test_PinFriend(t *testing.T) {
//	pinParams := &sdk_params_callback.SetFriendPinParams{
//		ToUserIDs: []string{"2", "3"},
//		IsPinned:  &wrapperspb.BoolValue{Value: false},
//	}
//
//	err := open_im_sdk.UserForSDK.Friend().PinFriends(ctx, pinParams)
//
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("CheckFriend success", ctx.Value("operationID"))
//}
//func Test_DeleteFriend(t *testing.T) {
//	err := open_im_sdk.UserForSDK.Friend().DeleteFriend(ctx, "863454357")
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("DeleteFriend success", ctx.Value("operationID"))
//}
//
//func Test_GetFriendList(t *testing.T) {
//	infos, err := open_im_sdk.UserForSDK.Friend().GetFriendList(ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("GetFriendList success", ctx.Value("operationID"))
//	for _, info := range infos {
//		t.Logf("PublicInfo: %#v, FriendInfo: %#v, BlackInfo: %#v", info.PublicInfo, info.FriendInfo, info.BlackInfo)
//	}
//}
//
//func Test_SearchFriends(t *testing.T) {
//	info, err := open_im_sdk.UserForSDK.Friend().SearchFriends(ctx, &sdk_params_callback.SearchFriendsParam{KeywordList: []string{"863454357"}, IsSearchUserID: true})
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("SearchFriends success", ctx.Value("operationID"))
//	for _, item := range info {
//		t.Log(*item)
//	}
//}
//
//func Test_SetFriendRemark(t *testing.T) {
//	err := open_im_sdk.UserForSDK.Friend().SetFriendRemark(ctx, &sdk_params_callback.SetFriendRemarkParams{ToUserID: "863454357", Remark: "testRemark"})
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("SetFriendRemark success", ctx.Value("operationID"))
//}
//
//func Test_AddBlack(t *testing.T) {
//	err := open_im_sdk.UserForSDK.Friend().AddBlack(ctx, "863454357", "ex")
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("AddBlack success", ctx.Value("operationID"))
//}
//
//func Test_RemoveBlack(t *testing.T) {
//	err := open_im_sdk.UserForSDK.Friend().RemoveBlack(ctx, "863454357")
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("RemoveBlack success", ctx.Value("operationID"))
//}
//
//func Test_GetBlackList(t *testing.T) {
//	info, err := open_im_sdk.UserForSDK.Friend().GetBlackList(ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("GetBlackList success", ctx.Value("operationID"))
//	for _, item := range info {
//		t.Log(*item)
//	}
//}
//func Test_SetFriendsEx(t *testing.T) {
//	err := open_im_sdk.UserForSDK.Friend().SetFriendsEx(ctx, []string{"2"}, "exx")
//	if err != nil {
//		t.Fatal(err)
//	}
//}
