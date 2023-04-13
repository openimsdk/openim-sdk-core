package group

import (
	"context"
	"errors"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/syncer"
	"open_im_sdk/pkg/utils"
	"sync"
)

func NewGroup(loginUserID string, db db_interface.DataBase, p *ws.PostApi,
	joinedSuperGroupCh chan common.Cmd2Value, heartbeatCmdCh chan common.Cmd2Value,
	conversationCh chan common.Cmd2Value) *Group {
	g := &Group{
		loginUserID:        loginUserID,
		db:                 db,
		p:                  p,
		joinedSuperGroupCh: joinedSuperGroupCh,
		heartbeatCmdCh:     heartbeatCmdCh,
		conversationCh:     conversationCh,
	}
	g.initSyncer()
	return g
}

// //utils.GetCurrentTimestampByMill()
type Group struct {
	listener                open_im_sdk_callback.OnGroupListener
	loginUserID             string
	db                      db_interface.DataBase
	groupSyncer             *syncer.Syncer[*model_struct.LocalGroup, string]
	groupMemberSyncer       *syncer.Syncer[*model_struct.LocalGroupMember, [2]string]
	groupRequestSyncer      *syncer.Syncer[*model_struct.LocalGroupRequest, [2]string]
	groupAdminRequestSyncer *syncer.Syncer[*model_struct.LocalAdminGroupRequest, [2]string]
	p                       *ws.PostApi
	loginTime               int64
	joinedSuperGroupCh      chan common.Cmd2Value
	heartbeatCmdCh          chan common.Cmd2Value

	conversationCh chan common.Cmd2Value
	//	memberSyncMutex sync.RWMutex

	listenerForService open_im_sdk_callback.OnListenerForService
}

func (g *Group) initSyncer() {
	g.groupSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalGroup) error {
		return g.db.InsertGroup(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalGroup) error {
		return g.db.DeleteGroup(ctx, value.GroupID)
	}, func(ctx context.Context, server, local *model_struct.LocalGroup) error {
		return g.db.UpdateGroup(ctx, server)
	}, func(value *model_struct.LocalGroup) string {
		return value.GroupID
	}, nil, nil)

	g.groupMemberSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalGroupMember) error {
		return g.db.InsertGroupMember(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalGroupMember) error {
		return g.db.DeleteGroupMember(ctx, value.GroupID, value.UserID)
	}, func(ctx context.Context, server, local *model_struct.LocalGroupMember) error {
		return g.db.UpdateGroupMember(ctx, server)
	}, func(value *model_struct.LocalGroupMember) [2]string {
		return [...]string{value.GroupID, value.UserID}
	}, nil, nil)

	g.groupRequestSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalGroupRequest) error {
		return g.db.InsertGroupRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalGroupRequest) error {
		return g.db.DeleteGroupRequest(ctx, value.GroupID, value.UserID)
	}, func(ctx context.Context, server, local *model_struct.LocalGroupRequest) error {
		return g.db.UpdateGroupRequest(ctx, server)
	}, func(value *model_struct.LocalGroupRequest) [2]string {
		return [...]string{value.GroupID, value.UserID}
	}, nil, nil)

	g.groupAdminRequestSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalAdminGroupRequest) error {
		return g.db.InsertAdminGroupRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalAdminGroupRequest) error {
		return g.db.DeleteAdminGroupRequest(ctx, value.GroupID, value.UserID)
	}, func(ctx context.Context, server, local *model_struct.LocalAdminGroupRequest) error {
		return g.db.UpdateAdminGroupRequest(ctx, server)
	}, func(value *model_struct.LocalAdminGroupRequest) [2]string {
		return [...]string{value.GroupID, value.UserID}
	}, nil, nil)

}

func (g *Group) SetGroupListener(callback open_im_sdk_callback.OnGroupListener) {
	if callback == nil {
		return
	}
	g.listener = callback
}

func (g *Group) LoginTime() int64 {
	return g.loginTime
}

func (g *Group) SetLoginTime(loginTime int64) {
	g.loginTime = loginTime
}

func (g *Group) SetListenerForService(listener open_im_sdk_callback.OnListenerForService) {
	g.listenerForService = listener
}

func (g *Group) GetGroupOwnerIDAndAdminIDList(ctx context.Context, groupID string) (ownerID string, adminIDList []string, err error) {
	localGroup, err := g.db.GetGroupInfoByGroupID(ctx, groupID)
	if err != nil {
		return "", nil, err
	}
	adminIDList, err = g.db.GetGroupAdminID(ctx, groupID)
	if err != nil {
		return "", nil, err
	}
	return localGroup.OwnerUserID, adminIDList, nil
}

