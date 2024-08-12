package manager

import (
	"context"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/decorator"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/reerrgroup"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/tools/log"
	"sync/atomic"
)

type TestMsgManager struct {
	*MetaManager
}

func NewMsgManager(m *MetaManager) *TestMsgManager {
	return &TestMsgManager{m}
}

// SendMessages send messages.The rules are: each user sends `vars.SingleMessageNum` messages to all friends,
// and sends `vars.GroupMessageNum` messages to all large groups and common groups created by themselves.
func (m *TestMsgManager) SendMessages(ctx context.Context) error {
	defer decorator.FuncLog(ctx)()

	gr, cctx := reerrgroup.WithContext(ctx, config.ErrGroupMiddleSmallLimit)

	var (
		total    atomic.Int64
		progress atomic.Int64
	)
	total.Add(int64(vars.UserNum * 2))
	utils.FuncProgressBarPrint(cctx, gr, &progress, &total)

	m.sendSingleMessages(ctx, gr)
	// prevent lock database
	//gr.WaitTaskDone()
	m.sendGroupMessages(ctx, gr)
	return gr.Wait()
}

// sendSingleMessages see SendMessages
func (m *TestMsgManager) sendSingleMessages(ctx context.Context, gr *reerrgroup.Group) {
	for userNum := 0; userNum < vars.UserNum; userNum++ {
		ctx := vars.Contexts[userNum]
		testSDK := sdk.TestSDKs[userNum]
		gr.Go(func() error {
			friends, err := testSDK.GetAllFriends(ctx)
			if err != nil {
				return err
			}
			for _, friend := range friends {
				if friend.FriendInfo != nil {
					for i := 0; i < vars.SingleMessageNum; i++ {
						msg, err := testSDK.SDK.Conversation().CreateTextMessage(ctx,
							fmt.Sprintf("count %d:my userID is %s", i, testSDK.UserID))
						if err != nil {
							return err
						}
						_, err = testSDK.SendSingleMsg(ctx, msg, friend.FriendInfo.FriendUserID)
						if err != nil {
							return err
						}
					}
				} else {
					fmt.Println("what`s this???")
				}
			}
			log.ZError(ctx, "send over", nil, "userID", userNum)
			return nil
		})
	}
	return
}

// sendGroupMessages see SendMessages
func (m *TestMsgManager) sendGroupMessages(ctx context.Context, gr *reerrgroup.Group) {
	for userNum := 0; userNum < vars.UserNum; userNum++ {
		ctx := vars.Contexts[userNum]
		testSDK := sdk.TestSDKs[userNum]
		gr.Go(func() error {
			groups, _, err := testSDK.GetAllJoinedGroups(ctx)
			if err != nil {
				return err
			}
			sendGroups := make([]string, 0)
			for _, group := range groups {
				if int(group.MemberCount) == vars.UserNum || group.OwnerUserID == testSDK.UserID {
					// is larger group or created by oneself
					sendGroups = append(sendGroups, group.GroupID)
				}
			}
			for _, group := range sendGroups {
				group := group
				for i := 0; i < vars.GroupMessageNum; i++ {
					msg, err := testSDK.SDK.Conversation().CreateTextMessage(ctx,
						fmt.Sprintf("count %d:my userID is %s", i, testSDK.UserID))
					if err != nil {
						return err
					}
					_, err = testSDK.SendGroupMsg(ctx, msg, group)
					if err != nil {
						return err
					}
				}
			}
			return nil
		})
	}
	return
}
