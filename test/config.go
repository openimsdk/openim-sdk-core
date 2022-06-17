package test

import (
	"open_im_sdk/internal/interaction"
	"sync"
)

var (
	//	TESTIP = "cjyy.zs.sjz.gov.cn"
	TESTIP = "43.128.5.63"
	//TESTIP       = "8.213.195.63"
	//TESTIP       = "103.116.45.174"
	APIADDR      = "http://" + TESTIP + ":10002"
	WSADDR       = "ws://" + TESTIP + ":10001"
	REGISTERADDR = APIADDR + "/auth/user_register"
	TOKENADDR    = APIADDR + "/auth/user_token"
	SECRET       = "tuoyun"
	SENDINTERVAL = 20
	ACCOUNTCHECK = APIADDR + "/manager/account_check"
)

var coreMgrLock sync.RWMutex
var allLoginMgr map[int]*CoreNode

var userLock sync.RWMutex

var allUserID []string
var allToken []string
var allWs []*interaction.Ws
var sendSuccessCount, sendFailedCount int
var sendSuccessLock sync.RWMutex
var sendFailedLock sync.RWMutex

var msgNumInOneClient = 0

//var Msgwg sync.WaitGroup
var sendMsgClient = 0

var MaxNumGoroutine = 100000
