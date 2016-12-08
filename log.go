package logs

import (
	"log"
	"io"
	"github.com/niflheims-io/logs/writer"
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
	l.Output(2, fmt.Sprint(writer.LEVEL_ERROR_NAME) + " " + fmt.Sprint(v...))
}

func (l *Logs) ErrorF(format string, v ...interface{}) {
	l.Output(2, fmt.Sprint(writer.LEVEL_ERROR_NAME) + " " + fmt.Sprintf(format, v...))
}

func (l *Logs) Warn(v ...interface{}) {
	l.Output(2, fmt.Sprint(writer.LEVEL_WARN_NAME) + " " + fmt.Sprint(v...))
}

func (l *Logs) WarnF(format string, v ...interface{}) {
	l.Output(2, fmt.Sprint(writer.LEVEL_WARN_NAME) + " " + fmt.Sprintf(format, v...))
}


func (l *Logs) Info(v ...interface{}) {
	l.Output(2, fmt.Sprint(writer.LEVEL_INFO_NAME) + " " + fmt.Sprint(v...))
}

func (l *Logs) InfoF(format string, v ...interface{}) {
	l.Output(2, fmt.Sprint(writer.LEVEL_INFO_NAME) + " " + fmt.Sprintf(format, v...))
}

func (l *Logs) Debug(v ...interface{}) {
	l.Output(2, fmt.Sprint(writer.LEVEL_DEBUG_NAME) + " " + fmt.Sprint(v...))
}

func (l *Logs) DebugF(format string, v ...interface{}) {
	l.Output(2, fmt.Sprint(writer.LEVEL_DEBUG_NAME) + " " + fmt.Sprintf(format, v...))
}

func (l *Logs) Trace(v ...interface{}) {
	l.Output(2, fmt.Sprint(writer.LEVEL_TRACE_NAME) + " " + fmt.Sprint(v...))
}

func (l *Logs) TraceF(format string, v ...interface{}) {
	l.Output(2, fmt.Sprint(writer.LEVEL_TRACE_NAME) + " " + fmt.Sprintf(format, v...))
}
