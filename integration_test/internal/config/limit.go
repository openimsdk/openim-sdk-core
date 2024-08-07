package config

const (
	MaxUserNum = 1e+5 // max users
)

const (
	ErrGroupSmallLimit  = 5   // max goroutine of a small error group
	ErrGroupCommonLimit = 150 // max goroutine of a common error group
)

const (
	SleepSec        = 30
	CheckWaitSec    = 5
	ProgressWaitSec = 1
)
