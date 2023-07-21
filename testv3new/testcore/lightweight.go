package testcore

import (
	"open_im_sdk/internal/interaction"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

type BaseCore struct {
	longConnMgr *interaction.LongConnMgr
	userID      string
	recvMsgNum  int
}

func NewBaseCore(userID string) *BaseCore {
	return &BaseCore{
		userID: userID,
	}
}

func (b *BaseCore) SetCallback() {

}

func (b *BaseCore) InitConn() error {
	return nil
}

func (b *BaseCore) SendMsg(index int) error {
	return nil
}

func (b *BaseCore) recvMsgCallback(msg *sdkws.MsgData) {
	b.recvMsgNum++
}

func (b *BaseCore) GetRecvMsgNum() int {
	return b.recvMsgNum
}
