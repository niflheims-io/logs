package logs

import (
	"encoding/json"
	"fmt"
)

type jsonFormatter struct  {

}

func NewJsonFormatter() jsonFormatter {
	return jsonFormatter{}
}

func (self jsonFormatter) Format(dataPkg *DataPkg) ([]byte, error)  {
	outputMap := make(map[string]interface{})
	outputMap["level"] = dataPkg.levelName()
	outputMap["time"] = dataPkg.Time.String()
	outputMap["msg"] = dataPkg.Message
	if dataPkg.Attr != nil {
		for k, v := range dataPkg.Attr {
			outputMap[k] = fmt.Sprint(v)
		}
	}
	if dataPkg.Position != "" {
		outputMap["line"] = dataPkg.Position
	}
	return json.Marshal(&outputMap)
}
