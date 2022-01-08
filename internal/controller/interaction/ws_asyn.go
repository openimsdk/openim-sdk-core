package ws

import (
	"errors"
	"open_im_sdk/pkg/utils"
	"sync"
	"time"
)

type GeneralWsResp struct {
	ReqIdentifier int    `json:"reqIdentifier"`
	ErrCode       int    `json:"errCode"`
	ErrMsg        string `json:"errMsg"`
	MsgIncr       string `json:"msgIncr"`
	OperationID   string `json:"operationID"`
	Data          []byte `json:"data"`
}

type GeneralWsReq struct {
	ReqIdentifier int32  `json:"reqIdentifier"`
	Token         string `json:"token"`
	SendID        string `json:"sendID"`
	OperationID   string `json:"operationID"`
	MsgIncr       string `json:"msgIncr"`
	Data          []byte `json:"data"`
}

type WsRespAsyn struct {
	wsNotification map[string]chan GeneralWsResp
	wsMutex        sync.RWMutex
}

func GenMsgIncr(userID string) string {
	return userID + "_" + utils.Int64ToString(utils.GetCurrentTimestampByNano())
}

func (u *WsRespAsyn) AddCh(userID string) (string, chan GeneralWsResp) {
	u.wsMutex.Lock()
	defer u.wsMutex.Unlock()
	msgIncr := GenMsgIncr(userID)

	ch := make(chan GeneralWsResp, 1)
	_, ok := u.wsNotification[msgIncr]
	if ok {
	}
	u.wsNotification[msgIncr] = ch
	return msgIncr, ch
}

func (u *WsRespAsyn) GetCh(msgIncr string) chan GeneralWsResp {
	ch, ok := u.wsNotification[msgIncr]
	if ok {
		return ch
	}

	return nil
}

func (u *WsRespAsyn) DelCh(msgIncr string) {
	u.wsMutex.Lock()
	defer u.wsMutex.Unlock()
	ch, ok := u.wsNotification[msgIncr]
	if ok {
		close(ch)
		delete(u.wsNotification, msgIncr)
	}
}

func notifyCh(ch chan GeneralWsResp, value GeneralWsResp, timeout int64) error {
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
		//	sdkLog("send cmd timeout, ", timeout, value)
		return errors.New("send cmd timeout")
	}
}

func (u *WsRespAsyn) notifyResp(wsResp GeneralWsResp) {
	utils.LogBegin(wsResp.OperationID)
	u.wsMutex.Lock()
	defer u.wsMutex.Unlock()

	ch := u.GetCh(wsResp.MsgIncr)
	if ch == nil {
		return
	}

	err := notifyCh(ch, wsResp, 1)
	if err != nil {
		//utils.sdkLog("notifyCh failed, ", err.Error(), ch, wsResp)
	}
	//utils.sdkLog("notify ch end", wsResp.OperationID)
}
