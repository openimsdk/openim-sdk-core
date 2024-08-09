package config

const (
	MaxUserNum = 1e+5 // max users
)

const (
	ErrGroupSmallLimit       = 1   // max goroutine of a small error group
	ErrGroupMiddleSmallLimit = 5   // max goroutine of a middle small error group
	ErrGroupCommonLimit      = 100 // max goroutine of a common error group
)

const (
	SleepSec        = 30
	CheckWaitSec    = 5 // check wait sec
	ProgressWaitSec = 1 // progress bar update wait sec
)

const (
	CheckMsgRate = 0.1 // Sampling and statistical message ratio. Max check message is MaxCheckMsg
	MaxCheckMsg  = 1e+5
)
