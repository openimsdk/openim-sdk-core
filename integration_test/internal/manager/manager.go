package manager

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/network"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	authPB "github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/tools/mcontext"
)

type MetaManager struct {
	sdk_struct.IMConfig
	secret string
	token  string
}

func NewMetaManager() *MetaManager {
	conf := config.GetConf()
	return &MetaManager{
		IMConfig: conf,
		secret:   config.Secret,
	}
}

func (m *MetaManager) ApiPost(ctx context.Context, route string, req, resp any) (err error) {
	return network.ApiPost(ctx, route, req, resp)
}

// PostWithCtx should only be used for scenarios such as registration and login that do not require a token or
// require an admin token.
// For scenarios that require a specific user token, please obtain the context from vars.Contexts for the request.
func (m *MetaManager) PostWithCtx(route string, req, resp any) error {
	return m.ApiPost(m.BuildCtx(nil), route, req, resp)
}

// BuildCtx build an admin token
func (m *MetaManager) BuildCtx(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = ccontext.WithInfo(ctx, &ccontext.GlobalConfig{
		Token:    m.token,
		IMConfig: m.IMConfig,
	})
	ctx = ccontext.WithOperationID(ctx, utils.OperationIDGenerator())
	ctx = mcontext.SetOpUserID(ctx, "admin")
	return ctx
}

func (m *MetaManager) GetSecret() string {
	return m.secret
}

func (m *MetaManager) GetAdminToken(userID string, platformID int32) (string, error) {
	req := authPB.GetAdminTokenReq{UserID: userID, Secret: m.secret}
	resp := authPB.GetAdminTokenResp{}
	err := m.PostWithCtx(api.GetAdminToken.Route(), &req, &resp)
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (m *MetaManager) GetUserToken(userID string, platformID int32) (string, error) {
	req := authPB.GetUserTokenReq{PlatformID: platformID, UserID: userID}
	resp := authPB.GetUserTokenResp{}
	err := m.PostWithCtx(api.GetUsersToken.Route(), &req, &resp)
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (m *MetaManager) WithAdminToken() (err error) {
	token, err := m.GetAdminToken(config.AdminUserID, config.PlatformID)
	if err != nil {
		return err
	}
	m.token = token
	return nil
}
