package server

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

type channelWriter struct {
	ch  []chan string
	In  chan string
	mux sync.Mutex
}

func NewChannelWriter() *channelWriter {
	cw := new(channelWriter)
	cw.In = make(chan string)
	return cw
}

func (cw *channelWriter) AddOutputChannel(out chan string) {
	cw.mux.Lock()
	defer cw.mux.Unlock()
	cw.ch = append(cw.ch, out)
}

func (cw *channelWriter) Start() {
	log.Info("Channel writer started")
	for msg := range cw.In {
		cw.mux.Lock()
		for _, out := range cw.ch {
			out <- msg
		}
		cw.mux.Unlock()
	}
}
