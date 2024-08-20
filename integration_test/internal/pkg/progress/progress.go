package progress

import (
	"bytes"
	"fmt"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/formatutil"
	"io"
	"os"
	"sync"
)

type proFlag uint8

const (
	AutoClose      proFlag = 1 << iota // if set, progress will close when all bars done
	ForbiddenWrite                     // if set, progress will if forbidden other goroutine print to Stdout through os.Stdout
)

type signalType uint8

const (
	start signalType = iota // do
	update
	stop
)

const (
	maxBars = 10 // default max print bars
)

func NewProgress(mode proFlag, m int) *Progress {
	if m == 0 {
		m = maxBars
	}
	return &Progress{
		signal: make(chan signalType, 1),
		stdout: os.Stdout,
		buf:    bytes.Buffer{},
		done:   make(chan struct{}),

		MaxPrintBar: m,
		mode:        mode,
	}
}

type Progress struct {
	stdout     *os.File // os.Stdout, read only!
	pipeWriter *os.File // Acts as a temporary writer for os.Stdout during the write prohibition period
	// buf used to store the data that is printed when input is disabled, to be output after the prohibition is lifted
	buf  bytes.Buffer
	done chan struct{} // Used to record whether copying to the buffer is complete

	Bars        []*Bar
	MaxPrintBar int // The maximum number of bars to print; any excess will be recorded but not displayed.
	printLine   int

	mode   proFlag
	signal chan signalType
	lock   sync.Mutex
}

func (p *Progress) forbiddenPrint() {
	r, w, _ := os.Pipe()
	p.pipeWriter = w
	os.Stdout = p.pipeWriter
	p.done = make(chan struct{})
	p.buf = bytes.Buffer{}
	go func() {
		_, _ = io.Copy(&p.buf, r)
		close(p.done)
	}()
}

func (p *Progress) allowPrint() {
	_ = p.pipeWriter.Close()
	os.Stdout = p.stdout

	<-p.done
	// print buf
	fmt.Print(p.buf.String())
}

func (p *Progress) AddBar(bs ...*Bar) {
	if len(bs) == 0 {
		return
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.Bars = append(p.Bars, bs...)
	p.notifyUpdate()
}

func (p *Progress) notifyUpdate() {
	select {
	case p.signal <- update:
	default:
	}
}

func (p *Progress) run() {
	for {
		signal := <-p.signal
		switch signal {
		case start:
			p.start()
			fallthrough
		case update:
			p.render()
		case stop:
			p.stop()
			return
		default:
			return
		}
	}
}

func (p *Progress) render() {
	printStr := ""

	printStr += saveCursor()
	printStr += cursorUpHead(p.printLine)
	p.printLine = 0
	for i, bar := range p.Bars {
		if bar.shouldRemove() && len(p.Bars) > p.MaxPrintBar {
			p.lock.Lock()
			datautil.DeleteAt(&p.Bars, i)
			p.lock.Unlock()
			continue
		}
		printStr += clearLine()
		printStr += formatutil.ProgressBar(bar.name, bar.now, bar.total)
		printStr += nextLine()
		p.printLine++
		if p.printLine >= p.MaxPrintBar {
			break
		}
	}
	printStr += loadCursor()
	p.print(printStr)

	// auto close
	if p.printLine == 0 && p.mode&AutoClose == AutoClose {
		select {
		case p.signal <- stop:
		default:
			// still have signal not been completed
		}
	}
}

func (p *Progress) start() {
	if p.mode&ForbiddenWrite == ForbiddenWrite {
		p.forbiddenPrint()
	}
}

func (p *Progress) stop() {
	if p.mode&ForbiddenWrite == ForbiddenWrite {
		p.allowPrint()
	}
	//close(p.signal)
}

func (p *Progress) print(s string) {
	_, _ = fmt.Fprint(p.stdout, s)
}

func nextLine() string {
	return "\n"
}

func (p *Progress) Start() {
	go p.run()
	p.signal <- start
}

func (p *Progress) Stop() {
	if p.IsStopped() {
		return
	}
	p.signal <- stop
	<-p.done
}

func (p *Progress) IsStopped() bool {
	select {
	case _, _ = <-p.done:
		// already done
		return true
	default:
		return false
	}
}

func (p *Progress) IncBar(b *Bar) {
	p.lock.Lock()
	defer p.lock.Unlock()
	b.now++
	p.notifyUpdate()
}

func (p *Progress) SetMaxPrintBar(n int) {
	p.MaxPrintBar = n
}
