package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/decorator"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/progress"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/reerrgroup"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/sdk"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	sdkUtils "github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
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

	gr, cctx := reerrgroup.WithContext(ctx, config.ErrGroupCommonLimit)

	var (
		total int
		now   int
	)
	total = vars.LoginUserNum * 2
	p := progress.FuncNameBarPrint(cctx, gr, now, total)

	m.sendSingleMessages(ctx, gr, p)
	m.sendGroupMessages(ctx, gr, p)
	return gr.Wait()
}

// sendSingleMessages see SendMessages
func (m *TestMsgManager) sendSingleMessages(ctx context.Context, gr *reerrgroup.Group, p *progress.Progress) {
	for userNum := 0; userNum < vars.LoginUserNum; userNum++ {
		userNum := userNum
		ctx := vars.Contexts[userNum]
		testSDK := sdk.TestSDKs[userNum]
		gr.Go(func() error {

			friends, err := testSDK.GetAllFriends(ctx)
			if err != nil {
				return err
			}

			bar := progress.NewRemoveBar(fmt.Sprintf("%s:%s", "sendSingleMessages", utils.GetUserID(userNum)),
				0, len(friends)*vars.SingleMessageNum)
			p.AddBar(bar)

			friends = datautil.ShuffleSlice(friends)
			for _, friend := range friends {
				if friend != nil {
					for i := 0; i < vars.SingleMessageNum; i++ {
						msg, err := testSDK.SDK.Conversation().CreateTextMessage(ctx,
							fmt.Sprintf("count %d:my userID is %s", i, testSDK.UserID))
						if err != nil {
							return err
						}
						ctx = ccontext.WithOperationID(ctx, sdkUtils.OperationIDGenerator())
						t := time.Now()
						log.ZWarn(ctx, "sendSingleMessages begin", nil)
						_, err = testSDK.SendSingleMsg(ctx, msg, friend.FriendUserID)
						if err != nil {
							return err
						}
						log.ZWarn(ctx, "sendSingleMessages end", nil, "time cost:", time.Since(t))
						p.IncBar(bar)

						time.Sleep(time.Millisecond * 500)

					}
				} else {
					fmt.Println("what`s this???")
				}
			}
			log.ZWarn(ctx, "send over", nil, "userID", userNum)
			return nil
		})
	}
	return
}

// sendGroupMessages see SendMessages
func (m *TestMsgManager) sendGroupMessages(ctx context.Context, gr *reerrgroup.Group, p *progress.Progress) {
	for userNum := 0; userNum < vars.LoginUserNum; userNum++ {
		userNum := userNum
		ctx := vars.Contexts[userNum]
		testSDK := sdk.TestSDKs[userNum]
		gr.Go(func() error {
			groups, _, err := testSDK.GetAllJoinedGroups(ctx)
			if err != nil {
				return err
			}
			sendGroups := make([]string, 0)
			for _, group := range groups {
				if int(group.MemberCount) == vars.LargeGroupMemberNum || group.OwnerUserID == testSDK.UserID {
					// is larger group or created by oneself
					sendGroups = append(sendGroups, group.GroupID)
				}
			}

			bar := progress.NewRemoveBar(fmt.Sprintf("%s:%s", "sendGroupMessages", utils.GetUserID(userNum)),
				0, len(sendGroups)*vars.GroupMessageNum)
			p.AddBar(bar)

			sendGroups = datautil.ShuffleSlice(sendGroups)
			for _, group := range sendGroups {
				group := group
				for i := 0; i < vars.GroupMessageNum; i++ {
					msg, err := testSDK.SDK.Conversation().CreateTextMessage(ctx,
						fmt.Sprintf("count %d:my userID is %s", i, testSDK.UserID))
					if err != nil {
						return err
					}

					ctx = ccontext.WithOperationID(ctx, sdkUtils.OperationIDGenerator())
					t := time.Now()
					log.ZWarn(ctx, "sendGroupMessages begin", nil)
					_, err = testSDK.SendGroupMsg(ctx, msg, group)
					if err != nil {
						return err
					}
					log.ZWarn(ctx, "sendGroupMessages end", nil, "time cost:", time.Since(t))

					p.IncBar(bar)
					time.Sleep(time.Millisecond * 500)
				}
			}
			return nil
		})
	}
	return
}
