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

const (
	ManagerUserID = "openIMAdmin"
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

func (m *MetaManager) PostWithCtx(route string, req, resp any) error {
	return m.ApiPost(m.BuildCtx(nil), route, req, resp)
}

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
