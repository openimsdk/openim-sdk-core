package config

import (
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/constant"
)

const (
	TestIP              = "127.0.0.1"
	APIAddr             = "http://" + TestIP + ":10002"
	WsAddr              = "ws://" + TestIP + ":10001"
	AdminUserID         = "imAdmin"
	Secret              = "openIM123"
	PlatformID          = constant.AndroidPlatformID
	LogLevel            = 3
	DataDir             = "./data/"
	LogFilePath         = "./logs/"
	IsLogStandardOutput = false
)

func GetConf() sdk_struct.IMConfig {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = APIAddr
	cf.WsAddr = WsAddr
	cf.DataDir = DataDir
	cf.LogLevel = LogLevel
	cf.IsExternalExtensions = true
	cf.PlatformID = int32(PlatformID)
	cf.LogFilePath = LogFilePath
	cf.IsLogStandardOutput = IsLogStandardOutput
	return cf
}
