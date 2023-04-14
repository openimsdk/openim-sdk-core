package full

import (
	"context"
	"fmt"
	"open_im_sdk/internal/cache"
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/group"
	"open_im_sdk/internal/super_group"
	"open_im_sdk/internal/user"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	api "open_im_sdk/pkg/server_api_params"
)

type Full struct {
	user       *user.User
	friend     *friend.Friend
	group      *group.Group
	ch         chan common.Cmd2Value
	userCache  *cache.Cache
	db         db_interface.DataBase
	SuperGroup *super_group.SuperGroup
}

func (u *Full) Group() *group.Group {
	return u.group
}

func NewFull(user *user.User, friend *friend.Friend, group *group.Group, ch chan common.Cmd2Value, userCache *cache.Cache, db db_interface.DataBase, superGroup *super_group.SuperGroup) *Full {
	return &Full{user: user, friend: friend, group: group, ch: ch, userCache: userCache, db: db, SuperGroup: superGroup}
}
func (u *Full) getUsersInfo(ctx context.Context, userIDList sdk.GetUsersInfoParam) (sdk.GetUsersInfoCallback, error) {
	friendList, err := u.db.GetFriendInfoList(ctx, userIDList)
	if err != nil {
		return nil, err
	}
	blackList, err := u.db.GetBlackInfoList(ctx, userIDList)
	if err != nil {
		return nil, err
	}
	notIn := make([]string, 0)
	for _, v := range userIDList {
		inFriendList := 0
		for _, friend := range friendList {
			if v == friend.FriendUserID {
				inFriendList = 1
				break
			}
		}
		inBlackList := 0
		for _, black := range blackList {
			if v == black.BlockUserID {
				inBlackList = 1
				break
			}
		}
		if inFriendList == 0 && inBlackList == 0 {
			notIn = append(notIn, v)
		}
	}
	//from svr
	publicList := make([]*api.PublicUserInfo, 0)
	if len(notIn) > 0 {
		publicList, err = u.user.GetUsersInfoFromSvr(ctx, notIn)
		go func() {
			for _, v := range publicList {
				u.userCache.Update(v.UserID, v.FaceURL, v.Nickname)
				//Update the faceURL and nickname information of the local chat history with non-friends
				//_ = u.user.UpdateMsgSenderFaceURLAndSenderNickname(v.UserID, v.FaceURL, v.Nickname, constant.SingleChatType)
				_ = common.TriggerCmdUpdateMessage(common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName, Args: common.UpdateMessageInfo{UserID: v.UserID, FaceURL: v.FaceURL, Nickname: v.Nickname}}, u.ch)
				//Update session information of local non-friends
				_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{Action: constant.UpdateConFaceUrlAndNickName, Args: common.SourceIDAndSessionType{SourceID: v.UserID, SessionType: constant.SingleChatType}}, u.ch)
			}
		}()
	}
	return common.MergeUserResult(publicList, friendList, blackList), nil
}

func (u *Full) GetGroupInfoFromLocal2Svr(ctx context.Context, groupID string, sessionType int32) (*model_struct.LocalGroup, error) {
	switch sessionType {
	case constant.GroupChatType:
		return u.group.GetGroupInfoFromLocal2Svr(ctx, groupID)
	case constant.SuperGroupChatType:
		return u.GetGroupInfoByGroupID(groupID)
	default:
		return nil, fmt.Errorf("sessionType is not support %d", sessionType)
	}
}
func (u *Full) GetReadDiffusionGroupIDList(ctx context.Context) ([]string, error) {
	g1, err1 := u.group.GetJoinedDiffusionGroupIDListFromSvr(ctx)
	g2, err2 := u.SuperGroup.GetJoinedGroupIDListFromSvr(ctx)
	var groupIDList []string
	if err1 == nil {
		groupIDList = append(groupIDList, g1...)
	}
	if err2 == nil {
		groupIDList = append(groupIDList, g2...)
	}
	var err error
	if err1 != nil {
		err = err1
	}
	if err2 != nil {
		err = err2
	}
	return groupIDList, err
}
