package controllers

import (
	"net/http"
	"os"
	"strings"

	about "config-service/about"

	"github.com/gin-gonic/gin"
)

func (c Controller) Actuator(ctx *gin.Context) {
	data := make(map[string]interface{}, 4)
	var baseUrl string = "/actuator"

	health := make(map[string]interface{}, 2)
	health["href"] = baseUrl + "/health"
	health["tepmlated"] = false
	data["health"] = health

	info := make(map[string]interface{}, 2)
	info["href"] = baseUrl + "/info"
	info["tepmlated"] = false
	data["info"] = info

	env := make(map[string]interface{}, 2)
	env["href"] = baseUrl + "/env"
	env["tepmlated"] = false
	data["env"] = env

	metrics := make(map[string]interface{}, 2)
	metrics["href"] = baseUrl + "/metrics"
	metrics["tepmlated"] = false
	data["metrics"] = metrics

	ctx.AsciiJSON(http.StatusOK, data)
}

func (c Controller) ActuatorHealth(ctx *gin.Context) {
	data := map[string]string{
		"status": "UP",
	}
	ctx.AsciiJSON(http.StatusOK, data)
}

func (c Controller) ActuatorMetrics(ctx *gin.Context) {
	ctx.AsciiJSON(http.StatusOK, gin.H{"status": "UP"})
}

func (c Controller) ActuatorInfo(ctx *gin.Context) {
	host, _ := os.Hostname()
	home, _ := os.UserHomeDir()
	data := map[string]interface{}{
		"Service": about.Application_Title,
		"Version": about.Application_Version,
		"Host":    host,
		"HomeDir": home,
	}
	ctx.AsciiJSON(http.StatusOK, data)
}

func (c Controller) ActuatorEnv(ctx *gin.Context) {
	env := os.Environ()
	data := map[string]string{}
	for _, val := range env {
		kv := strings.Split(val, "=")
		data[kv[0]] = kv[1]
	}
	ctx.AsciiJSON(http.StatusOK, data)
}
