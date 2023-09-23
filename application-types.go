package main

import (
	"config-service/appconfig"
	"config-service/core"
)

type Application struct {
	Logger               core.ILogger
	Database             core.IDatabase
	workDir              string
	Config               *appconfig.AppConfig
	http                 core.IHttpServer
	https                core.IHttpServer
	DiscoveryService     core.IDiscoveryService
	ConfigurationService core.IConfigurationService
	useSSL               bool
}
