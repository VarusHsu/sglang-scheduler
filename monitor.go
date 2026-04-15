package main

import (
	"context"
	"log"
	"time"
)

func StartMonitorJob(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Global queue length", length, "node1 length", len(node1.tasks), "node2 length", len(node2.tasks), "node3 length", len(node3.tasks))
		case <-ctx.Done():
			return
		}
	}
}
