package main

import (
	"context"
	"errors"
	"log"
	"sort"
	"sync"
	"time"
)

type TaskNode struct {
	next *TaskNode
	task string
}

var (
	globalQueue     *TaskNode
	globalQueueTail *TaskNode
	mutex           sync.Mutex
	length          int
)

func PushToGlobalQueue(nodes ...*TaskNode) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, node := range nodes {
		if node == nil {
			continue
		}
		node.next = nil
		if globalQueue == nil {
			globalQueue = node
			globalQueueTail = node
		} else {
			globalQueueTail.next = node
			globalQueueTail = node
		}
		length++
	}
}

func BatchPushToGlobalQueue(nodes ...*TaskNode) {
	PushToGlobalQueue(nodes...)
}

func BatchPopFromGlobalQueue(n int) []*TaskNode {
	mutex.Lock()
	defer mutex.Unlock()
	return batchPopFromGlobalQueueLocked(n)
}

func batchPopFromGlobalQueueLocked(n int) []*TaskNode {

	if n <= 0 || globalQueue == nil {
		return nil
	}

	if n > length {
		n = length
	}

	items := make([]*TaskNode, 0, n)
	for i := 0; i < n && globalQueue != nil; i++ {
		current := globalQueue
		globalQueue = current.next
		current.next = nil
		items = append(items, current)
		length--
	}

	if globalQueue == nil {
		globalQueueTail = nil
	}

	return items
}

func GetGlobalQueueLength() int {
	mutex.Lock()
	defer mutex.Unlock()
	return length
}

func ScheduleFromGlobalQueue() error {
	if ok := mutex.TryLock(); !ok {
		return nil
	}
	defer mutex.Unlock()

	nodeStatus := GetNodesStatus()

	if len(nodeStatus) == 0 {
		return errors.New("no available nodes")
	}

	sort.Slice(
		nodeStatus,
		func(i, j int) bool {
			return len(nodeStatus[i].tasks) < len(nodeStatus[j].tasks)
		},
	)

	node := nodeStatus[0]
	scheduleSize := threshold - len(node.tasks)
	if scheduleSize <= 0 {
		return nil
	}
	scheduleTasks := batchPopFromGlobalQueueLocked(scheduleSize)
	for _, task := range scheduleTasks {
		node.tasks = append(node.tasks, task.task)
	}
	return nil
}

func StarScheduleFromGlobalQueue(ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = ScheduleFromGlobalQueue()
			log.Printf("scheduled tasks from global queue, current length: %d", GetGlobalQueueLength())
		case <-ctx.Done():
			return
		}
	}

}
