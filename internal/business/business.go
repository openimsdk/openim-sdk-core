package business

import (
	"context"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/db/db_interface"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
)

type Business struct {
	listener open_im_sdk_callback.OnCustomBusinessListener
	db       db_interface.DataBase
}

func NewBusiness(db db_interface.DataBase) *Business {
	return &Business{
		db: db,
	}
}

func (b *Business) DoNotification(ctx context.Context, jsonDetailStr string) {
	if b.listener == nil {
		log.ZWarn(ctx, "listener is nil", nil)
		return
	}
	b.listener.OnRecvCustomBusinessMessage(jsonDetailStr)
}
