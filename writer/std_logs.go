package writer

import (
	"os"
	"syscall"
)

type LogsStandardOutputWriter struct {
	level Level
	out *os.File
	logPrefix string
	logFlags int
	transformFunc TransformFunc
}


func NewLogsStandardOutputWriter(level Level, prefix string, flags int) *LogsStandardOutputWriter {
	w := new(LogsStandardOutputWriter)
	w.out = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")
	w.logPrefix = prefix
	w.logFlags = flags
	w.level = level
	w.transformFunc = LogsStandardTextLineTransform
	return w
}

func NewLogsStandardOutputWriterWithTransform(level Level, prefix string, flags int, transformFunc TransformFunc) *LogsStandardOutputWriter {
	w := new(LogsStandardOutputWriter)
	w.out = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")
	w.logPrefix = prefix
	w.logFlags = flags
	w.transformFunc = transformFunc
	w.level = level
	return w
}

func (w *LogsStandardOutputWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	msg, decodeErr := LogMsgDecode(p, w.logPrefix, w.logFlags)
	if decodeErr != nil {
		panic(decodeErr)
		return 0, decodeErr
	}
	levelParseErr := parseLevelFromMsg(&msg)
	if levelParseErr != nil {
		msg.LevelNo = Level(10)
		msg.LevelName = "???"
	}
	if levelGate(msg.LevelNo, w.level) {
		b, transformErr := w.transformFunc(msg)
		if transformErr != nil {
			return 0, transformErr
		}
		if msg.NewLineFlag {
			b = append(b, '\n')
		}
		return w.out.Write(b)
	}
	return 0, nil
}

