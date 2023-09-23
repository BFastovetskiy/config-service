package controllers

import (
	"config-service/core"
	_ "net/http/pprof"
)

func StartController(app core.IApplication) *Controller {
	c := &Controller{
		app: app,
		log: app.GetLogger(),
	}
	return c
}
