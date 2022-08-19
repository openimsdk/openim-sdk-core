package test

import (
	"sync"
)

var LogLevel = 3
var PlatformID = int32(1)

var (
	TESTIP = "43.128.5.63"
	//TESTIP       = "121.37.25.71"
	APIADDR         = "http://" + TESTIP + ":10002"
	WSADDR          = "ws://" + TESTIP + ":10001"
	REGISTERADDR    = APIADDR + "/auth/user_register"
	TOKENADDR       = APIADDR + "/auth/user_token"
	SECRET          = "tuoyun"
	SENDINTERVAL    = 20
	GETSELFUSERINFO = APIADDR + "/user/get_self_user_info"
)

var coreMgrLock sync.RWMutex
var allLoginMgr map[int]*CoreNode

var userLock sync.RWMutex

var allUserID []string
var allToken []string

//var allWs []*interaction.Ws
var sendSuccessCount, sendFailedCount int
var sendSuccessLock sync.RWMutex
var sendFailedLock sync.RWMutex

var msgNumInOneClient = 0

//var Msgwg sync.WaitGroup
var sendMsgClient = 0

var MaxNumGoroutine = 100000
