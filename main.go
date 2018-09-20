package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	statusOk = iota
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Count   int32  `json:"count"`
}

type Counter struct {
	m           sync.Mutex
	Count       int32
	OpenSession int32
}

func (c *Counter) increment() {
	c.m.Lock()
	defer c.m.Unlock()
	c.Count++
	c.OpenSession++
}

func (c *Counter) decrement() {
	c.m.Lock()
	defer c.m.Unlock()
	c.OpenSession--
}

func (c Counter) getOpenSession() int32 {
	return c.OpenSession - 1
}

var conter Counter

func main() {
	fmt.Println("Launching server...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8000")

	for {
		// accept connection on port
		conn, err := ln.Accept()
		if nil != err {
			fmt.Println(err)
		}

		go handleConnection(conn)

	}

	var input string
	fmt.Scanln(&input)
}

func handleConnection(c net.Conn) {
	defer c.Close()
	conter.increment()
	defer conter.decrement()

	c.SetDeadline(time.Now().Add(100 * time.Millisecond))
	c.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	c.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))

	// will listen for message to process ending in newline (\n)
	_, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		return
		log.Fatal(err)
	}

	jsonData, err := json.Marshal(Response{
		Status:  statusOk,
		Message: "Cool",
		Count:   conter.getOpenSession(),
	})

	if err != nil {
		return
		log.Fatal(err)
	}

	// send new string back to client
	c.Write(jsonData)

}
