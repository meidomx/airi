package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/meidomx/airi/model"
	"github.com/meidomx/airi/service"

	scaffold "github.com/moetang/webapp-scaffold"

	"github.com/gin-gonic/gin"
)

const (
	RespStatusSuccess       = 0
	RespStatusErrParam      = 1
	RespStatusErrGeneral    = 2
	RespStatusListenTimeout = 3
)

func InitRestController(webscaf *scaffold.WebappScaffold) {
	g := webscaf.GetGin().Group("/api/v1")
	g.POST("/simple_task", CreateSimpleTask)
	g.GET("/simple_task/:task_key", GetTask)
}

type CreateSimpleTaskReq struct {
	TaskKey     string `json:"task_key" binding:"required"`
	Description string `json:"description" binding:"required"`
	Every       string `json:"every" binding:"required"`
	At          *int   `json:"at" binding:"required"`
}

func CreateSimpleTask(ctx *gin.Context) {
	var req = new(CreateSimpleTaskReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		log.Println("[ERROR] bind json req error.", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":       RespStatusErrParam,
			"errormessage": "parameter error",
		})
		return
	}

	var et service.EveryType
	switch req.Every {
	case "day":
		et = service.EveryDay
		if *req.At < 0 || *req.At > 23 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":       RespStatusErrParam,
				"errormessage": "parameter error",
			})
			return
		}
	case "hour":
		et = service.EveryHour
		if *req.At < 0 || *req.At > 59 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":       RespStatusErrParam,
				"errormessage": "parameter error",
			})
			return
		}
	case "minute":
		et = service.EveryMinute
		if *req.At < 0 || *req.At > 59 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":       RespStatusErrParam,
				"errormessage": "parameter error",
			})
			return
		}
	case "second":
		et = service.EverySecond
		*req.At = 0
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":       RespStatusErrParam,
			"errormessage": "parameter error",
		})
		return
	}
	s := new(service.SimpleConfig)
	s.At = *req.At
	s.Et = et

	t := new(model.Task)
	t.TaskKey = req.TaskKey
	t.Description = req.Description
	t.Type = model.TaskTypeSimple
	t.Status = model.TaskStatusNormal
	t.Config = service.ConvertSimpleToConfig(s)
	// save to db
	err := model.SaveTask(t)
	if err != nil {
		log.Println("[ERROR] save task error.", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":       RespStatusErrGeneral,
			"errormessage": "internal error",
		})
		return
	}

	// add to schedule
	service.AddTask(t)

	ctx.JSON(http.StatusOK, gin.H{
		"status": RespStatusSuccess,
	})
}

func GetTask(ctx *gin.Context) {
	t := time.NewTimer(10 * time.Second)

	taskKey := ctx.Param("task_key")
	if len(taskKey) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":       RespStatusErrParam,
			"errormessage": "parameter error",
		})
		return
	}

	ch := service.GetInflightChannel(taskKey)
	if ch == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":       RespStatusErrParam,
			"errormessage": "task doesn't exist",
		})
		return
	}

	select {
	case <-t.C:
		t.Stop()
		ctx.JSON(http.StatusOK, gin.H{
			"status": RespStatusListenTimeout,
		})
	case e := <-ch:
		ctx.JSON(http.StatusOK, gin.H{
			"status":       RespStatusSuccess,
			"task_key":     e.TaskKey,
			"parameter":    e.Parameter,
			"trigger_time": e.TriggerTime.Unix(),
		})
	}
}
