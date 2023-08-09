package testv3new

import (
	"context"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/sdk_struct"
)

func NewUserCtx(userID, token string) context.Context {
	return ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID: userID,
		Token:  token,
		IMConfig: sdk_struct.IMConfig{
			PlatformID: int32(PLATFORMID),
			ApiAddr:    APIADDR,
			WsAddr:     WSADDR,
		}})
}
