package controllers

import (
	"config-service/core"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (c Controller) GetAll(ctx *gin.Context) {
	data := c.app.GetDiscoveryService().GetAllService()
	ctx.AsciiJSON(http.StatusNotFound, data)
}

func (c Controller) LookuprService(ctx *gin.Context) {
	var service Propertie
	if err := ctx.ShouldBindUri(&service); err != nil {
		ctx.JSON(500, gin.H{"Error": err})
		return
	}
	c.log.Debugf("Input param: Configuration name = [%s]", service.Application)
	si := c.app.GetDiscoveryService().GetService(service.Application)
	if len(si) != 0 {
		ctx.AsciiJSON(http.StatusOK, si)
	} else {
		data := map[string]interface{}{
			"Service": service.Application,
			"Profile": service.Profile,
			"Message": "Not registred",
		}
		ctx.AsciiJSON(http.StatusNotFound, data)
	}

}

func (c Controller) RegisterService(ctx *gin.Context) {
	var service RegisterService

	if err := ctx.ShouldBindUri(&service); err != nil {
		ctx.JSON(500, gin.H{"Error": err})
		return
	}

	c.log.Debugf("Input param: Configuration name = [%s]", service.Application)

	si := core.DiscoveryItem{}
	si.Address = service.Address
	si.Port = service.Port
	si.Health = service.Health
	si.LastAccessTime = time.Now()
	c.app.GetDiscoveryService().SetService(si)

	ctx.JSON(http.StatusOK, gin.H{"Message": "Service " + service.Application + " registred"})
}
