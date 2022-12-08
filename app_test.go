package main

import (
	"testing"
	"time"
)

func TestApp_SetupN2N(t *testing.T) {
	a := &App{}
	a.SetupN2N()

	<-time.After(time.Second * 5)
	a.ShutdownN2N()
}
