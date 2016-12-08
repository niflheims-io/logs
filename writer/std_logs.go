package writer

import (
	"os"
	"syscall"
	"errors"
	"bytes"
	"fmt"
	"strings"
)

type LogsStandardOutputWriter struct {
	level Level
	out *os.File
	prefix string
	flags int
	transformFunc TransformFunc
}

func LogsStandardTextLineTransform(msg LogMsg) ([]byte, error) {
	line := ""
	if msg.Prefix != "" {
		line = line + msg.Prefix + " "
	}
	if msg.LevelNo != Level(-1) {
		line = line + getLevelColor(msg.LevelNo) + "[" + msg.LevelName + "]" + color_default + " "
	}
	if !msg.DateTime.IsZero() {
		line = line + color_blue + msg.DateTime.String() + color_default + " "
	}
	if msg.File != "" {
		line = line + msg.File + " "
	}
	if msg.Msg != "" {
		line = line + msg.Msg
	}
	if len(msg.ExtendFields) > 0 {
		line = line + " "
		for key, value := range msg.ExtendFields {
			line = line + "[" + color_attr + key + color_default + ":" + value + "]"
		}
	}
	return []byte(line), nil
}

func NewLogsStandardOutputWriter(prefix string, flags int, transformFunc TransformFunc) *LogsStandardOutputWriter {
	w := new(LogsStandardOutputWriter)
	w.out = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")
	w.prefix = prefix
	w.flags = flags
	w.transformFunc = transformFunc
	return w
}

func NewLogsStandardOutputWriterWithLevel(prefix string, flags int, level Level, transformFunc TransformFunc) *LogsStandardOutputWriter {
	w := new(LogsStandardOutputWriter)
	w.out = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")
	w.prefix = prefix
	w.flags = flags
	w.transformFunc = transformFunc
	w.level = level
	return w
}

func (w *LogsStandardOutputWriter) Write(p []byte) (int, error) {
	msg, decodeErr := LogMsgDecode(p, w.prefix, w.flags)
	if decodeErr != nil {
		return 0, decodeErr
	}
	levelParseErr := parseLevelFromMsg(&msg)
	if levelParseErr != nil {
		msg.LevelNo = Level(-1)
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


func parseLevelFromMsg(msg *LogMsg) error {
	if msg == nil {
		return errors.New("No message to log.")
	}
	bodyBytesBuf := bytes.NewBufferString(msg.Msg)
	level, levelReadErr := bodyBytesBuf.ReadString(' ')
	if levelReadErr != nil {
		return  errors.New("Can not get message level. " + msg.Msg + " . error:" + fmt.Sprint(levelReadErr))
	}
	level = strings.TrimSpace(level)
	isLevel := level == LEVEL_ERROR_NAME || level == LEVEL_WARN_NAME || level == LEVEL_INFO_NAME || level == LEVEL_DEBUG_NAME || level == LEVEL_TRACE_NAME
	if !isLevel {
		return  errors.New("Can not get message level. " + msg.Msg)
	}
	msg.LevelName = level
	msg.LevelNo = transferLevelNameToLevelNo(level)
	msg.Msg = bodyBytesBuf.String()
	if level == LEVEL_ERROR_NAME {
		msg.Msg = color_red + msg.Msg + color_default
	} else if level == LEVEL_WARN_NAME {
		msg.Msg = color_yellow + msg.Msg + color_default
	} else if level == LEVEL_INFO_NAME {
		msg.Msg = color_default + msg.Msg + color_default
	} else if level == LEVEL_DEBUG_NAME {
		msg.Msg = color_green + msg.Msg + color_default
	} else if level == LEVEL_TRACE_NAME {
		msg.Msg = color_trace + msg.Msg + color_default
	}
	return nil
}

