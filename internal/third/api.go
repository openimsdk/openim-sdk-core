package third

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"

	"github.com/openimsdk/protocol/third"
)

func (c *Third) UpdateFcmToken(ctx context.Context, fcmToken string, expireTime int64) error {
	return api.FcmUpdateToken.Execute(ctx, &third.FcmUpdateTokenReq{
		PlatformID: c.platformID,
		FcmToken:   fcmToken,
		Account:    c.loginUserID,
		ExpireTime: expireTime,
	})
}

func (c *Third) SetAppBadge(ctx context.Context, appUnreadCount int32) error {
	return api.SetAppBadge.Execute(ctx, &third.SetAppBadgeReq{
		UserID:         c.loginUserID,
		AppUnreadCount: appUnreadCount,
	})
}

func (c *Third) UploadLogs(ctx context.Context, line int, ex string, progress Progress) (err error) {
	return c.uploadLogs(ctx, line, ex, progress)
}

func (c *Third) Log(ctx context.Context, logLevel int, file string, line int, msg, err string, keysAndValues []any) {
	c.printLog(ctx, logLevel, file, line, msg, err, keysAndValues)
}
