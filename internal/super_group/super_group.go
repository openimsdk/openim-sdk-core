package super_group

import (
	"context"
	"errors"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/syncer"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/golang/protobuf/proto"
)

func NewSuperGroup(loginUserID string, db db_interface.DataBase, joinedSuperGroupCh chan common.Cmd2Value, heartbeatCmdCh chan common.Cmd2Value) *SuperGroup {
	s := &SuperGroup{loginUserID: loginUserID, db: db, joinedSuperGroupCh: joinedSuperGroupCh, heartbeatCmdCh: heartbeatCmdCh}
	s.initSyncer()
	return s
}

type SuperGroup struct {
	loginUserID        string
	db                 db_interface.DataBase
	loginTime          int64
	joinedSuperGroupCh chan common.Cmd2Value
	heartbeatCmdCh     chan common.Cmd2Value
	syncerGroup        *syncer.Syncer[*model_struct.LocalGroup, string]
}

func (s *SuperGroup) initSyncer() {
	s.syncerGroup = syncer.New(func(ctx context.Context, value *model_struct.LocalGroup) error {
		return s.db.InsertSuperGroup(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalGroup) error {
		return s.db.DeleteGroup(ctx, value.GroupID)
	}, func(ctx context.Context, server, local *model_struct.LocalGroup) error {
		return s.db.UpdateGroup(ctx, server)
	}, func(value *model_struct.LocalGroup) string {
		return value.GroupID
	}, nil, nil)
}

func (s *SuperGroup) SetLoginTime(loginTime int64) {
	s.loginTime = loginTime
}

func (s *SuperGroup) DoNotification(msg *sdkws.MsgData, ch chan common.Cmd2Value, operationID string) {
	ctx := context.Background()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID, msg.String())
	if msg.SendTime < s.loginTime || s.loginTime == 0 {
		log.Warn(operationID, "ignore notification ", msg.ClientMsgID, msg.ServerMsgID, msg.Seq, msg.ContentType)
		return
	}
	switch msg.ContentType {
	case constant.SuperGroupUpdateNotification:
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{Action: constant.SyncConversation, Args: operationID}, ch)
		if err := s.SyncJoinedGroupList(ctx); err != nil {
			// todo log
		}
		cmd := sdk_struct.CmdJoinedSuperGroup{OperationID: operationID}
		err := common.TriggerCmdJoinedSuperGroup(cmd, s.joinedSuperGroupCh)
		if err != nil {
			log.Error(operationID, "TriggerCmdJoinedSuperGroup failed ", err.Error(), cmd)
			return
		}
		err = common.TriggerCmdWakeUp(s.heartbeatCmdCh)
		if err != nil {
			log.Error(operationID, "TriggerCmdWakeUp failed ", err.Error())
		}

		log.Info(operationID, "constant.SuperGroupUpdateNotification", msg.String())

	case constant.MsgDeleteNotification:
		var tips api.TipsComm
		var elem api.MsgDeleteNotificationElem
		_ = proto.Unmarshal(msg.Content, &tips)
		_ = utils.JsonStringToStruct(tips.JsonDetail, &elem)
		//if elem.GroupID != nil {
		//
		//}
	default:
		log.Error(operationID, "ContentType tip failed ", msg.ContentType)
	}
}

func (g *SuperGroup) getJoinedGroupListFromSvr(ctx context.Context) ([]*sdkws.GroupInfo, error) {
	fn := func(resp *group.GetJoinedGroupListResp) []*sdkws.GroupInfo { return resp.Groups }
	req := &group.GetJoinedGroupListReq{FromUserID: g.loginUserID, Pagination: &sdkws.RequestPagination{}}
	return util.GetPageAll(ctx, constant.GetJoinedGroupListRouter, req, fn)
}

func (s *SuperGroup) GetGroupInfoFromLocal2Svr(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
	localGroup, err := s.db.GetSuperGroupInfoByGroupID(ctx, groupID)
	if err == nil {
		return localGroup, nil
	}
	groupIDList := []string{groupID}
	//operationID := utils.OperationIDGenerator()
	svrGroup, err := s.getGroupsInfoFromSvr(ctx, groupIDList)
	if err != nil {
		return nil, err
	}
	if len(svrGroup) == 0 {
		return nil, utils.Wrap(errors.New("no group"), "")
	}
	return ServerGroupToLocalGroup(svrGroup[0]), nil
}

func (s *SuperGroup) getGroupsInfoFromSvr(ctx context.Context, groupIDList []string) ([]*sdkws.GroupInfo, error) {
	resp, err := util.CallApi[group.GetGroupsInfoResp](ctx, constant.GetSuperGroupsInfoRouter, group.GetGroupsInfoReq{GroupIDs: groupIDList})
	if err != nil {
		return nil, err
	}
	return resp.GroupInfos, nil
}

func (s *SuperGroup) GetJoinedGroupIDListFromSvr(ctx context.Context) ([]string, error) {
	result, err := s.getJoinedGroupListFromSvr(ctx)
	if err != nil {
		return nil, utils.Wrap(err, "SuperGroup get err")
	}
	var groupIDList []string
	for _, v := range result {
		groupIDList = append(groupIDList, v.GroupID)
	}
	return groupIDList, nil
}
