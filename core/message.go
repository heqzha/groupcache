package core

import (
	"sync"

	"github.com/heqzha/goutils/container"
)

type MessageQueue struct {
	q    map[string]*container.Queue
	sync *sync.Mutex
}

func (m *MessageQueue) Init() {
	m.q = make(map[string]*container.Queue)
	m.sync = &sync.Mutex{}
}

func (m *MessageQueue) Push(chn string, msg interface{}) {
	_, ok := m.q[chn]
	if !ok {
		m.q[chn] = new(container.Queue)
	}
	m.q[chn].Push(msg)
}

func (m *MessageQueue) Pop(chn string) interface{} {
	c, ok := m.q[chn]
	if ok && c != nil {
		return c.Pop()
	}
	return nil
}

func (m *MessageQueue) Len(chn string) int {
	return m.q[chn].Len()
}
