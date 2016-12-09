package writer

import (
	"os"
	"sync"
	"bufio"
	"time"
	"fmt"
	"strings"
)


type DailyFileWriter struct {
	logFileDir string
	transform TransformFunc
	logPrefix string
	logFlags int
	fileOpenFlg bool
	file *os.File
	fileWriterBuf *bufio.Writer
	wg *sync.WaitGroup
	processChan chan int
	today string
}

func NewDailyFileWriter(logFileDir string) *DailyFileWriter {
	w := DailyFileWriter{
		logFileDir:logFileDir,
		transform:StandardTextLineTransform,
		fileOpenFlg:false,
		wg:new(sync.WaitGroup),
		processChan:make(chan int, 1024),
	}
	if strings.LastIndex(w.logFileDir, "/") < len(w.logFileDir) - 1 {
		w.logFileDir = w.logFileDir + "/"
	}
	w.today = time.Now().Format("2006-01-02")
	w.closeFileWhenIdle()
	return &w
}

func NewDailyFileWriterWithTransform(logFileDir string, logPrefix string, logFlags int, transform TransformFunc) *DailyFileWriter {
	w := DailyFileWriter{
		logFileDir:logFileDir,
		transform:transform,
		logFlags:logFlags,
		logPrefix:logPrefix,
		fileOpenFlg:false,
		wg:new(sync.WaitGroup),
		processChan:make(chan int, 1024),
	}
	if strings.LastIndex(w.logFileDir, "/") < len(w.logFileDir) - 1 {
		w.logFileDir = w.logFileDir + "/"
	}
	w.today = time.Now().Format("2006-01-02")
	w.closeFileWhenIdle()
	return &w
}


func (w *DailyFileWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	msg, decodeErr := LogMsgDecode(p, w.logPrefix, w.logFlags)
	if decodeErr != nil {
		panic(decodeErr)
		return 0, decodeErr
	}
	msgDateTime := msg.DateTime.Format("2006-01-02")
	if msgDateTime != w.today {
		w.wg.Add(1)
		if w.fileOpenFlg {
			w.fileWriterBuf.Flush()
			syncErr := w.file.Sync()
			closeErr := w.file.Close()
			w.fileOpenFlg = false
			w.wg.Done()
			if syncErr != nil || closeErr != nil {
				panic(fmt.Sprintf("Sync file error:%v, close file error:%v.", &syncErr, &closeErr))
			}
		}
		w.today = msgDateTime
		w.wg.Done()
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
		w.file, fileOpenErr = os.OpenFile(w.logFileDir + w.today + ".log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
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

func (w *DailyFileWriter) Close() error {
	close(w.processChan)
	if !w.fileOpenFlg {
		return nil
	}
	if w.fileWriterBuf.Buffered() > 0 {
		w.fileWriterBuf.Flush()
	}
	syncErr := w.file.Sync()
	closeErr := w.file.Close()
	w.fileOpenFlg = false
	if syncErr != nil {
		return syncErr
	}
	return closeErr
}

func (w *DailyFileWriter) closeFileWhenIdle() {
	go func(w *DailyFileWriter) {
		for {
			select {
			case _, ok := <- w.processChan :
				if !ok {
					return
				}
			case <- time.After(time.Second * 5) :
				w.wg.Add(1)
				defer w.wg.Done()
				if w.fileOpenFlg {
					if w.fileWriterBuf.Buffered() > 0 {
						w.fileWriterBuf.Flush()
					}
					syncErr := w.file.Sync()
					closeErr := w.file.Close()
					w.fileOpenFlg = false
					if syncErr != nil || closeErr != nil {
						panic(fmt.Sprintf("Sync file error:%v, close file error:%v.", &syncErr, &closeErr))
					}
				}
			}
		}
	}(w)
}

