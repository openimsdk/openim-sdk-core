package cliconf

import "context"

var globalClientConfig = &clientConfig{}

func SetLoginUserID(userID string) {
	globalClientConfig.userID = userID
	globalClientConfig.ClearConfig()
}

func ClearConfig() {
	globalClientConfig.ClearConfig()
}

func GetClientConfig(ctx context.Context) (*ClientConfig, error) {
	return globalClientConfig.GetConfig(ctx)
}
