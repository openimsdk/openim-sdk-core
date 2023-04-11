package syncer

type State uint8

func (s State) String() string {
	switch s {
	case StateUnchanged:
		return "Unchanged"
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
	StateUnchanged State = 0
	StateInsert    State = 1
	StateUpdate    State = 2
	StateDelete    State = 3
)
