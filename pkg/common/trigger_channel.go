package common

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"
	"time"
	"errors"
)

func TriggerCmdNewMsgCome(msg utils.ArrMsg, conversationCh chan Cmd2Value) error {
	c2v := Cmd2Value{Cmd:   constant.CmdNewMsgCome, Value: msg}
	return sendCmd(conversationCh, c2v, 1)
}



type Cmd2Value struct {
	Cmd   string
	Value interface{}
}



func sendCmd(ch chan Cmd2Value, value Cmd2Value, timeout int64) error {
	var flag = 0
	select {
	case ch <- value:
		flag = 1
	case <-time.After(time.Second * time.Duration(timeout)):
		flag = 2
	}
	if flag == 1 {
		return nil
	} else {
		return errors.New("send cmd timeout")
	}
}