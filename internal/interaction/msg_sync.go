package interaction

import (
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

type SeqPair struct {
	BeginSeq uint32
	EndSeq   uint32
}

type MsgSync struct {
	*db.DataBase
	*Ws
	loginUserID        string
	conversationCh     chan common.Cmd2Value
	PushMsgAndMaxSeqCh chan common.Cmd2Value

	selfMsgSync *SelfMsgSync
	//selfMsgSyncLatestModel *SelfMsgSyncLatestModel
	superGroupMsgSync *SuperGroupMsgSync
}

func (m *MsgSync) compareSeq() {
	operationID := utils.OperationIDGenerator()
	m.selfMsgSync.compareSeq(operationID)
	m.superGroupMsgSync.compareSeq(operationID)
}

func (m *MsgSync) doMaxSeq(cmd common.Cmd2Value) {
	m.selfMsgSync.doMaxSeq(cmd)
	m.superGroupMsgSync.doMaxSeq(cmd)
}

func (m *MsgSync) doPushMsg(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync).Msg
	switch msg.SessionType {
	case constant.SuperGroupChatType:
		m.superGroupMsgSync.doPushMsg(cmd)
	default:
		m.selfMsgSync.doPushMsg(cmd)
	}
}

func (m *MsgSync) Work(cmd common.Cmd2Value) {
	switch cmd.Cmd {
	case constant.CmdPushMsg:
		m.doPushMsg(cmd)
	case constant.CmdMaxSeq:
		m.doMaxSeq(cmd)
	default:
		log.Error("", "cmd failed ", cmd.Cmd)
	}
}

func (m *MsgSync) GetCh() chan common.Cmd2Value {
	return m.PushMsgAndMaxSeqCh
}

func NewMsgSync(dataBase *db.DataBase, ws *Ws, loginUserID string, ch chan common.Cmd2Value, pushMsgAndMaxSeqCh chan common.Cmd2Value, joinedSuperGroupCh chan common.Cmd2Value) *MsgSync {
	p := &MsgSync{DataBase: dataBase, Ws: ws, loginUserID: loginUserID, conversationCh: ch, PushMsgAndMaxSeqCh: pushMsgAndMaxSeqCh}
	p.superGroupMsgSync = NewSuperGroupMsgSync(dataBase, ws, loginUserID, ch, joinedSuperGroupCh)
	p.selfMsgSync = NewSelfMsgSync(dataBase, ws, loginUserID, ch)
	//	p.selfMsgSync = NewSelfMsgSyncLatestModel(dataBase, ws, loginUserID, ch)
	p.compareSeq()
	go common.DoListener(p)
	return p
}
