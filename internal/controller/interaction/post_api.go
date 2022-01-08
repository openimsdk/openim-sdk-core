package interaction

import (
	"open_im_sdk/pkg/commom"
	"open_im_sdk/pkg/network"
	"open_im_sdk/pkg/server_api_params"
)

type PostApi struct {
	token string
}

func (p *PostApi) PostFatalCallback(url string, data interface{}, callback commom.Base, operationID string) *server_api_params.CommDataResp {
	content, err := network.Post2Api(url, data, p.token)
	c := commom.CheckErrAndResp(callback, err, content, operationID)
	return c
}
