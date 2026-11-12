package interaction

import (
	"context"
	"errors"
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/utils/datautil"
)

func TestName(t *testing.T) {
	sub := newSubscription()

	//sub.setUserState([]*sdkws.SubUserOnlineStatusElem{
	//	{
	//		UserID:            "1",
	//		OnlinePlatformIDs: []int32{1, 2, 3},
	//	},
	//	{
	//		UserID:            "2",
	//		OnlinePlatformIDs: []int32{1},
	//	},
	//})

	exist, wait, subUserIDs, unsubUserIDs := sub.getUserOnline([]string{"1", "2", "3"})

	t.Logf("exist: %v", exist)
	t.Logf("wait: %v", wait)
	t.Logf("subUserIDs: %v", subUserIDs)
	t.Logf("unsubUserIDs: %v", unsubUserIDs)

	sub.writeFailed(wait, errors.New("todo test"))

}

func TestSubscriptionRemembersUsersForStatusQuery(t *testing.T) {
	sub := newSubscription()
	_, _, subscribed, _ := sub.getUserOnline([]string{"user-1", "user-2"})
	if !datautil.Contain("user-1", subscribed...) || !datautil.Contain("user-2", subscribed...) {
		t.Fatalf("initial subscription request = %v", subscribed)
	}

	sub.setUserState([]*sdkws.SubUserOnlineStatusElem{
		{UserID: "user-1", OnlinePlatformIDs: []int32{5}},
		{UserID: "user-2", OnlinePlatformIDs: []int32{}},
	})
	existing, wait, subscribeAgain, unsubscribe := sub.getUserOnline(nil)
	if len(wait) != 0 || len(subscribeAgain) != 0 || len(unsubscribe) != 0 {
		t.Fatalf(
			"cached status query unexpectedly changed subscriptions: wait=%v subscribe=%v unsubscribe=%v",
			wait,
			subscribeAgain,
			unsubscribe,
		)
	}
	if !equal(existing["user-1"], []int32{5}) {
		t.Fatalf("user-1 status = %v", existing["user-1"])
	}
	if status, ok := existing["user-2"]; !ok || len(status) != 0 {
		t.Fatalf("user-2 status = %v, present=%v", status, ok)
	}
}

func TestGetSubscribeUsersStatusAndReconnectPush(t *testing.T) {
	ctx := context.Background()
	mgr := &LongConnMgr{sub: newSubscription()}

	emptyStatuses, err := mgr.GetSubscribeUsersStatus(ctx)
	if err != nil || len(emptyStatuses) != 0 {
		t.Fatalf("initial GetSubscribeUsersStatus = %v, err = %v", emptyStatuses, err)
	}

	mgr.sub.getUserOnline([]string{"user-1", "user-2"})
	mgr.sub.setUserState([]*sdkws.SubUserOnlineStatusElem{
		{UserID: "user-1", OnlinePlatformIDs: []int32{5}},
		{UserID: "user-2", OnlinePlatformIDs: []int32{}},
	})

	statuses, err := mgr.GetSubscribeUsersStatus(ctx)
	if err != nil || len(statuses) != 2 {
		t.Fatalf("GetSubscribeUsersStatus = %v, err = %v", statuses, err)
	}

	// Simulate connection drop and reconnect re-subscription + server push.
	mgr.sub.onConnClosed(nil)
	reconnUsers := mgr.sub.getNewConnSubUserIDs()
	if !datautil.Contain("user-1", reconnUsers...) || !datautil.Contain("user-2", reconnUsers...) {
		t.Fatalf("getNewConnSubUserIDs = %v", reconnUsers)
	}

	changes := mgr.sub.setUserState([]*sdkws.SubUserOnlineStatusElem{
		{UserID: "user-1", OnlinePlatformIDs: []int32{1, 5}},
		{UserID: "user-2", OnlinePlatformIDs: []int32{}},
	})
	if !equal(changes["user-1"], []int32{1, 5}) {
		t.Fatalf("pushed changes after reconnect = %v", changes)
	}

	statusesAfterReconn, err := mgr.GetSubscribeUsersStatus(ctx)
	if err != nil || len(statusesAfterReconn) != 2 {
		t.Fatalf("GetSubscribeUsersStatus after reconnect push = %v, err = %v", statusesAfterReconn, err)
	}
	for _, s := range statusesAfterReconn {
		switch s.UserID {
		case "user-1":
			if s.Status != constant.Online || !equal(s.PlatformIDs, []int32{1, 5}) {
				t.Fatalf("unexpected user-1 status: %+v", s)
			}
		case "user-2":
			if s.Status != constant.Offline || len(s.PlatformIDs) != 0 {
				t.Fatalf("unexpected user-2 status: %+v", s)
			}
		default:
			t.Fatalf("unexpected userID: %s", s.UserID)
		}
	}

	if err := mgr.UnsubscribeUsersStatus(ctx, []string{"user-1"}); err != nil {
		t.Fatalf("UnsubscribeUsersStatus err = %v", err)
	}
	remaining, err := mgr.GetSubscribeUsersStatus(ctx)
	if err != nil || len(remaining) != 1 || remaining[0].UserID != "user-2" {
		t.Fatalf("GetSubscribeUsersStatus after unsubscribe = %v, err = %v", remaining, err)
	}
}
