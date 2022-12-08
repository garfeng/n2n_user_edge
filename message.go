package main

import (
	"errors"
	"fmt"
)

func NewMessageReceiver(length int) *MessageReceiver {
	que := make(chan *Message, length)
	return &MessageReceiver{que: &que}
}

type MessageSender struct {
	que   *chan *Message
	topic string
}

func (m *MessageSender) Write(buff []byte) (n int, err error) {
	if m.que == nil {
		return 0, errors.New("msg que is nil")
	}

	s := string(buff)
	Log.Info(fmt.Sprintf("[%s]", m.topic), s)

	defer func() {
		if r := recover(); r != nil {
			Log.Error("Panic:", r)
			m.que = nil
		}
	}()

	*m.que <- &Message{
		Topic:   m.topic,
		Message: s,
	}

	return len(buff), nil
}

type MessageReceiver struct {
	que *chan *Message
}

func (m *MessageReceiver) Read() *Message {
	return <-*m.que
}

func (m *MessageReceiver) Close() error {
	close(*m.que)
	return nil
}

func (m *MessageReceiver) NewSender(topic string) *MessageSender {
	return &MessageSender{
		que:   m.que,
		topic: topic,
	}
}

type Message struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}
