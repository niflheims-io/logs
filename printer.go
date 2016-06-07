package logs

import (
	"fmt"
	"runtime"
	"os"
	"strings"
	"time"
)

type printer struct  {
	logger *Logs
	dp DataPkg
}

func newPrinter(logger *Logs) *printer {
	return &printer{logger:logger, dp:DataPkg{Time:time.Now()}}
}

func (self *printer) Attr(attr AttrMap) *printer {
	self.dp.Attr = attr
	return self
}

func (self *printer) LineNo() *printer {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	self.dp.Position = file + ":" + fmt.Sprint(line)
	return self
}

func (self *printer) Debug(v ...interface{})  {
	msg := ""
	for _, i := range v {
		msg = msg + fmt.Sprint(i) + " "
	}
	msg = strings.TrimSpace(msg)
	self.dp.Message = msg
	self.dp.Level = LEVEL_DEBUG
	self.send()
}

func (self *printer) DebugF(format string, v ...interface{})  {
	msg := ""
	for _, i := range v {
		msg = msg + fmt.Sprintf(format, i) + " "
	}
	self.dp.Message = msg
	self.dp.Level = LEVEL_DEBUG
	self.send()
}

func (self *printer) Info(v ...interface{})  {
	msg := ""
	for _, i := range v {
		msg = msg + fmt.Sprint(i) + " "
	}
	msg = strings.TrimSpace(msg)
	self.dp.Message = msg
	self.dp.Level = LEVEL_INFO
	self.send()
}

func (self *printer) InfoF(format string, v ...interface{})  {
	msg := ""
	for _, i := range v {
		msg = msg + fmt.Sprintf(format, i) + " "
	}
	self.dp.Message = msg
	self.dp.Level = LEVEL_INFO
	self.send()
}

func (self *printer) Warn(v ...interface{})  {
	msg := ""
	for _, i := range v {
		msg = msg + fmt.Sprint(i) + " "
	}
	msg = strings.TrimSpace(msg)
	self.dp.Message = msg
	self.dp.Level = LEVEL_WARN
	self.send()
}

func (self *printer) WarnF(format string, v ...interface{})  {
	msg := ""
	for _, i := range v {
		msg = msg + fmt.Sprintf(format, i) + " "
	}
	self.dp.Message = msg
	self.dp.Level = LEVEL_WARN
	self.send()
}

func (self *printer) Error(v ...interface{})  {
	msg := ""
	for _, i := range v {
		msg = msg + fmt.Sprint(i) + " "
	}
	msg = strings.TrimSpace(msg)
	self.dp.Message = msg
	self.dp.Level = LEVEL_ERROR
	self.send()
}

func (self *printer) ErrorF(format string, v ...interface{})  {
	msg := ""
	for _, i := range v {
		msg = msg + fmt.Sprintf(format, i) + " "
	}
	self.dp.Message = msg
	self.dp.Level = LEVEL_ERROR
	self.send()
}

func (self *printer) Panic(v ...interface{})  {
	msg := ""
	for _, i := range v {
		msg = msg + fmt.Sprint(i) + " "
	}
	msg = strings.TrimSpace(msg)
	self.dp.Message = msg
	self.dp.Level = LEVEL_PANIC
	self.send()
}

func (self *printer) PanicF(format string, v ...interface{})  {
	msg := ""
	for _, i := range v {
		msg = msg + fmt.Sprintf(format, i) + " "
	}
	self.dp.Message = msg
	self.dp.Level = LEVEL_PANIC
	self.send()
}

func (self *printer) Fatal(v ...interface{})  {
	msg := ""
	for _, i := range v {
		msg = msg + fmt.Sprint(i) + " "
	}
	msg = strings.TrimSpace(msg)
	self.dp.Message = msg
	self.dp.Level = LEVEL_FATAL
	self.send()
	os.Exit(1)
}

func (self *printer) FatalF(format string, v ...interface{})  {
	msg := ""
	for _, i := range v {
		msg = msg + fmt.Sprintf(format, i) + " "
	}
	msg = strings.TrimSpace(msg)
	self.dp.Message = msg
	self.dp.Level = LEVEL_FATAL
	self.send()
	os.Exit(1)
}

func (self *printer) send()  {
	self.logger.doLog(self.dp)
}