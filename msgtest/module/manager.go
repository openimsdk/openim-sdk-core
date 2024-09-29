package module

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/network"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	authPB "github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
)

const (
	ManagerUserID = "openIMAdmin"
)

type MetaManager struct {
	managerUserID string
	apiAddr       string
	secret        string
	token         string
}

func NewMetaManager(apiAddr, secret, managerUserID string) *MetaManager {
	return &MetaManager{
		managerUserID: managerUserID,
		apiAddr:       apiAddr,
		secret:        secret,
	}
}

func (m *MetaManager) NewUserManager() *TestUserManager {
	return &TestUserManager{m}
}

func (m *MetaManager) NewGroupMananger() *TestGroupManager {
	return &TestGroupManager{m}
}

func (m *MetaManager) NewFriendManager() *TestFriendManager {
	return &TestFriendManager{m}
}

func (m *MetaManager) NewApiMsgSender() *ApiMsgSender {
	return &ApiMsgSender{m}
}

func (m *MetaManager) apiPost(ctx context.Context, route string, req, resp any) (err error) {
	operationID, _ := ctx.Value("operationID").(string)
	if operationID == "" {
		return errs.ErrArgs.WrapMsg("call api operationID is empty")
	}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return errs.ErrArgs.WrapMsg("json.Marshal(req) failed " + err.Error())
	}
	reqUrl := m.apiAddr + route
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl, bytes.NewReader(reqBody))
	if err != nil {
		return errs.ErrArgs.WrapMsg("sdk http.NewRequestWithContext failed " + err.Error())
	}
	start := time.Now()
	log.ZDebug(ctx, "ApiRequest", "url", reqUrl, "body", string(reqBody))
	request.ContentLength = int64(len(reqBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("operationID", operationID)
	if m.token != "" {
		request.Header.Set("token", m.token)
	}
	response, err := new(http.Client).Do(request)
	if err != nil {
		return errs.ErrArgs.WrapMsg("ApiPost http.Client.Do failed " + err.Error())
	}
	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.ZError(ctx, "ApiResponse", err, "type", "read body", "status", response.Status)
		return errs.ErrArgs.WrapMsg("io.ReadAll(ApiResponse) failed " + err.Error())
	}
	log.ZDebug(ctx, "ApiResponse", "url", reqUrl, "status", response.Status,
		"body", string(respBody), "time", time.Since(start).Milliseconds())
	var baseApi network.ApiResponse
	if err := json.Unmarshal(respBody, &baseApi); err != nil {
		return sdkerrs.ErrSdkInternal.WrapMsg(fmt.Sprintf("api %s json.Unmarshal(%q, %T) failed %s", m.apiAddr, string(respBody), &baseApi, err.Error()))
	}
	if baseApi.ErrCode != 0 {
		err := sdkerrs.New(baseApi.ErrCode, baseApi.ErrMsg, baseApi.ErrDlt)
		return err
	}
	if resp == nil || len(baseApi.Data) == 0 || string(baseApi.Data) == "null" {
		return nil
	}
	if err := json.Unmarshal(baseApi.Data, resp); err != nil {
		return sdkerrs.ErrSdkInternal.WrapMsg(fmt.Sprintf("json.Unmarshal(%q, %T) failed %s", string(baseApi.Data), resp, err.Error()))
	}
	return nil
}

func (m *MetaManager) postWithCtx(route string, req, resp any) error {
	return m.apiPost(m.buildCtx(), route, req, resp)
}

func (m *MetaManager) buildCtx() context.Context {
	return mcontext.NewCtx(utils.OperationIDGenerator())
}

func (m *MetaManager) getAdminToken(userID string) (string, error) {
	req := authPB.GetAdminTokenReq{UserID: userID, Secret: m.secret}
	resp := authPB.GetAdminTokenResp{}
	err := m.postWithCtx(api.GetAdminToken.Route(), &req, &resp)
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (m *MetaManager) getUserToken(userID string, platform int32) (string, error) {
	req := authPB.GetUserTokenReq{UserID: userID, PlatformID: platform}
	resp := authPB.GetUserTokenResp{}
	err := m.postWithCtx(api.GetUsersToken.Route(), &req, &resp)
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (m *MetaManager) initToken() error {
	token, err := m.getAdminToken(m.managerUserID)
	if err != nil {
		return err
	}
	m.token = token
	return nil
}
func (m *MetaManager) GetServerTime() (int64, error) {
	req := msg.GetServerTimeReq{}
	resp := msg.GetServerTimeResp{}
	err := m.postWithCtx(api.GetServerTime.Route(), &req, &resp)
	if err != nil {
		return 0, err
	} else {
		return resp.ServerTime, nil
	}
}
