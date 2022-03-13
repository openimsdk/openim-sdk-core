package rtc

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/log"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
)

func (s *LiveSignaling) Invite(signalInviteReq string, callback open_im_sdk_callback.Base, operationID string){
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", signalInviteReq)
		var unmarshalReq api.SignalInviteReq
		common.JsonUnmarshalCallback(signalInviteReq, &unmarshalReq, callback, operationID)
		result := s.invite(&unmarshalReq, callback, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonStringDefault(result))
	}()
}
