package logs

import (
	"time"
)

type DataPkg struct {
	Time 		time.Time
	Position 	string
	Level 		int
	Attr		AttrMap
	Message		string
}

func (self *DataPkg) levelName() string {
	if self.Level == LEVEL_DEBUG {
		return level_name_debug
	} else if self.Level == LEVEL_INFO {
		return level_name_info
	} else if self.Level == LEVEL_WARN {
		return level_name_warn
	} else if self.Level == LEVEL_ERROR {
		return level_name_error
	} else if self.Level == LEVEL_PANIC {
		return level_name_panic
	} else if self.Level == LEVEL_FATAL {
		return level_name_fatal
	} else {
		return "N/A"
	}
}

type AttrMap map[string]interface{}
