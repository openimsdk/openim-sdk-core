package super_group

import (
	"context"
	"errors"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/syncer"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"

	"github.com/golang/protobuf/proto"

	"github.com/golang/protobuf/proto"
)

func NewSuperGroup(loginUserID string, db db_interface.DataBase, p *ws.PostApi, joinedSuperGroupCh chan common.Cmd2Value, heartbeatCmdCh chan common.Cmd2Value) *SuperGroup {
	return &SuperGroup{loginUserID: loginUserID, db: db, p: p, joinedSuperGroupCh: joinedSuperGroupCh, heartbeatCmdCh: heartbeatCmdCh}
}

type SuperGroup struct {
	loginUserID        string
	db                 db_interface.DataBase
	p                  *ws.PostApi
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

func (s *SuperGroup) DoNotification(msg *api.MsgData, ch chan common.Cmd2Value, operationID string) {
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

func (s *SuperGroup) getJoinedGroupListFromSvr(operationID string) ([]*api.GroupInfo, error) {
	apiReq := api.GetJoinedSuperGroupReq{}
	apiReq.OperationID = operationID
	apiReq.FromUserID = s.loginUserID
	var result []*api.GroupInfo
	log.Debug(operationID, "super group api args: ", apiReq)
	err := s.p.PostReturn(constant.GetJoinedSuperGroupListRouter, apiReq, &result)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	log.Debug(operationID, "super group api result: ", result)
	return result, nil
}

func (s *SuperGroup) GetGroupInfoFromLocal2Svr(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
	localGroup, err := s.db.GetSuperGroupInfoByGroupID(ctx, groupID)
	if err == nil {
		return localGroup, nil
	}
	groupIDList := []string{groupID}
	operationID := utils.OperationIDGenerator()
	svrGroup, err := s.getGroupsInfoFromSvr(groupIDList, operationID)
	if err == nil && len(svrGroup) == 1 {
		transfer := common.TransferToLocalGroupInfo(svrGroup)
		return transfer[0], nil
	}
	if err != nil {
		return nil, utils.Wrap(err, "")
	} else {
		return nil, utils.Wrap(errors.New("no group"), "")
	}
}

func (s *SuperGroup) getGroupsInfoFromSvr(groupIDList []string, operationID string) ([]*api.GroupInfo, error) {
	apiReq := api.GetSuperGroupsInfoReq{}
	apiReq.GroupIDList = groupIDList
	apiReq.OperationID = operationID
	var groupInfoList []*api.GroupInfo
	err := s.p.PostReturn(constant.GetSuperGroupsInfoRouter, apiReq, &groupInfoList)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return groupInfoList, nil
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
