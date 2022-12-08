package main

import (
	"fmt"
	"github.com/lesismal/nbio"
	"net"
	"os"
)

var (
	conn   *nbio.Conn
	gopher *nbio.Gopher
)

func connect(addr string) {
	gopher = nbio.NewGopher(nbio.Config{})

	gopher.OnData(func(c *nbio.Conn, data []byte) {
		fmt.Println("Receive:", string(data))
	})
	gopher.OnClose(func(c *nbio.Conn, err error) {
		gopher = nil
		conn = nil
	})

	err := gopher.Start()
	if err != nil {
		panic(err)
	}
	c, err := net.Dial("udp", addr)
	if err != nil {
		panic(err)
	}
	conn, err = gopher.AddConn(c)
	if err != nil {
		panic(err)
	}
}

func main() {
	connect(os.Args[1])
	for {
		msg := ""
		fmt.Scanf("%s", &msg)
		if msg == "q" {
			break
		}
		conn.Write([]byte(msg))
	}
	conn.Close()
}