func (g *Group) GetGroupInfoFromLocal2Svr(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
	localGroup, err := g.db.GetGroupInfoByGroupID(ctx, groupID)
	if err == nil {
		return localGroup, nil
	}
	groupIDList := []string{groupID}
	operationID := utils.OperationIDGenerator()
	svrGroup, err := g.getGroupsInfoFromSvr(groupIDList, operationID)
	if err == nil && len(svrGroup) == 1 {
		transfer := common.TransferToLocalGroupInfo(svrGroup)
		return transfer[0], nil
	}
	if err != nil {
		return nil, utils.Wrap(err, "get groupInfo from server err ")
	} else {
		return nil, utils.Wrap(errors.New("server not this group"), "")
	}
}

func (g *Group) getGroupsInfoFromSvr(groupIDList []string, operationID string) ([]*api.GroupInfo, error) {
	apiReq := api.GetGroupInfoReq{}
	apiReq.GroupIDList = groupIDList
	apiReq.OperationID = operationID
	var groupInfoList []*api.GroupInfo
	err := g.p.PostReturn(constant.GetGroupsInfoRouter, apiReq, &groupInfoList)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return groupInfoList, nil
}

func (g *Group) getGroupAbstractInfoFromSvr(groupID string, operationID string) (*api.GetGroupAbstractInfoResp, error) {
	apiReq := api.GetGroupAbstractInfoReq{}
	apiReq.GroupID = groupID
	apiReq.OperationID = operationID
	var groupAbstractInfoResp api.GetGroupAbstractInfoResp
	err := g.p.Post2UnmarshalRespReturn(constant.GetGroupAbstractInfoRouter, apiReq, &groupAbstractInfoResp)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID+" "+groupID)
	}
	return &groupAbstractInfoResp, nil
}

func (g *Group) getJoinedGroupListFromSvr(operationID string) ([]*api.GroupInfo, error) {
	apiReq := api.GetJoinedGroupListReq{}
	apiReq.OperationID = operationID
	apiReq.FromUserID = g.loginUserID
	var result []*api.GroupInfo
	log.Debug(operationID, "api args: ", apiReq)
	err := g.p.PostReturn(constant.GetJoinedGroupListRouter, apiReq, &result)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return result, nil
}

func (g *Group) GetJoinedDiffusionGroupIDListFromSvr(operationID string) ([]string, error) {
	result, err := g.getJoinedGroupListFromSvr(operationID)
	if err != nil {
		return nil, utils.Wrap(err, "working group get err")
	}
	var groupIDList []string
	for _, v := range result {
		if v.GroupType == constant.WorkingGroup {
			groupIDList = append(groupIDList, v.GroupID)
		}
	}
	return groupIDList, nil
}

func (g *Group) SyncJoinedGroupList(ctx context.Context) {
	if err := g.SyncJoinedGroup(ctx); err != nil {
		// tood log
	}
}

func (g *Group) SyncJoinedGroupMemberForFirstLogin(ctx context.Context) {
	groups, err := g.GetAndSyncJoinedGroup(ctx)
	if err != nil {
		log.Error("SyncJoinedGroupMemberForFirstLogin", "GetAndSyncJoinedGroup failed", err.Error())
		return
	}
	var wg sync.WaitGroup
	for _, group := range groups {
		wg.Add(1)
		go func(groupID string) {
			defer wg.Done()
			if err := g.SyncGroupMember(ctx, groupID); err != nil {
				log.Error("SyncJoinedGroupMemberForFirstLogin", "SyncGroupMember failed", err.Error())
			}
		}(group.GroupID)
	}
	wg.Wait()
}

func (g *Group) getGroupAllMemberSplitByGroupIDFromSvr(groupID string, operationID string) ([]*api.GroupMemberFullInfo, error) {
	var apiReq api.GetGroupAllMemberReq
	apiReq.OperationID = operationID
	apiReq.GroupID = groupID
	var result []*api.GroupMemberFullInfo
	var page int32
	for {
		apiReq.Offset = page * constant.SplitGetGroupMemberNum
		apiReq.Count = constant.SplitGetGroupMemberNum
		var realData []*api.GroupMemberFullInfo
		err := g.p.PostReturn(constant.GetGroupAllMemberListRouter, apiReq, &realData)
		if err != nil {
			log.Error(operationID, "GetGroupAllMemberListRouter failed ", constant.GetGroupAllMemberListRouter, apiReq)
			return result, utils.Wrap(err, apiReq.OperationID)
		}
		log.Info(operationID, "GetGroupAllMemberListRouter result len: ", len(realData), groupID)
		result = append(result, realData...)
		if apiReq.Count > int32(len(realData)) {
			break
		}
		page++
	}
	return result, nil
}
