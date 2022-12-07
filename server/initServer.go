package main

import (
	"changeme/model"
	"fmt"
	"github.com/lesismal/nbio"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net"
	"sync"
	"text/template"
)

var (
	gopher             *nbio.Gopher
	udpConnToSuperNode *nbio.Conn
	udpMutex           = &sync.Mutex{}
)

func connectDB() {
	var err error
	globalConn, err = gorm.Open(sqlite.Open(cfg.Data), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	globalConn.Debug().AutoMigrate(&model.UserTable{})
}

func connectSuperNode() {
	gopher = nbio.NewGopher(nbio.Config{})

	gopher.OnData(func(c *nbio.Conn, data []byte) {
		fmt.Println(string(data))
	})
	gopher.OnClose(func(c *nbio.Conn, err error) {
		udpMutex.Lock()
		defer udpMutex.Unlock()
		gopher = nil
		udpConnToSuperNode = nil
	})

	err := gopher.Start()
	if err != nil {
		panic(err)
	}
	c, err := net.Dial("udp", cfg.SuperNodeAddr)
	if err != nil {
		panic(err)
	}
	udpConnToSuperNode, err = gopher.AddConn(c)
	if err != nil {
		panic(err)
	}
}

var (
	communityTemplate *template.Template
)

func initTemplate() {
	var err error
	communityTemplate, err = template.ParseFiles(cfg.CommunityTemplate)
	if err != nil {
		panic(err)
	}
}
