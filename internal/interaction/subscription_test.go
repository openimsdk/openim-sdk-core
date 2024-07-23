package interaction

import (
	"errors"
	"testing"
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
