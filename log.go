package logs

import (
	"log"
	"io"
	"fmt"
)

type Logs struct  {
	*log.Logger
}

func New(out io.Writer, prefix string, flag int) *Logs {
	logger := Logs{log.New(out, prefix, flag)}
	return &logger
}

func (l *Logs) Error(v ...interface{}) {
	l.Output(2, "ERROR" + " " + fmt.Sprint(v...))
}

func (l *Logs) ErrorF(format string, v ...interface{}) {
	l.Output(2, "ERROR" + " " + fmt.Sprintf(format, v...))
}

func (l *Logs) Warn(v ...interface{}) {
	l.Output(2, "WARN" + " " + fmt.Sprint(v...))
}

func (l *Logs) WarnF(format string, v ...interface{}) {
	l.Output(2, "WARN" + " " + fmt.Sprintf(format, v...))
}


func (l *Logs) Info(v ...interface{}) {
	l.Output(2, "INFO" + " " + fmt.Sprint(v...))
}

func (l *Logs) InfoF(format string, v ...interface{}) {
	l.Output(2, "INFO" + " " + fmt.Sprintf(format, v...))
}

func (l *Logs) Debug(v ...interface{}) {
	l.Output(2, "DEBUG" + " " + fmt.Sprint(v...))
}

func (l *Logs) DebugF(format string, v ...interface{}) {
	l.Output(2, "DEBUG" + " " + fmt.Sprintf(format, v...))
}

func (l *Logs) Trace(v ...interface{}) {
	l.Output(2, "TRACE" + " " + fmt.Sprint(v...))
}

func (l *Logs) TraceF(format string, v ...interface{}) {
	l.Output(2, "TRACE" + " " + fmt.Sprintf(format, v...))
}
