package server

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	s := NewMessagingServer("127.0.0.1:8000")
	os.Exit(m.Run())
	s.Close()
}
