package syncer

type State uint8

func (s State) String() string {
	switch s {
	case StateNoChange:
		return "NoChange"
	case StateInsert:
		return "Insert"
	case StateUpdate:
		return "Update"
	case StateDelete:
		return "Delete"
	default:
		return "Unknown"
	}
}

const (
	StateNoChange State = 0
	StateInsert   State = 1
	StateUpdate   State = 2
	StateDelete   State = 3
)
