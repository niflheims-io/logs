package writer

const (
	LEVEL_ERROR = Level(1)
	LEVEL_WARN = Level(2)
	LEVEL_INFO = Level(3)
	LEVEL_DEBUG = Level(4)
	LEVEL_TRACE = Level(5)

	LEVEL_ERROR_NAME = "ERROR"
	LEVEL_WARN_NAME = "WARN"
	LEVEL_INFO_NAME = "INFO"
	LEVEL_DEBUG_NAME = "DEBUG"
	LEVEL_TRACE_NAME = "TRACE"

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
		return LEVEL_ERROR_NAME
	} else if level == LEVEL_WARN {
		return LEVEL_WARN_NAME
	} else if level == LEVEL_INFO {
		return LEVEL_INFO_NAME
	} else if level == LEVEL_DEBUG {
		return LEVEL_DEBUG_NAME
	} else if level == LEVEL_TRACE {
		return LEVEL_TRACE_NAME
	} else {
		return "???"
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
	if levelName == LEVEL_ERROR_NAME {
		return LEVEL_ERROR
	} else if levelName == LEVEL_WARN_NAME {
		return LEVEL_WARN
	} else if levelName == LEVEL_INFO_NAME {
		return LEVEL_INFO
	} else if levelName == LEVEL_DEBUG_NAME {
		return LEVEL_DEBUG
	} else if levelName == LEVEL_TRACE_NAME {
		return LEVEL_TRACE
	} else {
		return Level(-1)
	}
}

func levelGate(levelLog Level, levelWriter Level) bool {
	return int(levelLog) <= int(levelWriter)
}

