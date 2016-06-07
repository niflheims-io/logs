package writer

import (
	"fmt"
)

type stdWriter struct  {
	async bool
}

func NewStdWriter(async bool) *stdWriter {
	return &stdWriter{async:async}
}

func (self stdWriter) Write(p []byte) (int, error) {
	if self.async {
		go func(p []byte) {
			fmt.Println(string(p))
		}(p)
		return 0, nil
	} else {
		return fmt.Println(string(p))
	}
}
