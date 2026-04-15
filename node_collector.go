package main

import (
	"context"
	"log"
	"math/rand"
	"time"
)

type PrefillNode struct {
	nodeName string
	tasks    []string
}

var (
	node1 = PrefillNode{nodeName: "node-a"}
	node2 = PrefillNode{nodeName: "node-b"}
	node3 = PrefillNode{nodeName: "node-c"}
)

func GetNodesStatus() []*PrefillNode {
	return []*PrefillNode{
		&node1,
		&node2,
		&node3,
	}
}

type NodeCollector interface {
	GetNodesStatus() []*PrefillNode
}

func StartMockTaskConsume(ctx context.Context) {
	run := func(node *PrefillNode) {
		n := rand.Intn(10)
		if n > len(node.tasks) {
			n = len(node.tasks)
		}

		for i := 0; i < n; i++ {
			log.Println(node.nodeName, "consuming task:", node.tasks[i])
		}
		node.tasks = node.tasks[n:]
	}

	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for _, node := range GetNodesStatus() {
				run(node)
			}
		case <-ctx.Done():
			return
		}
	}
}
