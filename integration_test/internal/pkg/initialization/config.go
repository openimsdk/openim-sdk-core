package initialization

import (
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
)

func GetConf() sdk_struct.IMConfig {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = vars.APIAddr
	cf.WsAddr = vars.WsAddr
	cf.DataDir = vars.DataDir
	cf.LogLevel = vars.LogLevel
	cf.IsExternalExtensions = true
	cf.PlatformID = int32(vars.PlatformID)
	cf.LogFilePath = vars.LogFilePath
	cf.IsLogStandardOutput = true
	return cf
}
