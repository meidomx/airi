package service

import (
	"sync"
	"time"

	"github.com/meidomx/airi/model"
)

type Event struct {
	TaskKey   string
	Parameter string

	TriggerTime time.Time
}

var newTaskQueue chan *model.Task
var inflightTaskQueueMap *sync.Map

func AddTask(t *model.Task) {
	newTaskQueue <- t
}

func GetInflightChannel(taskKey string) <-chan *Event {
	ch := getTaskQueue(taskKey, inflightTaskQueueMap)
	if ch == nil {
		return nil
	}
	return ch
}
