package delayqueue

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestQueueRun(t *testing.T) {
	now := time.Now()
	fun := func(msg *Entry[string]) {
		startTime := now
		fmt.Println("消费数据", msg, time.Now().Format(time.DateTime))
		assert.Equal(t, time.Now().Format(time.DateTime), startTime.Add(msg.Delay).Format(time.DateTime), "消费时间应该相同")
	}
	queue := NewQueue(fun, 10)
	queue.Run()
	queue.Push(time.Second*2, "2秒后消费")
	queue.Push(time.Second*5, "5秒后消费")
	queue.Push(time.Second*10, "10秒后消费")
	time.Sleep(time.Second * 15)
}

func TestExecQueueDuration(t *testing.T) {
	now := time.Now()
	fmt.Println("当前时间:", now.Format(time.DateTime))
	weelNum := 10
	testData := []time.Duration{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 2, 3, 3, 22, 3, 2, 5}
	fun := func(msg *Entry[string]) {
		duration := int(time.Now().Sub(now).Seconds())
		fmt.Printf("实际：%d秒后消费,期望：%s\n", duration, msg.Body)
		assert.Equal(t, msg.Body, fmt.Sprintf("%d秒后消费", duration), "时间环验证")
	}
	queue := NewQueue(fun, weelNum)
	queue.Run()
	for _, item := range testData {
		queue.Push(time.Second*item, fmt.Sprintf("%d秒后消费", item))
	}
	time.Sleep(time.Second * 30)
}

func TestQueue(t *testing.T) {
	weelNum := 3600
	testData := []time.Duration{2, 5, 10}
	fun := func(msg *Entry[string]) {}
	queue := NewQueue(fun, weelNum)
	for _, item := range testData {
		queue.Push(time.Second*item, fmt.Sprintf("%d秒后消费", item))
	}
Label:
	for _, item := range testData {
		list := queue.elements[(int(item)-1)%weelNum]
		cycleNum := (int(item) - 1) / weelNum
		curNode := list.Head
		for curNode != nil {
			if curNode.cycleNum == cycleNum && curNode.Delay == item*time.Second {
				continue Label
			}
		}
		assert.Failf(t, "环形数组中的值错误", "")
	}
}
