package open_im_sdk

func triggerCmdFriend() error {
	return nil

}
func triggerCmdReLogin() error {
	return sendCmd(SdkInitManager.ch, cmd2Value{Cmd: CmdReLogin}, 2)
}

func triggerCmdBlackList() error {
	return nil

}

func triggerCmdFriendApplication() error {
	return nil

}

type deleteConNode struct {
	SourceID       string
	ConversationID string
}

func triggerCmdDeleteConversationAndMessage(sourceID, conversationID string) error {
	c2v := cmd2Value{
		Cmd:   CmdDeleteConversation,
		Value: deleteConNode{SourceID: sourceID, ConversationID: conversationID},
	}

	return sendCmd(ConversationCh, c2v, 1)
}

func triggerCmdGetLoginUserInfo() error {
	c2v := cmd2Value{
		Cmd: CmdGeyLoginUserInfo,
	}
	return sendCmd(InitCh, c2v, 1)
}

type updateConNode struct {
	ConId  string
	Action int //1 Delete the conversation; 2 Update the latest news in the conversation or add a conversation; 3 Put a conversation on the top;
	// 4 Cancel a conversation on the top, 5 Messages are not read and set to 0, 6 New conversations
	Args interface{}
}

func triggerCmdNewMsgCome(msg ArrMsg) error {
	c2v := cmd2Value{
		Cmd:   CmdNewMsgCome,
		Value: msg,
	}
	return sendCmd(ConversationCh, c2v, 1)
}

func triggerCmdAcceptFriend(sendUid string) error {
	return nil

}

func triggerCmdRefuseFriend(receiveUid string) error {
	return nil
}

func triggerCmdUpdateConversation(node updateConNode) error {
	c2v := cmd2Value{
		Cmd:   CmdUpdateConversation,
		Value: node,
	}

	return sendCmd(ConversationCh, c2v, 1)
}

func unInitAll() {
	c2v := cmd2Value{Cmd: CmdUnInit}
	_ = sendCmd(InitCh, c2v, 1)
	_ = sendCmd(ConversationCh, c2v, 1)
}



type goroutine interface {
	work(cmd cmd2Value)
	getCh() chan cmd2Value
}

func doListener(Li goroutine) {
	for {
		select {
		case cmd := <-Li.getCh():
			if cmd.Cmd == CmdUnInit {
				return
			}
			Li.work(cmd)
		}
	}
}
