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

	LoginEndUserNum int
	NowLoginNum     atomic.Int64
)
