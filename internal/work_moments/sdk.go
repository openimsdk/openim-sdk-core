package workMoments

import (
	"context"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/db/model_struct"
)

func (w *WorkMoments) SetListener(callback open_im_sdk_callback.OnWorkMomentsListener) {
	if callback == nil {
		return
	}
	w.listener = callback
}

func (w *WorkMoments) GetWorkMomentsUnReadCount(ctx context.Context) (model_struct.LocalWorkMomentsNotificationUnreadCount, error) {
	return w.getWorkMomentsNotificationUnReadCount(ctx)
}

func (w *WorkMoments) GetWorkMomentsNotification(ctx context.Context, offset, count int) ([]*model_struct.WorkMomentNotificationMsg, error) {
	return w.getWorkMomentsNotification(ctx, offset, count)
}

func (w *WorkMoments) ClearWorkMomentsNotification(ctx context.Context) error {
	return w.clearWorkMomentsNotification(ctx)
}
