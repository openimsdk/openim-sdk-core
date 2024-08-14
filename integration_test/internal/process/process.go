package process

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"reflect"
	"runtime"
	"runtime/debug"
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
			if task.ShouldRun {
				if err := p.call(task.Func, task.Args...); err != nil {
					return err
				}
			} else {
				if err := p.call(task.NegativeFunc, task.NegativeArgs...); err != nil {
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

func (p *Process) call(fn any, args ...any) (err error) {
	funcPtr := reflect.ValueOf(fn).Pointer()
	funcName := runtime.FuncForPC(funcPtr).Name()
	defer func() {
		ctx := context.TODO()
		if r := recover(); r != nil {
			fmt.Printf("panic: %+v\n%s", r, debug.Stack())
			err = fmt.Errorf("call panic: %+v", r)
		} else {
			if err == nil {
				log.ZInfo(ctx, "fn call success", "function name", funcName)
			} else {
				log.ZError(ctx, "fn call error", err, "function name", funcName)
			}
		}
	}()

	fnv := reflect.ValueOf(fn)
	if fnv.Kind() != reflect.Func {
		return errs.New("call input type is not func").Wrap()
	}
	if fnv.IsNil() {
		return nil
	}

	fnt := fnv.Type()
	nin := fnt.NumIn()
	var (
		ins         = make([]reflect.Value, 0, nin)
		fieldOffset = 0 // if have ctx, parse field offset = 1
		isVariadic  = fnt.IsVariadic()
	)

	if nin != 0 {
		argsLen := len(args)
		// If there are parameters, the first parameter must be ctx
		if fnt.In(0).Implements(reflect.ValueOf(new(context.Context)).Elem().Type()) {
			fieldOffset = 1
			ins = append(ins, reflect.ValueOf(p.ctx))
			argsLen++
		}
		if isVariadic {
			nin--
		}
		if (isVariadic && argsLen < nin) || (!isVariadic && argsLen != nin) {
			return errs.New(fmt.Sprintf("call input args num not correct. nin: %d, argsLen: %d", nin, args)).Wrap()
		}
	}

	for i := 0; i < len(args); i++ {
		inFnField := fnt.In(i + fieldOffset)

		arg := reflect.TypeOf(args[i])

		var expectedType reflect.Type
		if isVariadic && i+fieldOffset >= nin {
			expectedType = fnt.In(nin).Elem()
		} else {
			expectedType = fnt.In(i + fieldOffset)
		}

		if arg.AssignableTo(expectedType) {
			ins = append(ins, reflect.ValueOf(args[i]))
			continue
		}
		if arg.Kind() == reflect.String { // json
			var ptr int
			for inFnField.Kind() == reflect.Ptr {
				inFnField = inFnField.Elem()
				ptr++
			}
			switch inFnField.Kind() {
			case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map:
				v := reflect.New(inFnField)
				if err := json.Unmarshal([]byte(args[i].(string)), v.Interface()); err != nil {
					return errs.New(fmt.Sprintf("go call json.Unmarshal error: %s", err)).Wrap()
				}
				if ptr == 0 {
					v = v.Elem()
				} else if ptr != 1 {
					for i := ptr - 1; i > 0; i-- {
						temp := reflect.New(v.Type())
						temp.Elem().Set(v)
						v = temp
					}
				}
				ins = append(ins, v)
				continue
			}
		}
		return errs.New(fmt.Sprintf("go code error: fn in args type is not match. index:%d, field type:%s, arg type:%s",
			i, inFnField.String(), arg.String())).Wrap()
	}
	outs := fnv.Call(ins)
	if len(outs) == 0 {
		return nil
	}
	if fnt.Out(len(outs) - 1).Implements(reflect.ValueOf(new(error)).Elem().Type()) {
		if errValueOf := outs[len(outs)-1]; !errValueOf.IsNil() {
			return errValueOf.Interface().(error)
		}
	}

	return nil
}

func (p *Process) Clear() {
	p.Tasks = nil
	p.RunConditions = nil
	p.nowTaskNum = 0
}
