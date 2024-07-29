package manager

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/config"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	authPB "github.com/openimsdk/protocol/auth"
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
	return util.ApiPost(ctx, route, req, resp)
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
	return ctx
}

func (m *MetaManager) GetSecret() string {
	return m.secret
}

func (m *MetaManager) GetToken(userID string, platformID int32) (string, error) {
	req := authPB.UserTokenReq{PlatformID: platformID, UserID: userID, Secret: m.secret}
	resp := authPB.UserTokenResp{}
	err := m.PostWithCtx(constant.GetUsersToken, &req, &resp)
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (m *MetaManager) WithAdminToken() (err error) {
	token, err := m.GetToken(config.AdminUserID, config.PlatformID)
	if err != nil {
		return err
	}
	m.token = token
	return nil
}
