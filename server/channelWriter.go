package server

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

type channelWriter struct {
	ch  []chan string
	in  chan string
	mux sync.Mutex
}

func newChannelWriter() *channelWriter {
	cw := new(channelWriter)
	cw.in = make(chan string)
	return cw
}

func (cw *channelWriter) addOutputChannel(out chan string) {
	cw.mux.Lock()
	defer cw.mux.Unlock()
	cw.ch = append(cw.ch, out)
}

func (cw *channelWriter) start() {
	log.Info("Channel writer started")
	for msg := range cw.in {
		cw.mux.Lock()
		for _, out := range cw.ch {
			out <- msg
		}
		cw.mux.Unlock()
	}
}
