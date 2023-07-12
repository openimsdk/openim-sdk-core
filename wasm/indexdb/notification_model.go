//go:build js && wasm
// +build js,wasm

package indexdb

import (
	"context"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

type NotificationSeqs struct {
}

func NewNotificationSeqs() *NotificationSeqs {
	return &NotificationSeqs{}
}

func (i *NotificationSeqs) SetNotificationSeq(ctx context.Context, conversationID string, seq int64) error {
	_, err := Exec(conversationID, seq)
	return err
}

func (i *NotificationSeqs) GetNotificationAllSeqs(ctx context.Context) (result []*model_struct.NotificationSeqs, err error) {
	gList, err := Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}
