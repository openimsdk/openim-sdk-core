package testcore

import (
	"fmt"
	"open_im_sdk/pkg/constant"
)

// config here

// system
var (
	TESTIP  = "43.128.72.19"
	APIADDR = fmt.Sprintf("http://%v:10002", TESTIP)
	WSADDR  = fmt.Sprintf("ws://%v:10001", TESTIP)
	SECRET  = "openIM123"

	REGISTERADDR = APIADDR + constant.UserRegister
	TOKENADDR    = APIADDR + constant.GetUsersToken
	PLATFORMID   = constant.AndroidPlatformID
)
