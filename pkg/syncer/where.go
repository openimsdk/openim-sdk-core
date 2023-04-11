package syncer

type Column struct {
	Name  string
	Value any
}

type Where struct {
	Columns []*Column
}
