package main

import "testing"

func resetTestState() {
	mutex.Lock()
	globalQueue = nil
	globalQueueTail = nil
	length = 0
	mutex.Unlock()

	node1.tasks = nil
	node2.tasks = nil
	node3.tasks = nil
}

func TestBatchPopFromGlobalQueue(t *testing.T) {
	resetTestState()

	BatchPushToGlobalQueue(
		&TaskNode{task: "task-1"},
		&TaskNode{task: "task-2"},
		&TaskNode{task: "task-3"},
	)

	items := BatchPopFromGlobalQueue(2)
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].task != "task-1" || items[1].task != "task-2" {
		t.Fatalf("unexpected FIFO order: %#v", []string{items[0].task, items[1].task})
	}
	if got := GetGlobalQueueLength(); got != 1 {
		t.Fatalf("expected queue length 1, got %d", got)
	}
}

func TestScheduleFromGlobalQueueDrainsQueuedTasks(t *testing.T) {
	resetTestState()

	BatchPushToGlobalQueue(
		&TaskNode{task: "task-1"},
		&TaskNode{task: "task-2"},
		&TaskNode{task: "task-3"},
	)

	if err := ScheduleFromGlobalQueue(); err != nil {
		t.Fatalf("unexpected schedule error: %v", err)
	}
	if got := GetGlobalQueueLength(); got != 0 {
		t.Fatalf("expected queue to be drained, got length %d", got)
	}
	if len(node1.tasks) != 3 {
		t.Fatalf("expected node1 to receive 3 tasks, got %d", len(node1.tasks))
	}
}

