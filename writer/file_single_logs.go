package writer

import (
	"os"
	"sync"
	"bufio"
	"time"
	"fmt"
)


type LogsSingleFileWriter struct {
	level Level
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

func NewLogsSingleFileWriter(fileName string, level Level, logPrefix string, logFlags int) *LogsSingleFileWriter {
	w := LogsSingleFileWriter{
		level:level,
		logFlags:logFlags,
		logPrefix:logPrefix,
		fileName:fileName,
		fileOpenFlg:false,
		transform:LogsStandardTextLineTransform,
		wg:new(sync.WaitGroup),
		processChan:make(chan int, 1024),
	}
	w.closeFileWhenIdle()
	return &w
}

func NewLogsSingleFileWriterWithTransform(fileName string, level Level, logPrefix string, logFlags int, transform TransformFunc) *LogsSingleFileWriter {
	w := LogsSingleFileWriter{
		level:level,
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


func (w *LogsSingleFileWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	msg, decodeErr := LogMsgDecode(p, w.logPrefix, w.logFlags)
	if decodeErr != nil {
		panic(decodeErr)
		return 0, decodeErr
	}
	levelParseErr := parseLevelFromMsg(&msg)
	if levelParseErr != nil {
		msg.LevelNo = Level(10)
		msg.LevelName = "???"
	}
	if !levelGate(msg.LevelNo, w.level) {
		return 0, nil
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

func (w *LogsSingleFileWriter) Close() error {
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

func (w *LogsSingleFileWriter) closeFileWhenIdle() {
	go func(w *LogsSingleFileWriter) {
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
