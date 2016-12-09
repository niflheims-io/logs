package writer

import (
	"bytes"
	"fmt"
	"strings"
	"errors"
)

const (
	LEVEL_ERROR = Level(1)
	LEVEL_WARN = Level(2)
	LEVEL_INFO = Level(3)
	LEVEL_DEBUG = Level(4)
	LEVEL_TRACE = Level(5)

	level_error_name = "ERROR"
	level_warn_name = "WARN"
	level_info_name = "INFO"
	level_debug_name = "DEBUG"
	level_trace_name = "TRACE"

	color_default	= "\x1b[0m"
	color_red 	= "\x1b[31m"
	color_green 	= "\x1b[32m"
	color_yellow	= "\x1b[33m"
	color_blue 	= "\x1b[34m"
	color_attr	= "\x1b[35m"
	color_trace	= "\x1b[37m"
)

type Level int

func transferLevelNoToLevelName(level Level) string {
	if level == LEVEL_ERROR {
		return level_error_name
	} else if level == LEVEL_WARN {
		return level_warn_name
	} else if level == LEVEL_INFO {
		return level_info_name
	} else if level == LEVEL_DEBUG {
		return level_debug_name
	} else if level == LEVEL_TRACE {
		return level_trace_name
	} else {
		return "*"
	}
}

func getLevelColor(level Level) string {
	if level == LEVEL_ERROR {
		return color_red
	} else if level == LEVEL_WARN {
		return color_yellow
	} else if level == LEVEL_INFO {
		return color_default
	} else if level == LEVEL_DEBUG {
		return color_green
	} else if level == LEVEL_TRACE {
		return color_trace
	} else {
		return color_default
	}
}

func transferLevelNameToLevelNo(levelName string) Level {
	if levelName == level_error_name {
		return LEVEL_ERROR
	} else if levelName == level_warn_name {
		return LEVEL_WARN
	} else if levelName == level_info_name {
		return LEVEL_INFO
	} else if levelName == level_debug_name {
		return LEVEL_DEBUG
	} else if levelName == level_trace_name {
		return LEVEL_TRACE
	} else {
		return LEVEL_INFO
	}
}

func levelGate(levelLog Level, levelWriter Level) bool {
	return int(levelLog) <= int(levelWriter)
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
	isLevel := level == level_error_name || level == level_warn_name || level == level_info_name || level == level_debug_name || level == level_trace_name
	if !isLevel {
		level = level_info_name
	}
	msg.LevelName = level
	msg.LevelNo = transferLevelNameToLevelNo(level)
	msg.Msg = bodyBytesBuf.String()
	if level == level_error_name {
		msg.Msg = color_red + msg.Msg + color_default
	} else if level == level_warn_name {
		msg.Msg = color_yellow + msg.Msg + color_default
	} else if level == level_info_name {
		msg.Msg = color_default + msg.Msg + color_default
	} else if level == level_debug_name {
		msg.Msg = color_green + msg.Msg + color_default
	} else if level == level_trace_name {
		msg.Msg = color_trace + msg.Msg + color_default
	}
	return nil
}