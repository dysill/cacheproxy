package cache

import (
	"time"
)

type node struct {
	key        string
	value      []byte
	expiresAt  time.Time
	prev, next *node
}

type dll struct {
	head, tail *node
}

func newDLL() *dll {
	head := &node{}
	tail := &node{}
	head.next = tail
	tail.prev = head
	return &dll{head: head, tail: tail}
}

func (l *dll) insertFront(n *node) {
	n.prev = l.head
	n.next = l.head.next
	l.head.next.prev = n
	l.head.next = n
}

func (l *dll) remove(n *node) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

func (l *dll) moveToFront(n *node) {
	l.remove(n)
	l.insertFront(n)
}

func (l *dll) removeBack() *node {
	last := l.tail.prev
	if last == l.head {
		return nil
	}
	l.remove(last)
	return last
}
