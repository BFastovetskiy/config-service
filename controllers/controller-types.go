package controllers

import "config-service/core"

type Controller struct {
	app core.IApplication
	log core.ILogger
}

type Propertie struct {
	Application string `uri:"application" binding:"required"`
	Profile     string `uri:"profile"`
}

type RegisterService struct {
	Name        string `json:"name" binding:"required"`
	Application string `json:"application" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Port        int    `json:"port" binding:"required"`
	Health      string `json:"context" binding:"required"`
}
