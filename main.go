package main

import (
	"github.com/meidomx/airi/controller"
	"github.com/meidomx/airi/model"
	"github.com/meidomx/airi/service"

	scaffold "github.com/moetang/webapp-scaffold"

	"github.com/gin-gonic/gin"
)

func main() {
	webscaf, err := scaffold.NewFromConfigFile("airi.toml")
	if err != nil {
		panic(err)
	}

	webscaf.GetGin().Use(gin.Logger())
	webscaf.GetGin().Use(gin.Recovery())

	if err := webscaf.PreInitDb(); err != nil {
		panic(err)
	}
	model.Init(webscaf)
	service.Init()
	service.StartBackgroundJob()
	controller.InitRestController(webscaf)

	err = webscaf.SyncStart()
	if err != nil {
		panic(err)
	}
}
