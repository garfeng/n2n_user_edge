package main

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type MessageSender struct {
	topic string
	ctx   context.Context
}

func NewMessageSender(ctx context.Context, topic string) *MessageSender {
	return &MessageSender{
		topic: topic,
		ctx:   ctx,
	}
}

const (
	EventMessage = "message"
)

func (m *MessageSender) Write(buff []byte) (n int, err error) {
	s := string(buff)
	Log.Info(fmt.Sprintf("[%s]", m.topic), s)

	runtime.EventsEmit(m.ctx, EventMessage, &Message{
		Topic:   m.topic,
		Message: s,
	})

	return len(buff), nil
}

type Message struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}
