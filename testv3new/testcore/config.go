package testcore

import "open_im_sdk/pkg/constant"

// config here

// system
var (
	TESTIP  = "203.56.175.233"
	APIADDR = "http://" + TESTIP + ":10002"
	WSADDR  = "ws://" + TESTIP + ":10001"
	SECRET  = "tuoyun"

	REGISTERADDR = APIADDR + "/auth/user_register"
	TOKENADDR    = APIADDR + "/auth/user_token"
	PLATFORMID   = constant.AndroidPlatformID
)
