package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	Log *logrus.Logger
)

func init() {
	os.MkdirAll("./log", 0755)
	Log = logrus.New()
	w, err := os.Create("./log/log.txt")
	if err != nil {
		panic(err)
	}
	Log.SetOutput(w)
}
