package obj_storage

import (
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/open_im_sdk_callback"

	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/server_api_params"
)

type Push struct {
	p          *ws.PostApi
	platformID int32
}

func NewPush(p *ws.PostApi, platformID int32) *Push {
	return &Push{p: p, platformID: platformID}
}

func (c *Push) UpdateFcmToken(callback open_im_sdk_callback.Base, fcmToken, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "UpdateFcmToken args: ", fcmToken)
		c.fmcUpdateToken(callback, fcmToken, operationID)
		callback.OnSuccess(sdk_params_callback.UpdateFcmTokenCallback)
		log.NewInfo(operationID, "UpdateFcmToken callback: ", sdk_params_callback.UpdateFcmTokenCallback)
	}()

}

func (c *Push) fmcUpdateToken(callback open_im_sdk_callback.Base, fcmToken, operationID string) {
	apiReq := server_api_params.FcmUpdateTokenReq{}
	apiReq.OperationID = operationID
	apiReq.Platform = int(c.platformID)
	apiReq.FcmToken = fcmToken
	c.p.PostFatalCallback(callback, constant.FcmUpdateTokenRouter, apiReq, nil, apiReq.OperationID)
}
