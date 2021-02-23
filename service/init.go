package service

import (
	"sync"

	"github.com/meidomx/airi/model"
)

func Init() {
	newTaskQueue = make(chan *model.Task, 1024)

	inflightTaskQueueMap = new(sync.Map)
}

func getOrCreateTaskQueue(taskKey string, m *sync.Map) chan *Event {
	i, ok := m.Load(taskKey)
	if ok {
		return i.(chan *Event)
	}
	i, _ = m.LoadOrStore(taskKey, make(chan *Event)) // use no cache channel by default
	return i.(chan *Event)
}

func getTaskQueue(taskKey string, m *sync.Map) chan *Event {
	i, ok := m.Load(taskKey)
	if ok {
		return i.(chan *Event)
	} else {
		return nil
	}
}
