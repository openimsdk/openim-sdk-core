package vars

import (
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
)

const (
	TestIP      = "125.124.195.201"
	APIAddr     = "http://" + TestIP + ":10002"
	WsAddr      = "ws://" + TestIP + ":10001"
	Secret      = "openIM123"
	PlatformID  = constant.WindowsPlatformID
	LogLevel    = uint32(5)
	DataDir     = "../"
	LogFilePath = "../"
)
