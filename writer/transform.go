package writer

type TransformFunc func(LogMsg) ([]byte, error)

func StandardTextLineTransform(msg LogMsg) ([]byte, error) {
	line := ""
	if msg.Prefix != "" {
		line = line + msg.Prefix + " "
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

func LogsStandardTextLineTransform(msg LogMsg) ([]byte, error) {
	line := ""
	if msg.Prefix != "" {
		line = line + msg.Prefix + " "
	}
	line = line + getLevelColor(msg.LevelNo) + "[" + msg.LevelName + "]" + color_default + " "
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

