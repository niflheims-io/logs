package writer

import (
	"io"
	"sync"
)

type GroupWriter struct {
	outs []*writerNode
	wg *sync.WaitGroup
}

func NewGroupWriter(outs ...io.Writer) *GroupWriter {
	outLen := len(outs)
	if outLen == 0 {
		panic("no writers.")
	}
	wg := new(sync.WaitGroup)
	wns := make([]*writerNode,0, outLen)
	for i := 0 ; i < outLen ; i ++ {
		out := outs[i]
		wn := new(writerNode)
		wn.out = out
		wn.msgChan = make(chan msgPackage, 1024)
		wn.wg = wg
		wn.closeFlg = false
		wn.receive()
		wns = append(wns, wn)
	}
	gw := new(GroupWriter)
	gw.outs = wns
	gw.wg = wg
	return gw
}

type msgPackage struct {
	p []byte
	wg *sync.WaitGroup
}

type writerNode struct {
	out io.Writer
	msgChan chan msgPackage
	wg *sync.WaitGroup
	closeFlg bool
}

func (n *writerNode) close() {
	n.closeFlg = true
	close(n.msgChan)
}

func (n *writerNode) send(msg msgPackage) {
	if n.closeFlg {
		return
	}
	n.msgChan <- msg
}

func (n *writerNode) receive()  {
	n.wg.Add(1)
	go func(n *writerNode) {
		defer n.wg.Done()
		for {
			msg, ok := <- n.msgChan
			if !ok {
				break
			}
			_, flushedErr := n.out.Write(msg.p)
			msg.wg.Done()
			if flushedErr != nil {
				panic(flushedErr)
			}
		}
	}(n)
}

func (w *GroupWriter) Write(p []byte) (int, error) {
	outLen := len(w.outs)
	wg := new(sync.WaitGroup)
	for i := 0 ; i < outLen ; i ++ {
		wg.Add(1)
		w.outs[i].send(msgPackage{p:p, wg:wg})
	}
	wg.Wait()
	return len(p), nil
}

func (w *GroupWriter) Close() {
	outLen := len(w.outs)
	for i := 0 ; i < outLen ; i ++ {
		out := w.outs[i]
		out.close()
	}
	w.wg.Wait()
}