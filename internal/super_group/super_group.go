package super_group

import (
	"errors"
	"github.com/golang/protobuf/proto"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

type SuperGroup struct {
	loginUserID        string
	db                 db_interface.DataBase
	p                  *ws.PostApi
	loginTime          int64
	joinedSuperGroupCh chan common.Cmd2Value
	heartbeatCmdCh     chan common.Cmd2Value
}

func (s *SuperGroup) SetLoginTime(loginTime int64) {
	s.loginTime = loginTime
}

func NewSuperGroup(loginUserID string, db db_interface.DataBase, p *ws.PostApi, joinedSuperGroupCh chan common.Cmd2Value, heartbeatCmdCh chan common.Cmd2Value) *SuperGroup {
	return &SuperGroup{loginUserID: loginUserID, db: db, p: p, joinedSuperGroupCh: joinedSuperGroupCh, heartbeatCmdCh: heartbeatCmdCh}
}

func (s *SuperGroup) DoNotification(msg *api.MsgData, ch chan common.Cmd2Value) {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID, msg.String())
	if msg.SendTime < s.loginTime || s.loginTime == 0 {
		log.Warn(operationID, "ignore notification ", msg.ClientMsgID, msg.ServerMsgID, msg.Seq, msg.ContentType)
		return
	}
	go func() {
		switch msg.ContentType {
		case constant.SuperGroupUpdateNotification:
			_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{Action: constant.SyncConversation, Args: operationID}, ch)
			s.SyncJoinedGroupList(operationID)
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
	}()
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

func (s *SuperGroup) GetGroupInfoFromLocal2Svr(groupID string) (*model_struct.LocalGroup, error) {
	localGroup, err := s.db.GetSuperGroupInfoByGroupID(groupID)
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

func (s *SuperGroup) GetJoinedGroupIDListFromSvr(operationID string) ([]string, error) {
	result, err := s.getJoinedGroupListFromSvr(operationID)
	if err != nil {
		return nil, utils.Wrap(err, "SuperGroup get err")
	}
	var groupIDList []string
	for _, v := range result {
		groupIDList = append(groupIDList, v.GroupID)
	}
	return groupIDList, nil
}
