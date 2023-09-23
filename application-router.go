package main

import (
	controller "config-service/controllers"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func (a Application) CreateRouter() *gin.Engine {
	ctrl := controller.StartController(a)
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(sloggin.New(a.Logger.GetSysLogger().WithGroup("http")))

	actuator := router.Group("/actuator")
	actuator.GET("/", ctrl.Actuator)
	actuator.GET("/health", ctrl.ActuatorHealth)
	actuator.GET("/metrics", ctrl.ActuatorMetrics)
	actuator.GET("/info", ctrl.ActuatorInfo)
	actuator.GET("/env", ctrl.ActuatorEnv)

	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	system := router.Group("/system")
	//discovery.Use(a.ServiceAuth())
	system.POST("/backup", ctrl.Backup)

	config := router.Group("/api")
	//discovery.Use(a.ServiceAuth())
	config.GET("/service/:application", ctrl.GetServiceConfiguration)
	config.GET("/service/:application/:profile", ctrl.GetServiceConfiguration)

	discovery := router.Group("/api")
	//discovery.Use(a.ServiceAuth())
	discovery.GET("/discovery/", ctrl.GetAll)
	discovery.GET("/discovery/:application", ctrl.LookuprService)
	discovery.POST("/discovery/", ctrl.RegisterService)

	return router
}
