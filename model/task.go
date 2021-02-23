package model

import (
	"context"
	"time"

	"github.com/moetang/webapp-scaffold/frmpg"
	"github.com/moetang/webapp-scaffold/utils"
)

type TaskType int16
type TaskStatus int16

const (
	TaskTypeSimple TaskType = 1
)

const (
	TaskStatusNormal  TaskStatus = 0
	TaskStatusDeleted TaskStatus = 1
)

type Task struct {
	TaskId      int64      `mx.orm:"task_id"`
	TaskKey     string     `mx.orm:"task_key"`
	Description string     `mx.orm:"description"`
	Type        TaskType   `mx.orm:"type"`
	Status      TaskStatus `mx.orm:"status"`
	Config      string     `mx.orm:"config"`
	TimeCreated int64      `mx.orm:"time_created"`
	TimeUpdated int64      `mx.orm:"time_updated"`
}

func SaveTask(t *Task) error {
	now := utils.UnixTime(time.Now())
	t.TimeUpdated = now
	t.TimeCreated = now
	rs, err := db.GetPostgresPool().Query(context.Background(),
		"insert into airi_task(task_key, description, type, status, config, time_created, time_updated) values ($1, $2, $3, $4, $5, $6, $7) returning task_id",
		t.TaskKey, t.Description, t.Type, t.Status, t.Config, t.TimeCreated, t.TimeUpdated)
	if err != nil {
		return err
	}
	if rs.Next() {
		err = rs.Scan(&t.TaskId)
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return frmpg.ErrNoRecordFound
	}
}

func LoadAllAvailableTasks() ([]*Task, error) {
	var r []*Task
	err := frmpg.QueryMulti(db.GetPostgresPool(), &r, context.Background(),
		"select * from airi_task where status = $1",
		TaskStatusNormal)
	if err != nil {
		return nil, err
	}
	return r, nil
}
