package testcore

import (
	"open_im_sdk/sdk_struct"
	"sync"
)

// config here
var AdminToken = ""

var Config = sdk_struct.IMConfig{
	ApiAddr:             APIADDR,
	WsAddr:              WSADDR,
	PlatformID:          PlatformID,
	DataDir:             DataDir,
	LogLevel:            LogLevel,
	IsLogStandardOutput: true,
}

var (
	rotateCount         = uint(0)
	LogLevel            = uint32(6)
	PlatformID          = int32(1)
	Secret              = "tuoyun"
	IsLogStandardOutput = true
	isLogJson           = false
	LogName             = ""
	LogFilePath         = ""
	DataDir             = "./../"
)

// system
var (
	// TESTIP       = "59.36.173.89"
	TESTIP       = "203.56.175.233"
	APIADDR      = "http://" + TESTIP + ":10002"
	WSADDR       = "ws://" + TESTIP + ":10001"
	REGISTERADDR = APIADDR + "/auth/user_register"
	TOKENADDR    = APIADDR + "/auth/user_token"
)

var userLock sync.RWMutex
var AllUserID []string
