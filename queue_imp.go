package delayqueue

import (
	"fmt"
	"runtime/debug"
	"time"
)

func (q *Queue[T]) Push(delay time.Duration, ele T) {
	slotOffset := int(delay.Seconds())
	if slotOffset <= 0 {
		return
	}
	// 偏移位置需要减掉1，因为计时器不能立即执行
	slotOffset = slotOffset - 1
	absPosition := slotOffset + q.curPosition
	position := absPosition % len(q.elements)
	cycleNum := absPosition / len(q.elements)
	entry := &Entry[T]{
		Body:     ele,
		Delay:    delay,
		cycleNum: cycleNum,
	}
	elementList := q.elements[position]
	elementList.Lock()
	defer elementList.Unlock()
	if elementList.Head == nil {
		elementList.Head = entry
		elementList.Tail = entry
		return
	}
	if q.elements[position].Tail != nil {
		elementList.Tail.next = entry
		entry.prv = elementList.Tail
		elementList.Tail = entry
	}
	return
}

func (q *Queue[T]) asyncDeal(node *Entry[T]) {
	q.async(func() {
		q.consume(node)
	})
}

func (q *Queue[T]) deal(ll *LinkList[T]) {
	ll.Lock()
	defer ll.Unlock()
	curNode := ll.Head
	for curNode != nil {
		if curNode.cycleNum <= 0 {
			q.asyncDeal(curNode)
			ll.Remove(curNode)
		} else {
			curNode.cycleNum--
		}
		curNode = curNode.next
	}
	return
}
func (ll *LinkList[T]) Remove(node *Entry[T]) {
	if node.prv == nil && node.next == nil {
		ll.Head = nil
		ll.Tail = nil
	}
	if node.prv == nil && node.next != nil {
		nextNode := node.next
		ll.Head = nextNode
		nextNode.prv = nil
	}
	if node.prv != nil && node.next == nil {
		prvNode := node.prv
		ll.Tail = prvNode
		prvNode.next = nil
	}
	if node.prv != nil && node.next != nil {
		prvNode := node.prv
		nextNode := node.next
		prvNode.next = nextNode
		nextNode.prv = prvNode
	}
}

func (q *Queue[T]) Run() {
	q.async(func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				if q.elements[q.curPosition].Head != nil {
					q.asyncDealLinkList(q.elements[q.curPosition])
				}
				if q.curPosition >= (len(q.elements) - 1) {
					q.curPosition = 0
				} else {
					q.curPosition++
				}
			}
		}
	})
}

func (q *Queue[T]) asyncDealLinkList(list *LinkList[T]) {
	q.async(func() {
		q.deal(list)
	})
}

func (q *Queue[T]) async(fun func()) {
	go func(f func()) {
		defer func() {
			if err := recover(); err != nil {
				errMsg := fmt.Sprintf("======== Panic ========\nPanic: %v\nTraceBack:\n%s\n======== Panic ========", err, string(debug.Stack()))
				fmt.Println(errMsg)
			}
		}()
		f()
	}(fun)
}
