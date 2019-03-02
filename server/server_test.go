package server

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("running")
	s := NewMessagingServer("127.0.0.1:8000")
	go s.Run()
	os.Exit(m.Run())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n uint) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createConnections(t *testing.T) []*bufio.ReadWriter {
	t.Helper()
	var connections []*bufio.ReadWriter
	for i := 0; i < 2; i++ {
		newConn, err := net.Dial("tcp", "127.0.0.1:8000")
		if err != nil {
			t.Fatal(err)
		}
		rw := bufio.NewReadWriter(bufio.NewReader(newConn), bufio.NewWriter(newConn))
		_, err = rw.ReadString('\n')
		if err != nil {
			t.Fatal(err)
		}
		rw.WriteString(randomString(5))
		rw.WriteString("\n")
		rw.Flush()
		connections = append(connections, rw)
	}
	return connections
}

func Test_HappyPath(t *testing.T) {
	conns := createConnections(t)
	conns[0].WriteString("hello\n")
	conns[0].Flush()
	data, err := conns[1].ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if strings.Compare(data, "hello\n") != 0 {
		t.Fail()
	}
}
