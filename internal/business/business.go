package business

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/db/db_interface"
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
