package server

import (
	"bufio"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
)

//INIT is the initial state of a client connection
const INIT = "INIT"

type clientState struct {
	state string
	name  string
	in    *bufio.Reader
}

func newClientState(in io.Reader) *clientState {
	res := new(clientState)
	res.state = INIT
	res.in = bufio.NewReader(in)
	return res
}

//Connection handling code
func handleConnection(c net.Conn, input, output chan string) {
	go readThread(c, output)
	go writeThread(c, input)
}

func readThread(c net.Conn, output chan string) {
	state := newClientState(c)
	for {
		switch state.state {
		case INIT:

		}
		data, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Error(err)
			break
		}
		log.WithFields(log.Fields{
			"connection": c.RemoteAddr,
			"message":    data,
		}).Info("Read Message")

		output <- data
	}
	c.Close()
}

func writeThread(c net.Conn, input chan string) {
	for msg := range input {
		log.WithFields(log.Fields{
			"connection": c.RemoteAddr(),
			"message":    msg,
		}).Info("Writing message")
		_, err := c.Write([]byte(msg))
		if err != nil {
			log.Error(err)
			break
		}
	}
	c.Close()
}

//StartServer starts a new messaging server
func StartServer(portNum string) {
	PORT := ":" + portNum
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.WithFields(log.Fields{
		"Port": PORT,
	}).Info("Server started")

	cw := newChannelWriter()
	go cw.start()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Error(err)
		}

		log.WithFields(log.Fields{
			"connection": c.RemoteAddr(),
		}).Info("New connection made")

		newConnCh := make(chan string)
		cw.addOutputChannel(newConnCh)
		go handleConnection(c, newConnCh, cw.in)
	}
}
