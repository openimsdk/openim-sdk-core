package third

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/protocol/third"
)

func (c *Third) UpdateFcmToken(ctx context.Context, fcmToken string, expireTime int64) error {
	req := third.FcmUpdateTokenReq{
		PlatformID: c.platformID,
		FcmToken:   fcmToken,
		Account:    c.loginUserID,
		ExpireTime: expireTime}
	_, err := util.CallApi[third.FcmUpdateTokenResp](ctx, constant.FcmUpdateTokenRouter, &req)
	return err

}

func (c *Third) SetAppBadge(ctx context.Context, appUnreadCount int32) error {
	req := third.SetAppBadgeReq{
		UserID:         c.loginUserID,
		AppUnreadCount: appUnreadCount,
	}
	_, err := util.CallApi[third.SetAppBadgeResp](ctx, constant.SetAppBadgeRouter, &req)
	return err
}

func (c *Third) UploadLogs(ctx context.Context, line int, ex string, progress Progress) (err error) {
	return c.uploadLogs(ctx, line, ex, progress)
}

func (c *Third) Log(ctx context.Context, logLevel int, relativePath string, line string, msg, err string, keysAndValues []any) {
	c.printLog(ctx, logLevel, relativePath, line, msg, err, keysAndValues)
}
