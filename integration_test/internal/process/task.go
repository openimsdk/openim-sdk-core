package process

type Task struct {
	ShouldRun    bool  // determine if task will run
	Func         any   // must be func. run funcs
	Args         []any // func args
	NegativeFunc any   // must be func. if !ShouldRun, will run this
	NegativeArgs []any // negative args
}

func NewTask(shouldRun bool, f any, args ...any) *Task {
	return &Task{
		ShouldRun: shouldRun,
		Func:      f,
		Args:      args,
	}
}

func (t *Task) AddNegativeFunc(f any, args ...any) *Task {
	t.NegativeFunc = f
	t.NegativeArgs = args
	return t
}
