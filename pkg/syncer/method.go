package syncer

type Method uint8

func (s Method) String() string {
	switch s {
	case MethodChange:
		return "Change"
	case MethodDelete:
		return "Delete"
	case MethodComplete:
		return "complete"
	default:
		return "Unknown"
	}
}

const (
	MethodChange   = 1
	MethodDelete   = 2
	MethodComplete = 3
)
