package server

import (
	"bytes"
	"io"
	"net"
	"reflect"
	"testing"

	testLog "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func testConn(t *testing.T) (client, server net.Conn) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error(err)
	}
	go func() {
		defer ln.Close()
		server, err = ln.Accept()
		if err != nil {
			t.Error(err)
		}
	}()

	client, err = net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Error(err)
	}
	return
}

func Test_writeThread(t *testing.T) {
	tests := []struct {
		name  string
		input []string
	}{
		{
			"Normal Data",
			[]string{"a\n", "hellon\n"},
		},
		{
			"No new line",
			[]string{"a", "hellon"},
		},
		{
			"Non ascii",
			[]string{"✅℅ℬℶ∝∞‡"},
		},
	}
	for _, tt := range tests {
		inputChannel := make(chan string)
		clientConn, serverConn := testConn(t)
		t.Run(tt.name, func(t *testing.T) {
			hook := testLog.NewGlobal()
			go writeThread(serverConn, inputChannel)
			for _, data := range tt.input {
				inputChannel <- data
				buffer := make([]byte, len(data))
				n, err := clientConn.Read(buffer)
				if err != nil {
					t.Error(err)
				}
				if n != len(data) {
					t.Error("Wrong number of bytes written to connection")
				}
				if !bytes.Equal(buffer, []byte(data)) {
					t.Error("Incorrect data written to socket")
				}
			}
			assert.Equal(t, len(tt.input), len(hook.AllEntries()))
		})
		clientConn.Close()
		close(inputChannel)
	}
}

func Test_newClientState(t *testing.T) {
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name string
		args args
		want *clientState
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newClientState(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newClientState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleConnection(t *testing.T) {
	type args struct {
		c      net.Conn
		input  chan string
		output chan string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleConnection(tt.args.c, tt.args.input, tt.args.output)
		})
	}
}

func Test_readThread(t *testing.T) {
	type args struct {
		c      net.Conn
		output chan string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readThread(tt.args.c, tt.args.output)
		})
	}
}

func TestStartServer(t *testing.T) {
	type args struct {
		portNum string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartServer(tt.args.portNum)
		})
	}
}
