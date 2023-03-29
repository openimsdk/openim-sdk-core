package workMoments

import (
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/utils"
)

type WorkMoments struct {
	listener    open_im_sdk_callback.OnWorkMomentsListener
	loginUserID string
	db          db_interface.DataBase
	p           *ws.PostApi
}

func NewWorkMoments(loginUserID string, db db_interface.DataBase, p *ws.PostApi) *WorkMoments {
	return &WorkMoments{loginUserID: loginUserID, db: db, p: p}
}

func (w *WorkMoments) DoNotification(jsonDetailStr string, operationID string) {
	if w.listener == nil {
		log.NewDebug(operationID, "WorkMoments listener is null", jsonDetailStr)
		return
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "json_detail: ", jsonDetailStr)
	if err := w.db.InsertWorkMomentsNotification(jsonDetailStr); err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "InsertWorkMomentsNotification failed", err.Error())
		return
	}
	if err := w.db.IncrWorkMomentsNotificationUnreadCount(); err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "IncrWorkMomentsNotificationUnreadCount failed", err.Error())
		return
	}
	w.listener.OnRecvNewNotification()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "do notification callback success")
}

func (w *WorkMoments) getWorkMomentsNotification(offset, count int, callback open_im_sdk_callback.Base, operationID string) sdk_params_callback.GetWorkMomentsNotificationCallback {
	log.NewInfo(operationID, utils.GetSelfFuncName(), offset, count)
	err := w.db.MarkAllWorkMomentsNotificationAsRead()
	common.CheckDBErrCallback(callback, err, operationID)
	workMomentsNotifications, err := w.db.GetWorkMomentsNotification(offset, count)
	common.CheckDBErrCallback(callback, err, operationID)
	msgs := make([]sdk_params_callback.WorkMomentNotificationMsg, len(workMomentsNotifications))
	for i, v := range workMomentsNotifications {
		workMomentNotificationMsg := sdk_params_callback.WorkMomentNotificationMsg{}
		if err := utils.JsonStringToStruct(v.JsonDetail, &workMomentNotificationMsg); err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "JsonStringToStruct failed", err.Error())
			continue
		}
		msgs[i] = workMomentNotificationMsg
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "success")
	return msgs
}

func (w *WorkMoments) clearWorkMomentsNotification(callback open_im_sdk_callback.Base, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	err := w.db.ClearWorkMomentsNotification()
	common.CheckDBErrCallback(callback, err, operationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "success")
}

func (w *WorkMoments) getWorkMomentsNotificationUnReadCount(callback open_im_sdk_callback.Base, operationID string) sdk_params_callback.GetWorkMomentsUnReadCountCallback {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	unreadCount, err := w.db.GetWorkMomentsUnReadCount()
	common.CheckDBErrCallback(callback, err, operationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "success")
	return sdk_params_callback.GetWorkMomentsUnReadCountCallback(unreadCount)
}
