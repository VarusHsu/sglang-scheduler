package main

const threshold = 50

func Push(tasks ...string) error {
	// select node

	taskNodes := make([]*TaskNode, 0, len(tasks))
	for _, task := range tasks {
		taskNodes = append(taskNodes, &TaskNode{
			task: task,
		})
	}

	BatchPushToGlobalQueue(taskNodes...)

	return ScheduleFromGlobalQueue()
}
