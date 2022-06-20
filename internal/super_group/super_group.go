package super_group

import (
	"errors"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

type SuperGroup struct {
	loginUserID        string
	db                 *db.DataBase
	p                  *ws.PostApi
	loginTime          int64
	joinedSuperGroupCh chan common.Cmd2Value
}

func (s *SuperGroup) SetLoginTime(loginTime int64) {
	s.loginTime = loginTime
}

func NewSuperGroup(loginUserID string, db *db.DataBase, p *ws.PostApi, joinedSuperGroupCh chan common.Cmd2Value) *SuperGroup {
	return &SuperGroup{loginUserID: loginUserID, db: db, p: p, joinedSuperGroupCh: joinedSuperGroupCh}
}

func (s *SuperGroup) DoNotification(msg *api.MsgData, _ chan common.Cmd2Value) {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID, msg.String())
	if msg.SendTime < s.loginTime {
		log.Warn(operationID, "ignore notification ", msg.ClientMsgID, msg.ServerMsgID, msg.Seq, msg.ContentType)
		return
	}
	go func() {
		switch msg.ContentType {
		case constant.SuperGroupUpdateNotification:
			s.SyncJoinedGroupList(operationID)
			cmd := sdk_struct.CmdJoinedSuperGroup{OperationID: operationID}
			err := common.TriggerCmdJoinedSuperGroup(cmd, s.joinedSuperGroupCh)
			if err != nil {
				log.Error(operationID, "TriggerCmdJoinedSuperGroup failed ", err.Error(), cmd)
				return
			}
			log.Info(operationID, "constant.SuperGroupUpdateNotification", msg.String())
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
