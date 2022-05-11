package test

import (
	"open_im_sdk/internal/interaction"
	"sync"
)

var (
	TESTIP       = "43.128.5.63"
	APIADDR      = "http://" + TESTIP + ":10000"
	WSADDR       = "ws://" + TESTIP + ":17778"
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
