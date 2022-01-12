package interaction

import (
	"errors"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/network"
	"open_im_sdk/pkg/server_api_params"
	"sync"
)

//no share
type PostApi struct {
	token      string
	apiAddress string
}

func NewPostApi(token string, apiAddress string) *PostApi {
	return &PostApi{token: token, apiAddress: apiAddress}
}


func (p *PostApi) PostFatalCallback(callback common.Base, url string, data interface{}, operationID string) *server_api_params.CommDataResp {
	content, err := network.Post2Api(p.apiAddress+url, data, p.token)
	c := common.CheckErrAndResp(callback, err, content, operationID)
	return c
}

func (p *PostApi) OnError(errCode int32, errMsg string){
	p.err = errors.New(errMsg)
}

func (p *PostApi) OnSuccess(data string){
}



func (p *PostApi) PostReturn(url string, data interface{}, operationID string) (*server_api_params.CommDataResp, error) {
	var wg sync.WaitGroup
	wg.Add(1)
	var commData *server_api_params.CommDataResp
	go func() {
		commData = p.PostFatalCallback(p, url , data , operationID)
		wg.Done()
	}()

	wg.Wait()
	return commData, p.err
}


