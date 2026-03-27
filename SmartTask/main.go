package main

import (
	"container/heap"
	"fmt"
	"time"
)

type Task struct {
	ID        int
	Name      string
	Priority  int
	Retry     int
	ExecuteAt time.Time
}

type TaskManager struct {
	queue Queue
	pq    PriorityQueue
	isSeq int
}

func NewTaskManager() *TaskManager {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	return &TaskManager{
		queue: Queue{},
		pq:    pq,
	}
}

func (tm *TaskManager) AddTask(name string, priority int, delaySeconds int) {
	tm.isSeq++
	task := &Task{
		ID:        tm.isSeq,
		Name:      name,
		Priority:  priority,
		Retry:     0,
		ExecuteAt: time.Now().Add(time.Duration(delaySeconds) * time.Second),
	}
	if priority > 0 {
		heap.Push(&tm.pq, task)
	} else {
		tm.queue.Enqueue(task)
	}
}

func (tm *TaskManager) ProcessTask() {
	now := time.Now()

	if tm.pq.Len() > 0 {
		task := tm.pq[0]

		if task.ExecuteAt.Before(now) {
			task = heap.Pop(&tm.pq).(*Task)
			executeTask(tm, task)
			return
		}
	}
	task := tm.queue.Dequeue()
	if task != nil && task.ExecuteAt.Before(now) {
		executeTask(tm, task)
	}
}

func executeTask(tm *TaskManager, task *Task) {
	fmt.Printf("Processing Task: %s (Priority: %d)\n", task.Name, task.Priority)

	if task.Retry < 2 {
		fmt.Println("Task failed, retrying...")
		task.Retry++
		task.ExecuteAt = time.Now().Add(2 * time.Second)

		heap.Push(&tm.pq, task)
	} else {
		fmt.Println("Task Completed!")
	}
}

func main() {
	taskManager := NewTaskManager()
	taskManager.AddTask("Task 1", 1, 0)
	taskManager.AddTask("Task 2", 0, 0)
	taskManager.AddTask("Task 3", 2, 0)

	for {
		taskManager.ProcessTask()
		time.Sleep(1 * time.Second)
	}
}
