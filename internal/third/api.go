package third

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"

	"github.com/openimsdk/protocol/third"
)

func (c *Third) UpdateFcmToken(ctx context.Context, fcmToken string, expireTime int64) error {
	_, err := api.FcmUpdateToken.Invoke(ctx, &third.FcmUpdateTokenReq{
		PlatformID: c.platformID,
		FcmToken:   fcmToken,
		Account:    c.loginUserID,
		ExpireTime: expireTime})
	return err
}

func (c *Third) SetAppBadge(ctx context.Context, appUnreadCount int32) error {
	_, err := api.SetAppBadge.Invoke(ctx, &third.SetAppBadgeReq{
		UserID:         c.loginUserID,
		AppUnreadCount: appUnreadCount,
	})
	return err
}

func (c *Third) UploadLogs(ctx context.Context, line int, ex string, progress Progress) (err error) {
	return c.uploadLogs(ctx, line, ex, progress)
}

func (c *Third) Log(ctx context.Context, logLevel int, file string, line int, msg, err string, keysAndValues []any) {
	c.printLog(ctx, logLevel, file, line, msg, err, keysAndValues)
}
