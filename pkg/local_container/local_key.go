package local_container

import (
	"encoding/json"
	"errors"
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/utils"
	"github.com/mitchellh/mapstructure"
	constant2 "github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/log"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/network"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/server_api_params"
)

func SdkGetKey(ApiAddr, UserID, sessionID string, sType int32) (string, error) {
	if sessionID == "" {
		return "", utils.Wrap(errors.New("sessionID is nil"), "post failed ")
	}
	token := ""
	req := server_api_params.GetLocalKeyReq{
		SessionType: sType,
		UserID:      UserID,
		OperationID: utils.OperationIDGenerator(),
	}
	switch sType {
	case constant.SingleChatType:
		req.FriendID = sessionID
	case constant.GroupChatType:
		req.GroupID = sessionID
	}
	resp := server_api_params.GetLocalKeyResp{}
	content, err := network.Post2Api(ApiAddr+constant2.SDKGetKey, req, token)
	//content, err := network.Post2Api("http://116.62.111.115:10002/msg/sdk_get_key", req, token)
	if err != nil {
		log.Error("sdkGetKey err***************** ", utils.StructToJsonString(req))
		return "", utils.Wrap(err, "post failed ")
	}
	err = checkErrAndResp(err, content, &resp)
	if err != nil {
		log.Error("sdkGetKey err******************* ", utils.StructToJsonString(req))
		return "", utils.Wrap(err, "CheckErrAndResp failed ")
	}
	return resp.SessionInfo, nil
}
func SdkGetAllKey(ApiAddr, loginUserId string) (*server_api_params.GetAllLocalKeyBySessionIDResp, error) {
	req := server_api_params.GetAllLocalKeyBySessionIDReq{
		UserID:      loginUserId,
		OperationID: utils.OperationIDGenerator(),
	}
	resp := server_api_params.GetAllLocalKeyBySessionIDResp{}
	content, err := network.Post2Api(ApiAddr+constant2.GetAllKey, req, "p.token")
	//content, err := network.Post2Api("http://116.62.111.115:10002/msg/sdk_get_all_key", req, "p.token")
	if err != nil {
		utils.Wrap(err, "post failed "+string(content))
	}
	err = checkErrAndResp(err, content, &resp)
	return &resp, utils.Wrap(err, "CheckErrAndResp failed ")
}
func checkErrAndResp(err error, resp []byte, output interface{}) error {
	if err != nil {
		return utils.Wrap(err, "api resp failed")
	}
	var c server_api_params.CommDataResp
	err = json.Unmarshal(resp, &c)
	if err == nil {
		if c.ErrCode != 0 {
			return utils.Wrap(errors.New(c.ErrMsg), "")
		}
		if output != nil {
			err = mapstructure.Decode(c.Data, output)
			if err != nil {
				//	log.Error("mapstructure.Decode failed ", "err: ", err.Error(), c.Data)
				goto one
			}
			return nil
		}
		return nil
	} else {
		//	log.Error("json.Unmarshal failed ", string(resp), "err: ", err.Error())
	}

one:
	var c2 server_api_params.CommDataRespOne

	err = json.Unmarshal(resp, &c2)
	if err != nil {
		log.Error("json.Unmarshal failed ", string(resp), "err: ", err.Error())
		return utils.Wrap(err, "")
	}
	if c2.ErrCode != 0 {
		return utils.Wrap(errors.New(c2.ErrMsg), "")
	}
	if output != nil {
		err = mapstructure.Decode(c2.Data, output)
		if err != nil {
			return utils.Wrap(err, "")
		}
		return nil
	}
	return nil
}
