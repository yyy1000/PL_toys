package ProtocolMultiplexer

import "sync"

type Mux struct {
	srv  Service
	send chan Message

	mu      sync.Mutex
	pending map[int64]chan Message
}

func (m *Mux) Init(service Service) {
	m.srv = service
	m.pending = make(map[int64]chan Message)
	go m.sendLoop()
	go m.recvLoop()
}

func (m *Mux) sendLoop() {
	for args := range m.send {
		m.srv.Send(args)
	}
}

func (m *Mux) recvLoop() {
	for {
		reply := m.srv.Recv()
		tag := m.srv.ReadTag(reply)

		m.mu.Lock()
		done := m.pending[tag]
		delete(m.pending, tag)
		m.mu.Unlock()

		if done == nil {
			panic("unexpected reply")
		}
		done <- reply
	}
}

func (m *Mux) Call(args Message) (reply Message) {
	tag := m.srv.ReadTag(args)
	done := make(chan Message, 1)

	m.mu.Lock()
	if m.pending[tag] != nil {
		m.mu.Unlock()
		panic("mux: duplicate call tag")
	}
	m.pending[tag] = done
	m.mu.Unlock()

	m.send <- args
	return <-done
}
