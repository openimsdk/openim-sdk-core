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

func (pe *postErr) OnError(errCode int32, errMsg string){
	pe.err = errors.New(errMsg)
}

func (pe *postErr) OnSuccess(data string){
}

type postErr struct {
	err error
}


func (p *PostApi) PostReturn(url string, data interface{}, operationID string) (*server_api_params.CommDataResp, error) {
	 pe := postErr{}
	var wg sync.WaitGroup
	wg.Add(1)
	var commData *server_api_params.CommDataResp
	go func() {
		commData = p.PostFatalCallback(&pe, url , data , operationID)
		wg.Done()
	}()

	wg.Wait()
	return commData, pe.err
}


