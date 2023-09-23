package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sigs.k8s.io/yaml"
)

func (c Controller) GetServiceConfiguration(ctx *gin.Context) {
	var service Propertie
	if err := ctx.ShouldBindUri(&service); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}
	cf := c.app.GetConfigurationService().GetConfiguration(service.Application, service.Profile)
	if len(cf.Body) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "Not found configuration"})
		return
	}
	json, err := yaml.YAMLToJSON(cf.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}
	ctx.Data(http.StatusOK, "application/json", json)
}
