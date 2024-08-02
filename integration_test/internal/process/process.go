package process

import (
	"context"
	"github.com/openimsdk/tools/errs"
)

type Process struct {
	ctx   context.Context
	Tasks []*Task

	// Only all RunCondition is true, process will exec tasks.
	RunConditions []bool

	nowTaskNum int // now un exec task number
}

func NewProcess() *Process {
	return &Process{}
}

func (p *Process) GetTaskNum() int {
	return len(p.Tasks)
}

func (p *Process) AddTasks(task ...*Task) {
	if len(task) != 0 {
		p.Tasks = append(p.Tasks, task...)
	}
}

func (p *Process) AddConditions(condition ...bool) {
	if len(condition) != 0 {
		p.RunConditions = append(p.RunConditions, condition...)
	}
}

func (p *Process) ResetConditions(condition ...bool) {
	p.RunConditions = condition
}

func (p *Process) Exec() error {
	return p.ExecOffset(0)
}

func (p *Process) ContinueExec() error {
	return p.ExecOffset(p.nowTaskNum)
}

func (p *Process) ExecOffset(offset int) error {
	if offset < 0 || offset > len(p.Tasks) {
		return errs.New("err input offset is process exec").Wrap()
	}
	var (
		interrupt = -1
	)

	return p.execTasks(offset, interrupt)
}

func (p *Process) execTasks(offset, interrupt int) error {
	p.nowTaskNum = offset
	for _, task := range p.Tasks[offset:] {
		if p.nowTaskNum == interrupt {
			return nil
		}
		if p.shouldRun() && task.ShouldRun {
			for _, f := range task.Funcs {
				if err := f(p.ctx); err != nil {
					return err
				}
			}
		} else {
			for _, f := range task.NegativeFuncs {
				if err := f(p.ctx); err != nil {
					return err
				}
			}
		}
		p.nowTaskNum++
	}
	return nil
}

func (p *Process) SetContext(ctx context.Context) {
	p.ctx = ctx
}

func (p *Process) shouldRun() bool {
	if len(p.RunConditions) == 0 {
		return true
	}
	for _, cond := range p.RunConditions {
		if !cond {
			return false
		}
	}
	return true
}

func (p *Process) Clear() {
	p.Tasks = nil
	p.RunConditions = nil
	p.nowTaskNum = 0
}

type Task struct {
	ShouldRun     bool                              // determine if task will run
	Funcs         []func(ctx context.Context) error // run funcs
	NegativeFuncs []func(ctx context.Context) error // if !ShouldRun, will run this
}

func NewTask(shouldRun bool, funcs ...func(ctx context.Context) error) *Task {
	return &Task{
		ShouldRun: shouldRun,
		Funcs:     funcs,
	}
}

func (t *Task) AddNegativeFuncs(funcs ...func(ctx context.Context) error) *Task {
	t.NegativeFuncs = append(t.NegativeFuncs, funcs...)
	return t
}

// WrapFunc wrap common func
func WrapFunc(f func()) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		f()
		return nil
	}
}
