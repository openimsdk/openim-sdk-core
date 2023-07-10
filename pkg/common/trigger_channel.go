package common

import (
	"errors"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"runtime"
	"time"
)

func TriggerCmdJoinedSuperGroup(cmd sdk_struct.CmdJoinedSuperGroup, joinedSuperGroupCh chan Cmd2Value) error {
	if joinedSuperGroupCh == nil {
		return utils.Wrap(errors.New("ch == nil"), "")
	}
	c2v := Cmd2Value{Cmd: constant.CmdJoinedSuperGroup, Value: cmd}
	return sendCmd(joinedSuperGroupCh, c2v, 100)
}

func TriggerCmdNewMsgCome(msg sdk_struct.CmdNewMsgComeToConversation, conversationCh chan Cmd2Value) error {
	if conversationCh == nil {
		return utils.Wrap(errors.New("ch == nil"), "")
	}
	if len(msg.MsgList) == 0 {
		return nil
	}

	c2v := Cmd2Value{Cmd: constant.CmdNewMsgCome, Value: msg}
	return sendCmd(conversationCh, c2v, 100)
}
func TriggerCmdSuperGroupMsgCome(msg sdk_struct.CmdNewMsgComeToConversation, conversationCh chan Cmd2Value) error {
	if conversationCh == nil {
		return utils.Wrap(errors.New("ch == nil"), "")
	}
	//if len(msg.MsgList) == 0 {
	//	return nil
	//}

	c2v := Cmd2Value{Cmd: constant.CmdSuperGroupMsgCome, Value: msg}
	return sendCmd(conversationCh, c2v, 100)
}

func TriggerCmdLogout(ch chan Cmd2Value) error {
	if ch == nil {
		return utils.Wrap(errors.New("ch == nil"), "")
	}
	c2v := Cmd2Value{Cmd: constant.CmdLogout, Value: nil}
	return sendCmd(ch, c2v, 100)
}

func TriggerCmdWakeUp(ch chan Cmd2Value) error {
	if ch == nil {
		return utils.Wrap(errors.New("ch == nil"), "")
	}
	c2v := Cmd2Value{Cmd: constant.CmdWakeUp, Value: nil}
	return sendCmd(ch, c2v, 100)
}

func TriggerCmdDeleteConversationAndMessage(sourceID, conversationID string, sessionType int, conversationCh chan Cmd2Value) error {
	if conversationCh == nil {
		return utils.Wrap(errors.New("ch == nil"), "")
	}
	c2v := Cmd2Value{
		Cmd:   constant.CmdDeleteConversation,
		Value: DeleteConNode{SourceID: sourceID, ConversationID: conversationID, SessionType: sessionType},
	}

	return sendCmd(conversationCh, c2v, 100)
}

func TriggerCmdSyncReactionExtensions(node SyncReactionExtensionsNode, conversationCh chan Cmd2Value) error {
	if conversationCh == nil {
		return utils.Wrap(errors.New("ch == nil"), "")
	}
	c2v := Cmd2Value{
		Cmd:   constant.CmSyncReactionExtensions,
		Value: node,
	}

	return sendCmd(conversationCh, c2v, 100)
}
func TriggerCmdUpdateConversation(node UpdateConNode, conversationCh chan Cmd2Value) error {
	c2v := Cmd2Value{
		Cmd:   constant.CmdUpdateConversation,
		Value: node,
	}

	return sendCmd(conversationCh, c2v, 100)
}
func TriggerCmdUpdateMessage(node UpdateMessageNode, conversationCh chan Cmd2Value) error {
	c2v := Cmd2Value{
		Cmd:   constant.CmdUpdateMessage,
		Value: node,
	}

	return sendCmd(conversationCh, c2v, 100)
}

func TriggerCmdPushMsg(msg sdk_struct.CmdPushMsgToMsgSync, ch chan Cmd2Value) error {
	if ch == nil {
		return utils.Wrap(errors.New("ch == nil"), "")
	}

	c2v := Cmd2Value{Cmd: constant.CmdPushMsg, Value: msg}
	return sendCmd(ch, c2v, 100)
}

func TriggerCmdMaxSeq(seq sdk_struct.CmdMaxSeqToMsgSync, ch chan Cmd2Value) error {
	if ch == nil {
		return utils.Wrap(errors.New("ch == nil"), "")
	}
	c2v := Cmd2Value{Cmd: constant.CmdMaxSeq, Value: seq}
	return sendCmd(ch, c2v, 100)
}

type DeleteConNode struct {
	SourceID       string
	ConversationID string
	SessionType    int
}
type SyncReactionExtensionsNode struct {
	OperationID string
	Action      int
	Args        interface{}
}
type UpdateConNode struct {
	ConID  string
	Action int //1 Delete the conversation; 2 Update the latest news in the conversation or add a conversation; 3 Put a conversation on the top;
	// 4 Cancel a conversation on the top, 5 Messages are not read and set to 0, 6 New conversations
	Args interface{}
}
type UpdateMessageNode struct {
	Action int
	Args   interface{}
}

type Cmd2Value struct {
	Cmd   string
	Value interface{}
}

type UpdateMessageInfo struct {
	UserID   string
	FaceURL  string
	Nickname string
	GroupID  string
}

type SourceIDAndSessionType struct {
	SourceID    string
	SessionType int
}

func UnInitAll(conversationCh chan Cmd2Value) error {
	c2v := Cmd2Value{Cmd: constant.CmdUnInit}
	return sendCmd(conversationCh, c2v, 100)
}

type goroutine interface {
	Work(cmd Cmd2Value)
	GetCh() chan Cmd2Value
}

func DoListener(Li goroutine) {
	log.Info("internal", "doListener start.", Li.GetCh())
	for {
		select {
		case cmd := <-Li.GetCh():
			if cmd.Cmd == constant.CmdUnInit {
				log.Warn("", "close doListener channel ", Li.GetCh())
				//	close(Li.GetCh())
				runtime.Goexit()
			}
			//	log.Info("doListener work.")
			Li.Work(cmd)
		}
	}
}

func sendCmd(ch chan Cmd2Value, value Cmd2Value, timeout int64) error {
	var flag = 0
	select {
	case ch <- value:
		flag = 1
	case <-time.After(time.Millisecond * time.Duration(timeout)):
		flag = 2
	}
	if flag == 1 {
		return nil
	} else {
		return errors.New("send cmd timeout")
	}
}
