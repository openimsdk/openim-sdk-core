package testv3new

import (
	"fmt"
	"open_im_sdk/pkg/constant"
)

// config here

// system
var (
	TESTIP     = "59.36.173.89"
	APIADDR    = fmt.Sprintf("http://%v:10002", TESTIP)
	WSADDR     = fmt.Sprintf("ws://%v:10001", TESTIP)
	SECRET     = "openIM123"
	PLATFORMID = constant.WindowsPlatformID

	REGISTERADDR = APIADDR + constant.UserRegister
	TOKENADDR    = APIADDR + constant.GetUsersToken
)
