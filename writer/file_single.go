package writer

import (
	"os"
	"sync"
	"bufio"
	"time"
	"fmt"
)


type SingleFileWriter struct {
	fileName string
	transform TransformFunc
	logPrefix string
	logFlags int
	fileOpenFlg bool
	file *os.File
	fileWriterBuf *bufio.Writer
	wg *sync.WaitGroup
	processChan chan int
}

func NewSingleFileWriter(fileName string) *SingleFileWriter {
	w := SingleFileWriter{
		fileName:fileName,
		transform:StandardTextLineTransform,
		fileOpenFlg:false,
		wg:new(sync.WaitGroup),
		processChan:make(chan int, 1024),
	}
	w.closeFileWhenIdle()
	return &w
}

func NewSingleFileWriterWithTransform(fileName string, logPrefix string, logFlags int, transform TransformFunc) *SingleFileWriter {
	w := SingleFileWriter{
		fileName:fileName,
		transform:transform,
		logFlags:logFlags,
		logPrefix:logPrefix,
		fileOpenFlg:false,
		wg:new(sync.WaitGroup),
		processChan:make(chan int, 1024),
	}
	w.closeFileWhenIdle()
	return &w
}


func (w *SingleFileWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	msg, decodeErr := LogMsgDecode(p, w.logPrefix, w.logFlags)
	if decodeErr != nil {
		panic(decodeErr)
		return 0, decodeErr
	}
	b, transformErr := w.transform(msg)
	if transformErr != nil {
		panic(transformErr)
		return 0, transformErr
	}
	p = b
	if msg.NewLineFlag {
		p = append(p, '\n')
	}
	w.wg.Wait()
	if !w.fileOpenFlg {
		w.wg.Add(1)
		var fileOpenErr error
		w.file, fileOpenErr = os.OpenFile(w.fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if fileOpenErr != nil {
			return 0, fileOpenErr
		}
		w.fileWriterBuf = bufio.NewWriter(w.file)
		w.fileOpenFlg = true
		w.wg.Done()
	}
	defer w.fileWriterBuf.Flush()
	w.processChan <- 1
	return w.fileWriterBuf.Write(p)
}

func (w *SingleFileWriter) Close() error {
	close(w.processChan)
	if !w.fileOpenFlg {
		return nil
	}
	if w.fileWriterBuf.Buffered() > 0 {
		w.fileWriterBuf.Flush()
	}
	syncErr := w.file.Sync()
	closeErr := w.file.Close()
	if syncErr != nil {
		return syncErr
	}
	return closeErr
}

func (w *SingleFileWriter) closeFileWhenIdle() {
	go func(w *SingleFileWriter) {
		for {
			select {
			case _, ok := <- w.processChan :
				if !ok {
					return
				}
			case <- time.After(time.Second * 5) :
				if w.fileOpenFlg {
					w.wg.Add(1)
					if w.fileWriterBuf.Buffered() > 0 {
						w.fileWriterBuf.Flush()
					}
					syncErr := w.file.Sync()
					closeErr := w.file.Close()
					w.fileOpenFlg = false
					w.wg.Done()
					if syncErr != nil || closeErr != nil {
						panic(fmt.Sprintf("Sync file error:%v, close file error:%v.", &syncErr, &closeErr))
					}
				}
			}
		}
	}(w)
}
