package logs

import (
	"errors"
	"fmt"
	"strings"
)

type textLineFormatter struct  {

}

func NewTextLineFormatter() textLineFormatter {
	return textLineFormatter{}
}

func (self textLineFormatter) Format(dataPkg *DataPkg) ([]byte, error)  {
	if dataPkg == nil {
		return nil, errors.New("dataPkg is nil")
	}
	textLine := ""
	if dataPkg.Level > 0 {
		level := ""
		if dataPkg.Level == LEVEL_DEBUG {
			level = color_blue + dataPkg.levelName() + color_default
		} else if dataPkg.Level == LEVEL_INFO {
			level = color_green + dataPkg.levelName() + color_default
		} else if dataPkg.Level == LEVEL_WARN {
			level = color_yello + dataPkg.levelName() + color_default
		} else if dataPkg.Level == LEVEL_ERROR {
			level = color_red + dataPkg.levelName() + color_default
		} else if dataPkg.Level == LEVEL_PANIC {
			level = color_red + dataPkg.levelName() + color_default
		} else if dataPkg.Level == LEVEL_FATAL {
			level = color_red + dataPkg.levelName() + color_default
		}
		textLine = textLine + level + " | "
	}
	if !dataPkg.Time.IsZero() {
		textLine = textLine + color_yello + dataPkg.Time.String() + color_default + " | "
	}
	textLine = textLine + dataPkg.Message + " | "
	if dataPkg.Attr != nil {
		attrLine := ""
		for k, v := range dataPkg.Attr {
			key := color_attr + k + color_default
			val := fmt.Sprint(v)
			attrLine = attrLine + key + "=" + val + ";"
		}
		attrLine = strings.TrimSpace(attrLine)
		attrLine = attrLine[0:len(attrLine) - 1]
		if attrLine != "" {
			textLine = textLine + attrLine + " | "
		}
	}
	if dataPkg.Position != "" {
		textLine = textLine + dataPkg.Position + " "
	}
	return []byte(textLine), nil
}
