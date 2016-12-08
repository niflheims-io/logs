package writer

import (
	"time"
	"errors"
	"log"
	"bytes"
	"io"
	"strings"
)

type LogMsg struct {
	Prefix string
	DateTime time.Time
	File string
	Msg string
	NewLineFlag bool
	ExtendFields map[string]string
	LevelNo Level
	LevelName string
}

// choose log.flag
func LogMsgDecode(p []byte, prefix string, logFlags int) (msg LogMsg, err error) {
	if p == nil || len(p) == 0 {
		err = errors.New("Log message is empty.")
		return
	}
	if p[len(p)-1] == '\n' {
		p = p[0:len(p)-1]
		msg.NewLineFlag = true
	}
	readTimes := 0
	hasPrefix := prefix != ""
	if hasPrefix {
		msg.Prefix = prefix
		p = p[len(prefix):]
	}
	hasFile := logFlags & (log.Lshortfile|log.Llongfile) != 0
	if hasFile {
		readTimes = readTimes + 1
	}
	hasDateTime := logFlags & (log.Ldate|log.Ltime|log.Lmicroseconds) != 0
	hasDate := logFlags & log.Ldate != 0
	hasTime := logFlags & (log.Ltime|log.Lmicroseconds) != 0
	hasDateTimeUTC := false
	dateTimeLayout := ""
	if hasDateTime {
		hasDateTimeUTC = logFlags & log.LUTC != 0
		if hasDate {
			dateTimeLayout = "2006/01/02"
			readTimes = readTimes + 1
		}
		if hasTime {
			readTimes = readTimes + 1
			if hasDate {
				dateTimeLayout = dateTimeLayout + " 15:04:05"
			} else {
				dateTimeLayout = dateTimeLayout + "15:04:05"
			}
			if logFlags & (log.Lmicroseconds) != 0 {
				dateTimeLayout = dateTimeLayout + ".999999999"
			}
		}
	}

	if hasDateTime == false && hasFile == false {
		msg.Msg = string(p)
		return
	}
	lines := make([]string, 0, 1)
	bytesBuf := bytes.NewBuffer(p)
	line, readLineErr := bytesBuf.ReadString(' ')
	line = strings.TrimSpace(line)
	if readLineErr != nil {
		if readLineErr == io.EOF {
			line = bytesBuf.String()
		} else {
			err = readLineErr
			return
		}
	}
	lines = append(lines, line)
	for i := 1 ; i < readTimes ; i ++ {
		line, readLineErr = bytesBuf.ReadString(' ')
		line = strings.TrimSpace(line)
		if readLineErr != nil {
			if readLineErr == io.EOF {
				break
			} else {
				err = readLineErr
				return
			}
		}
		lines = append(lines, line)
	}
	lines = append(lines, bytesBuf.String())
	idx := 0
	if hasDateTime {
		dateTimeValue := ""
		if hasDate {
			dateTimeValue = dateTimeValue + lines[idx]
			idx = idx + 1
		}
		if hasTime {
			if hasDate {
				dateTimeValue = dateTimeValue + " " + lines[idx]
			} else {
				dateTimeValue = dateTimeValue + lines[idx]
			}
			idx = idx + 1
		}
		var dateTime time.Time
		var dateTimeErr error
		if hasDateTimeUTC {
			dateTime, dateTimeErr = time.ParseInLocation(dateTimeLayout, dateTimeValue, time.UTC)
		} else {
			dateTime, dateTimeErr = time.ParseInLocation(dateTimeLayout, dateTimeValue, time.Local)
		}
		if dateTimeErr != nil {
			err = dateTimeErr
			return
		}
		msg.DateTime = dateTime
	}
	if hasFile {
		msg.File = lines[idx]
	}
	msg.Msg = lines[len(lines) - 1]
	err = nil
	return
}