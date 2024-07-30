package config

const (
	MaxUserNum = 1e+5 // max users
)

const (
	ErrGroupSmallLimit  = 5   // max goroutine of a small error group
	ErrGroupCommonLimit = 200 // max goroutine of a common error group
)
