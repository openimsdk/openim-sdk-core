package interaction

import (
	"errors"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/network"
	"open_im_sdk/pkg/server_api_params"
)

//no share
type PostApi struct {
	token      string
	apiAddress string
}

func NewPostApi(token string, apiAddress string) *PostApi {
	return &PostApi{token: token, apiAddress: apiAddress}
}

func (p *PostApi) PostFatalCallback(callback common.Base, url string, data interface{}, output interface{}, operationID string) *server_api_params.CommDataResp {
	content, err := network.Post2Api(p.apiAddress+url, data, p.token)
	common.CheckErrAndRespCallback(callback, err, content, output, operationID)
}

func (pe *postErr) OnError(errCode int32, errMsg string) {
	pe.err = errors.New(errMsg)
}

func (pe *postErr) OnSuccess(data string) {
}

type postErr struct {
	err error
}

func (p *PostApi) PostReturn(url string, data interface{}, output interface{}) error {
	//log.Debug("000", utils.GetSelfFuncName(), p.apiAddress+url)
	content, err := network.Post2Api(p.apiAddress+url, data, p.token)
	return common.CheckErrAndResp(err, content, output)
}
