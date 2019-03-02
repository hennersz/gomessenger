package server

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

//A MessagingServer is a server that forwards messages from one client to all other clients
type MessagingServer struct {
	addr          string
	server        net.Listener
	idCounter     uint
	counterMux    sync.Mutex
	messageQueue  chan string
	connections   []*messenger
	connectionMux sync.RWMutex
	names         map[string]bool
	nameMapMux    sync.RWMutex
}

//NewMessagingServer creates a new MessagingServer that listens on the supplied address
func NewMessagingServer(addr string) *MessagingServer {
	newServer := new(MessagingServer)
	newServer.messageQueue = make(chan string)
	newServer.names = make(map[string]bool)
	newServer.addr = addr
	return newServer
}

type messenger struct {
	rw           *bufio.ReadWriter
	addr         string
	name         string
	id           uint
	mux          sync.RWMutex
	in           chan string
	out          chan string
	disconnected bool
}

//Run starts the MessagingServer listening for connections and forwarding messages
func (s *MessagingServer) Run() (err error) {
	s.server, err = net.Listen("tcp", s.addr)
	if err != nil {
		return
	}
	defer s.Close()

	if strings.Compare(s.server.Addr().String(), s.addr) != 0 {
		s.addr = s.server.Addr().String()
	}

	log.WithFields(log.Fields{
		"Address": s.addr,
	}).Info("Server started")
	go s.startMessageWriter()
	return s.handleConnections()
}

//Close stops the MessagingServer
func (s *MessagingServer) Close() (err error) {
	return s.server.Close()
}

func (s *MessagingServer) startMessageWriter() {
	for msg := range s.messageQueue {
		s.connectionMux.RLock()
		for _, c := range s.connections {
			c.in <- msg
		}
		s.connectionMux.RUnlock()
	}
}

func (s *MessagingServer) handleConnections() (err error) {
	for {
		conn, err := s.server.Accept()
		if err != nil || conn == nil {
			err = errors.New("could not accept connection")
			break
		}

		go s.handleConnection(conn)
	}
	return
}

func (s *MessagingServer) handleConnection(c net.Conn) {
	newConn, err := s.initialiseConnection(c)
	if err != nil {
		log.Error(err)
		return
	}
	s.addConnection(newConn)
	log.WithFields(log.Fields{
		"Address": newConn.addr,
		"Name":    newConn.name,
		"ID":      newConn.id,
	}).Info("New client connected")
	newConn.start()
}

func (s *MessagingServer) initialiseConnection(c net.Conn) (*messenger, error) {
	newConn := new(messenger)
	newConn.rw = bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
	newConn.addr = c.RemoteAddr().String()
	newConn.id = s.getID()
	newConn.out = s.messageQueue
	newConn.in = make(chan string)
	err := newConn.getNameFromUser(s)
	if err != nil {
		return nil, err
	}
	return newConn, nil
}

func (s *MessagingServer) addConnection(c *messenger) {
	s.connectionMux.Lock()
	s.connections = append(s.connections, c)
	s.connectionMux.Unlock()
}

func (s *MessagingServer) getID() (id uint) {
	s.counterMux.Lock()
	defer s.counterMux.Unlock()
	id = s.idCounter
	s.idCounter++
	return
}

func (s *MessagingServer) addName(name string) bool {
	s.nameMapMux.RLock()
	taken, ok := s.names[name]
	s.nameMapMux.RUnlock()
	if taken && ok {
		return false
	}
	s.nameMapMux.Lock()
	s.names[name] = true
	s.nameMapMux.Unlock()
	return true
}

func (m *messenger) getNameFromUser(s *MessagingServer) error {
	m.rw.WriteString("Please enter a name\n")
	m.rw.Flush()
	nameSet := false
	for !nameSet {
		name, err := m.rw.ReadString('\n')
		name = strings.TrimSuffix(name, "\n")
		if err != nil {
			return err
		}

		if s.addName(name) {
			m.name = name
			nameSet = true
		} else {
			m.rw.WriteString("Name taken, please choose another\n")
			m.rw.Flush()
		}
	}
	return nil
}

func (m *messenger) start() {
	go m.readMessages()
	go m.writeMessages()
}

func (m *messenger) readMessages() {
	for {
		data, err := m.rw.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Error(err)
			}
			log.WithFields(log.Fields{
				"connection": m.addr,
				"id":         m.id,
			}).Info("Connection closed")
			break
		}
		log.WithFields(log.Fields{
			"connection": m.addr,
			"message":    data,
		}).Info("Read Message")

		m.out <- data
	}
}

func (m *messenger) writeMessages() {
	for msg := range m.in {
		m.rw.WriteString(msg)
		m.rw.Flush()
		log.WithFields(log.Fields{
			"connection": m.addr,
			"message":    msg,
		}).Info("Wrote message")
	}
}
