package vars

import (
	"context"
	"sync/atomic"
)

const (
	UserIDPrefix = "test_v3_u"
)

var (
	UserIDs      []string // all user ids
	SuperUserIDs []string // user ids of users with all friends

	Contexts []context.Context // users contexts

	LoginEndUserNum int // e.g. if LoginEndUserNum = 5, login user is [0,1,2,3,4]
	NowLoginNum     atomic.Int64
)
