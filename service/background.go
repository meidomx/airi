package service

import (
	"errors"
	"log"
	"time"

	"github.com/meidomx/airi/model"

	"github.com/better-concurrent/guc"
)

func StartBackgroundJob() {
	ts := loadDbTasks()
	startJob(ts)
}

type ScheduleTask struct {
	*model.Task
	NextTime int
}

func (s *ScheduleTask) CompareTo(i interface{}) int {
	return s.NextTime - i.(*ScheduleTask).NextTime
}

func startJob(ts []*model.Task) {
	pq := guc.NewPriority()
	total := 0
	for _, v := range ts {
		st := genScheduleTask(v)
		pq.Add(st)
		total++
		// init queue
		getOrCreateTaskQueue(v.TaskKey, inflightTaskQueueMap)
	}
	log.Println("[INFO] init job count:", total)

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			t := <-ticker.C
			log.Println("[INFO] trigger at:", t)

			// process add task
		LOOP:
			for {
				select {
				case t := <-newTaskQueue:
					st := genScheduleTask(t)
					pq.Add(st)
					// init queue
					getOrCreateTaskQueue(t.TaskKey, inflightTaskQueueMap)
				default:
					break LOOP
				}
			}

			// process task
			for {
				now := time.Now()
				tu := int(now.Unix())
				if pq.Peek() == nil {
					break
				}
				if pq.Peek().(*ScheduleTask).NextTime > tu {
					break
				}
				st := pq.Poll().(*ScheduleTask)
				inflightTaskQueue := getOrCreateTaskQueue(st.TaskKey, inflightTaskQueueMap)
				select {
				case inflightTaskQueue <- &Event{
					TaskKey:     st.TaskKey,
					Parameter:   "",
					TriggerTime: now,
				}:
				default:
				}
				pq.Add(genScheduleTask(st.Task))
			}
		}
	}()
}

func genScheduleTask(v *model.Task) *ScheduleTask {
	st := new(ScheduleTask)
	st.Task = v
	switch v.Type {
	case model.TaskTypeSimple:
		st.NextTime = CalculateNextTimeSimple(FromConfigToSimple(v.Config))
	default:
		panic(errors.New("unknown task type"))
	}
	return st
}

func loadDbTasks() []*model.Task {
	r, err := model.LoadAllAvailableTasks()
	if err != nil {
		panic(err)
	}
	return r
}
