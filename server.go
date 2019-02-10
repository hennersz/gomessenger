package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

type conChanList struct {
	ch  []chan string
	mux sync.Mutex
}

func messageWriter(connections *conChanList, in chan string) {
	for msg := range in {
		fmt.Println("got message")
		connections.mux.Lock()
		for _, out := range connections.ch {
			out <- msg
		}
		connections.mux.Unlock()
	}
}

func handleConnection(c net.Conn, input, output chan string) {
	fmt.Println("handling connection")
	go readThread(c, output)
	go writeThread(c, input)
	fmt.Println("connection handled")
}

func readThread(c net.Conn, output chan string) {
	for {
		data, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("read data")

		output <- data
	}
	c.Close()
}

func writeThread(c net.Conn, input chan string) {
	for msg := range input {
		fmt.Println("writing")
		c.Write([]byte(msg))
	}
	c.Close()
}

func main() {
	arguments := os.Args

	if len(arguments) == 1 {
		fmt.Println("Please provide a port number")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	connections := new(conChanList)
	send := make(chan string)
	go messageWriter(connections, send)
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("New connection")
		newConnCh := make(chan string)
		connections.mux.Lock()
		connections.ch = append(connections.ch, newConnCh)
		connections.mux.Unlock()
		go handleConnection(c, newConnCh, send)
	}
}
