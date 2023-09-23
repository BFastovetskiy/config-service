package main

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

func (a Application) ServiceAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := a.tokenValid(ctx); err != nil {
			a.Logger.Error(err.Error())
			ctx.JSON(403, gin.H{"Error": fmt.Sprint(err)})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func (a Application) tokenValid(ctx *gin.Context) error {
	var tokenString = a.DiscoveryService.GetSecret()
	var secret = ctx.Request.Header.Get("Authorization")
	if secret != tokenString {
		return errors.New("Forbidden")
	}
	return nil
}
