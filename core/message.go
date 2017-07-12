package core

import "github.com/heqzha/goutils/container"

type MessageQueue map[string]*container.Queue

func (m MessageQueue) Push(chn string, msg interface{}) {
	_, ok := m[chn]
	if !ok {
		m[chn] = new(container.Queue)
	}
	m[chn].Push(msg)
}

func (m MessageQueue) Pop(chn string) interface{} {
	c, ok := m[chn]
	if ok && c != nil {
		return c.Pop()
	}
	return nil
}
