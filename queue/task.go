package queue

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"
)

// Task - Task Struct
type Task struct {
	Run      func()
	id       string
	metadata MetaData
}

type MetaData struct {
	Name     string
	Duration time.Duration
	Success  bool
	TaskType string
}

type BackgroundTask struct {
	Name string
	Task *Task
	job  int8
}

type ScheduledTask struct {
	Name string
	Task *Task
	job  int8
	time time.Duration
}

func (bt *BackgroundTask) GenerateMetaData(tt int8) {
	bt.Task.id = uuid()
	bt.Task.metadata.Name = bt.Name
	switch tt {
	case 0:
		bt.Task.metadata.TaskType = "Background"
	case 1:
		bt.Task.metadata.TaskType = "Scheduled"
	}

}

func (bt *ScheduledTask) GenerateMetaData(tt int8) {
	bt.Task.id = uuid()
	bt.Task.metadata.Name = bt.Name
	switch tt {
	case 0:
		bt.Task.metadata.TaskType = "Background"
	case 1:
		bt.Task.metadata.TaskType = "Scheduled"
	}
}

func uuid() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return (uuid)
}
