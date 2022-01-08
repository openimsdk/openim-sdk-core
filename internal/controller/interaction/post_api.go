package ws

import (
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/network"
	"open_im_sdk/pkg/server_api_params"
)

type PostApi struct {
	token      string
	apiAddress string
}

func (p *PostApi) PostFatalCallback(url string, data interface{}, callback common.Base, operationID string) *server_api_params.CommDataResp {
	content, err := network.Post2Api(p.apiAddress+url, data, p.token)
	c := common.CheckErrAndResp(callback, err, content, operationID)
	return c
}
