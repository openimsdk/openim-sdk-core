package manager

import (
	"context"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/decorator"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"golang.org/x/sync/errgroup"
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

	gr, _ := errgroup.WithContext(ctx)
	gr.SetLimit(config.ErrGroupCommonLimit)
	m.sendSingleMessages(ctx, gr)
	m.sendGroupMessages(ctx, gr)
	return gr.Wait()
}

// sendSingleMessages see SendMessages
func (m *TestMsgManager) sendSingleMessages(ctx context.Context, gr *errgroup.Group) {
	for userNum := 0; userNum < vars.UserNum; userNum++ {
		ctx := vars.Contexts[userNum]
		testSDK := sdk.TestSDKs[userNum]
		gr.Go(func() error {
			friends, err := testSDK.GetAllFriends(ctx)
			if err != nil {
				return err
			}

			grr, _ := errgroup.WithContext(ctx)
			grr.SetLimit(config.ErrGroupSmallLimit)
			for _, friend := range friends {
				friend := friend
				if friend.FriendInfo != nil {
					for i := 0; i < vars.SingleMessageNum; i++ {
						grr.Go(func() error {
							msg, err := testSDK.SDK.Conversation().CreateTextMessage(ctx, fmt.Sprintf("my userID is %s", testSDK.UserID))
							if err != nil {
								return err
							}
							_, err = testSDK.SendSingleMsg(ctx, msg, friend.FriendInfo.FriendUserID)
							if err != nil {
								return err
							}
							return nil
						})
					}
				} else {
					fmt.Println("what`s this???")
				}
			}
			return grr.Wait()
		})
	}
	return
}

// sendGroupMessages see SendMessages
func (m *TestMsgManager) sendGroupMessages(ctx context.Context, gr *errgroup.Group) {
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
			grr, _ := errgroup.WithContext(ctx)
			grr.SetLimit(config.ErrGroupSmallLimit)
			for _, group := range sendGroups {
				group := group
				for i := 0; i < vars.GroupMessageNum; i++ {
					grr.Go(func() error {
						msg, err := testSDK.SDK.Conversation().CreateTextMessage(ctx, fmt.Sprintf("my userID is %s", testSDK.UserID))
						if err != nil {
							return err
						}
						_, err = testSDK.SendGroupMsg(ctx, msg, group)
						if err != nil {
							return err
						}
						return nil
					})
				}
			}
			return grr.Wait()
		})
	}
	return
}
