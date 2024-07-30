package vars

import "context"

const (
	UserIDPrefix = "test_v3_u"
)

var (
	UserIDs      []string // all user ids
	SuperUserIDs []string // user ids of users with all friends

	Contexts []context.Context // users contexts
)
