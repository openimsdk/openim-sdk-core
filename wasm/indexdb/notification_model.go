//go:build js && wasm
// +build js,wasm

package indexdb

import (
	"context"
	"open_im_sdk/pkg/db/model_struct"
)

type NotificationSeqs struct {
}

func NewNotificationSeqs() *NotificationSeqs {
	return &NotificationSeqs{}
}

func (i *NotificationSeqs) SetNotificationSeq(ctx context.Context, conversationID string, seq int64) error {
	//TODO implement me
	panic("implement me")
}

func (i *NotificationSeqs) GetNotificationAllSeqs(ctx context.Context) ([]*model_struct.NotificationSeqs, error) {
	//TODO implement me
	panic("implement me")
}
