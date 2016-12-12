package logs

import (
	"log"
	"io"
	"fmt"
	"time"
	"runtime"
	"sync"
	"os"
)

type Logs struct  {
	*log.Logger
	out io.Writer
	async bool
	msgChan chan []byte
	wg *sync.WaitGroup
	closeFlg bool
}

func New(out io.Writer, prefix string, flag int) *Logs {
	logger := Logs{Logger:log.New(out, prefix, flag)}
	logger.out = out
	logger.async = false
	return &logger
}

func NewAsyncLogs(out io.Writer, prefix string, flag int) *Logs {
	logger := Logs{Logger:log.New(out, prefix, flag)}
	logger.out = out
	logger.async = true
	logger.msgChan = make(chan []byte, 1024)
	logger.wg = new(sync.WaitGroup)
	logger.closeFlg = false
	logger.receive()
	return &logger
}

func (l *Logs) Error(v ...interface{}) {
	if l.async {
		l.OutputAsync(2, "ERROR" + " " + fmt.Sprint(v...))
	} else {
		l.Output(2, "ERROR" + " " + fmt.Sprint(v...))
	}
}

func (l *Logs) ErrorF(format string, v ...interface{}) {
	if l.async {
		l.OutputAsync(2, "ERROR" + " " + fmt.Sprintf(format, v...))
	} else {
		l.Output(2, "ERROR" + " " + fmt.Sprintf(format, v...))
	}
}

func (l *Logs) Warn(v ...interface{}) {
	if l.async {
		l.OutputAsync(2, "WARN" + " " + fmt.Sprint(v...))
	} else {
		l.Output(2, "WARN" + " " + fmt.Sprint(v...))
	}
}

func (l *Logs) WarnF(format string, v ...interface{}) {
	if l.async {
		l.OutputAsync(2, "WARN" + " " + fmt.Sprintf(format, v...))
	} else {
		l.Output(2, "WARN" + " " + fmt.Sprintf(format, v...))
	}
}


func (l *Logs) Info(v ...interface{}) {
	if l.async {
		l.OutputAsync(2, "INFO" + " " + fmt.Sprint(v...))
	} else {
		l.Output(2, "INFO" + " " + fmt.Sprint(v...))
	}
}

func (l *Logs) InfoF(format string, v ...interface{}) {
	if l.async {
		l.OutputAsync(2, "INFO" + " " + fmt.Sprintf(format, v...))
	} else {
		l.Output(2, "INFO" + " " + fmt.Sprintf(format, v...))
	}
}

func (l *Logs) Debug(v ...interface{}) {
	if l.async {
		l.OutputAsync(2, "DEBUG" + " " + fmt.Sprint(v...))
	} else {
		l.Output(2, "DEBUG" + " " + fmt.Sprint(v...))
	}
}

func (l *Logs) DebugF(format string, v ...interface{}) {
	if l.async {
		l.OutputAsync(2, "DEBUG" + " " + fmt.Sprintf(format, v...))
	} else {
		l.Output(2, "DEBUG" + " " + fmt.Sprintf(format, v...))
	}
}

func (l *Logs) Trace(v ...interface{}) {
	if l.async {
		l.OutputAsync(2, "TRACE" + " " + fmt.Sprint(v...))
	} else {
		l.Output(2, "TRACE" + " " + fmt.Sprint(v...))
	}
}

func (l *Logs) TraceF(format string, v ...interface{}) {
	if l.async {
		l.OutputAsync(2, "TRACE" + " " + fmt.Sprintf(format, v...))
	} else {
		l.Output(2, "TRACE" + " " + fmt.Sprintf(format, v...))
	}
}

func (l *Logs) Printf(format string, v ...interface{}) {
	if l.async {
		l.OutputAsync(2, fmt.Sprintf(format, v...))
	} else {
		l.Output(2, fmt.Sprintf(format, v...))
	}
}

func (l *Logs) Println(v ...interface{})  {
	if l.async {
		l.OutputAsync(2, fmt.Sprintln(v...))
	} else {
		l.Output(2, fmt.Sprintln(v...))
	}
}

func (l *Logs) Print(v ...interface{})  {
	if l.async {
		l.OutputAsync(2, fmt.Sprint(v...))
	} else {
		l.Output(2, fmt.Sprint(v...))
	}
}

// fatal
func (l *Logs) Fatalf(format string, v ...interface{}) {
	if l.async {
		l.OutputAsync(2, fmt.Sprintf(format, v...))
		l.Close()
	} else {
		l.Output(2, fmt.Sprintf(format, v...))
	}
	os.Exit(1)
}

func (l *Logs) Fatalln(v ...interface{})  {
	if l.async {
		l.OutputAsync(2, fmt.Sprintln(v...))
		l.Close()
	} else {
		l.Output(2, fmt.Sprintln(v...))
	}
	os.Exit(1)
}

func (l *Logs) Fatal(v ...interface{})  {
	if l.async {
		l.OutputAsync(2, fmt.Sprint(v...))
		l.Close()
	} else {
		l.Output(2, fmt.Sprint(v...))
	}
	os.Exit(1)
}

// panic
func (l *Logs) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	if l.async {
		l.OutputAsync(2, s)
	} else {
		l.Output(2, s)
	}
	panic(s)
}

func (l *Logs) Panicln(v ...interface{})  {
	s := fmt.Sprintln(v...)
	if l.async {
		l.OutputAsync(2, s)
	} else {
		l.Output(2, s)
	}
	panic(s)
}

func (l *Logs) Panic(v ...interface{})  {
	s := fmt.Sprint(v...)
	if l.async {
		l.OutputAsync(2, s)
	} else {
		l.Output(2, s)
	}
	panic(s)
}

func (l *Logs) OutputAsync(callDepth int, s string) error {
	if l.closeFlg {
		return nil
	}
	now := time.Now() // get this early.
	var file string
	var line int
	if l.Flags()&(log.Lshortfile|log.Llongfile) != 0 {
		var ok bool
		_, file, line, ok = runtime.Caller(callDepth)
		if !ok {
			file = "???"
			line = 0
		}
	}
	msgBytes := make([]byte, 0, len(s))
	l.formatHeader(&msgBytes, now, file, line)
	msgBytes = append(msgBytes, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		msgBytes = append(msgBytes, '\n')
	}
	l.msgChan <- msgBytes
	return nil
}

func (l *Logs) formatHeader(buf *[]byte, t time.Time, file string, line int) {
	*buf = append(*buf, l.Prefix()...)
	if l.Flags()&log.LUTC != 0 {
		t = t.UTC()
	}
	if l.Flags()&(log.Ldate|log.Ltime|log.Lmicroseconds) != 0 {
		if l.Flags()&log.Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.Flags()&(log.Ltime|log.Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.Flags()&log.Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.Flags()&(log.Lshortfile|log.Llongfile) != 0 {
		if l.Flags()&log.Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}

func itoa(buf *[]byte, i int, wid int) {
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (l *Logs) receive() {
	go func(l *Logs) {
		for {
			l.wg.Add(1)
			msgBytes, ok := <- l.msgChan
			if !ok {
				l.wg.Done()
				return
			}
			n, err := l.out.Write(msgBytes)
			l.wg.Done()
			if err != nil {
				panic(fmt.Sprintf("write msg fail, writed size : %d, error : %v.", n, err))
			}

		}
	}(l)
}

func (l *Logs) Close()  {
	if l.async {
		l.closeFlg = true
		close(l.msgChan)
		l.wg.Wait()
	}

}