package checker

type CountChecker struct {
	TotalCount   int
	CorrectCount int
	IsEqual      bool
}

func NewCountChecker(total, correct int, isEqual bool) *CountChecker {
	return &CountChecker{
		TotalCount:   total,
		CorrectCount: correct,
		IsEqual:      isEqual,
	}
}
