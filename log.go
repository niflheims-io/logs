package logs

import (
	"sync"
	"github.com/niflheims-io/logs/writer"
)

type Logs struct  {
	outputs []output
	lk *sync.Mutex
}

func New() *Logs {
	bitLog := &Logs{outputs:make([]output,0,1), lk:new(sync.Mutex)}
	return bitLog
}

func (self *Logs) AddOutput(level int, f Formatter, w writer.Writer)  {
	self.lk.Lock()
	defer self.lk.Unlock()
	p := output{level:level, formatter:f, writer:w, logger: self}
	self.outputs = append(self.outputs, p)
}

func (self *Logs) Print() *printer {
	p := newPrinter(self)
	return p
}

func (self *Logs) doLog(dp DataPkg) {
	outputLen := len(self.outputs)
	if outputLen == 0 {
		w := writer.NewStdWriter(false)
		f := NewTextLineFormatter();
		self.AddOutput(LEVEL_DEBUG, f, w)
	}
	for i := 0 ; i < outputLen ; i ++ {
		op := self.outputs[i]
		op.on(dp)
	}
}









