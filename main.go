package main

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go StarScheduleFromGlobalQueue(ctx, time.Second*5)
	go StartMockTaskConsume(ctx)
	go StartMonitorJob(ctx)

	// 模拟推送任务
	for i := 0; i < 10000; i++ {
		taskName := "task-" + uuid.NewString()
		_ = Push(taskName)
		time.Sleep(time.Millisecond * 5)
	}

	select {}
}
