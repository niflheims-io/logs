package logs

import "github.com/niflheims-io/logs/writer"

const (
	LEVEL_FATAL		= 1
	LEVEL_PANIC		= 2
	LEVEL_ERROR 		= 3
	LEVEL_WARN 		= 4
	LEVEL_INFO 		= 5
	LEVEL_DEBUG 		= 6
	level_name_fatal 	= "FATAL"
	level_name_panic 	= "PANIC"
	level_name_error 	= "ERROR"
	level_name_warn 	= "WARN"
	level_name_info 	= "INFO"
	level_name_debug 	= "DEBUG"
)

type output struct  {
	level 		int
	writer 		writer.Writer
	formatter 	Formatter
	logger 		*Logs
}

func (self *output) on(dp DataPkg) {
	if self.level == 0 || self.level < dp.Level {
		return
	}
	outputData, fmtErr := self.formatter.Format(&dp)
	if fmtErr != nil {
		self.logger.Print().Panic(fmtErr)
		return
	}
	_, outputErr := self.writer.Write(outputData)
	if outputErr != nil {
		self.logger.Print().Panic(outputErr)
		return
	}
	if self.level == LEVEL_PANIC {
		panic(dp)
	}
}
