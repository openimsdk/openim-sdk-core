package super_group

import (
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
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

func (s *SuperGroup) DoNotification(msg *api.MsgData, conversationCh chan common.Cmd2Value) {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	if msg.SendTime < s.loginTime {
		log.Warn(operationID, "ignore notification ", msg.ClientMsgID, msg.ServerMsgID, msg.Seq, msg.ContentType)
		return
	}
	go func() {
		switch msg.ContentType {
		case constant.SuperGroupUpdateNotification:
			s.SyncJoinedGroupList(operationID)
			cmd := sdk_struct.CmdJoinedSuperGroup{}
			err := common.TriggerCmdJoinedSuperGroup(cmd, s.joinedSuperGroupCh)
			if err != nil {
				log.Error(operationID, "TriggerCmdJoinedSuperGroup failed ", err.Error(), cmd)
			}
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
	log.Debug(operationID, "api args: ", apiReq)
	err := s.p.PostReturn(constant.GetJoinedSuperGroupListRouter, apiReq, &result)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return result, nil
}
