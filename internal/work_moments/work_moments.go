// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package workMoments

//
//import (
//	"context"
//	"open_im_sdk/open_im_sdk_callback"
//	"open_im_sdk/pkg/db/db_interface"
//	"open_im_sdk/pkg/db/model_struct"
//	"open_im_sdk/pkg/utils"
//
//	"github.com/OpenIMSDK/tools/log"
//)
//
//type WorkMoments struct {
//	listener    open_im_sdk_callback.OnWorkMomentsListener
//	loginUserID string
//	db          db_interface.DataBase
//}
//
//funcation NewWorkMoments(loginUserID string, db db_interface.DataBase) *WorkMoments {
//	return &WorkMoments{loginUserID: loginUserID, db: db}
//}
//
//funcation (w *WorkMoments) DoNotification(ctx context.Context, jsonDetail string) {
//	if w.listener == nil {
//		return
//	}
//	if err := w.db.InsertWorkMomentsNotification(ctx, jsonDetail); err != nil {
//		log.ZError(ctx, "InsertWorkMomentsNotification failed", err, "jsonDetail", jsonDetail)
//		return
//	}
//	if err := w.db.IncrWorkMomentsNotificationUnreadCount(ctx); err != nil {
//		log.ZError(ctx, "IncrWorkMomentsNotificationUnreadCount failed", err)
//		return
//	}
//	w.listener.OnRecvNewNotification()
//}
//
//funcation (w *WorkMoments) getWorkMomentsNotification(ctx context.Context, offset, count int) ([]*model_struct.WorkMomentNotificationMsg, error) {
//	if err := w.db.MarkAllWorkMomentsNotificationAsRead(ctx); err != nil {
//		return nil, err
//	}
//	workMomentsNotifications, err := w.db.GetWorkMomentsNotification(ctx, offset, count)
//	if err != nil {
//		return nil, err
//	}
//	msgs := make([]*model_struct.WorkMomentNotificationMsg, len(workMomentsNotifications))
//	for i, v := range workMomentsNotifications {
//		workMomentNotificationMsg := model_struct.WorkMomentNotificationMsg{}
//		if err := utils.JsonStringToStruct(v.JsonDetail, &workMomentNotificationMsg); err != nil {
//			log.ZError(ctx, "invalid data", err, "jsonDetail", v.JsonDetail)
//			continue
//		}
//		msgs[i] = &workMomentNotificationMsg
//	}
//	return msgs, nil
//}
//
//funcation (w *WorkMoments) clearWorkMomentsNotification(ctx context.Context) error {
//	return w.db.ClearWorkMomentsNotification(ctx)
//}
//
//funcation (w *WorkMoments) getWorkMomentsNotificationUnReadCount(ctx context.Context) (model_struct.LocalWorkMomentsNotificationUnreadCount, error) {
//	return w.db.GetWorkMomentsUnReadCount(ctx)
//}
