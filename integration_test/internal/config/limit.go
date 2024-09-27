package config

const (
	MaxUserNum = 1e+5 // max users
)

const (
	ErrGroupSmallLimit       = 20  // max goroutine of a small error group
	ErrGroupMiddleSmallLimit = 50  // max goroutine of a middle small error group
	ErrGroupCommonLimit      = 100 // max goroutine of a common error group
)

const (
	SleepSec          = 30
	CheckWaitSec      = 5 // check wait sec
	BarRemoveWaiteSec = 1 // progress bar remove wait second
)

const (
	CheckMsgRate    = 1 // Sampling and statistical message ratio. Max check message is MaxCheckMsg
	MaxCheckMsg     = 1e+8
	MaxCheckLoopNum = 40
)

const (
	ApiParamLength = 1000
)
