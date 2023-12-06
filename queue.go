package delayqueue

import (
	"sync"
	"time"
)

const defSlotNum = 3600

type LinkList[T any] struct {
	Head *Entry[T]
	Tail *Entry[T]
	sync.Mutex
}

type Entry[T any] struct {
	Body     T
	Delay    time.Duration
	cycleNum int
	next     *Entry[T]
	prv      *Entry[T]
}

type Queue[T any] struct {
	consume     HandleFun[T]
	elements    []*LinkList[T]
	curPosition int
}

type HandleFun[T any] func(ele *Entry[T])

func NewQueue[T any](fun HandleFun[T], slotNum int) *Queue[T] {
	if slotNum == 0 {
		slotNum = defSlotNum
	}
	q := &Queue[T]{
		consume:  fun,
		elements: make([]*LinkList[T], slotNum, slotNum),
	}
	for i := 0; i < len(q.elements); i++ {
		q.elements[i] = &LinkList[T]{}
	}
	return q
}
