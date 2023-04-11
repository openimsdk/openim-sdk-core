package syncer

type Method uint8

func (s Method) String() string {
	switch s {
	case MethodLocally:
		return "Locally"
	case MethodDelete:
		return "Delete"
	case MethodGlobal:
		return "Global"
	default:
		return "Unknown"
	}
}

const (
	MethodLocally = 1
	MethodDelete  = 2
	MethodGlobal  = 3
)
