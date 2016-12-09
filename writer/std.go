package writer

import (
	"os"
	"syscall"
)

type StandardOutputWriter struct {
	out *os.File
	prefix string
	flags int
	transformFunc TransformFunc
}

func NewStandardOutputWriter(prefix string, flags int) *StandardOutputWriter {
	w := new(StandardOutputWriter)
	w.out = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")
	w.prefix = prefix
	w.flags = flags
	w.transformFunc = StandardTextLineTransform
	return w
}

func NewStandardOutputWriterWithTransform(prefix string, flags int, transformFunc TransformFunc) *StandardOutputWriter {
	w := new(StandardOutputWriter)
	w.out = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")
	w.prefix = prefix
	w.flags = flags
	w.transformFunc = transformFunc
	return w
}

func (w *StandardOutputWriter) Write(p []byte) (int, error) {
	msg, decodeErr := LogMsgDecode(p, w.prefix, w.flags)
	if decodeErr != nil {
		panic(decodeErr)
		return 0, decodeErr
	}
	bytes, transformErr := w.transformFunc(msg)
	if transformErr != nil {
		panic(transformErr)
		return 0, transformErr
	}
	if msg.NewLineFlag {
		bytes = append(bytes, '\n')
	}
	return w.out.Write(bytes)
}

