//go:build js && wasm
// +build js,wasm

package indexdb

import (
	"context"
	"open_im_sdk/pkg/db/model_struct"
)

func (i IndexDB) SetNotificationSeq(ctx context.Context, conversationID string, seq int64) error {
	//TODO implement me
	panic("implement me")
}

func (i IndexDB) GetNotificationAllSeqs(ctx context.Context) ([]*model_struct.NotificationSeqsModel, error) {
	//TODO implement me
	panic("implement me")
}
