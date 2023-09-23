package core

import (
	appconfig "config-service/appconfig"
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	ILogger interface {
		GetSysLogger() *slog.Logger
		Info(msg string)
		Infow(msg string, args ...any)
		Infof(msg string, args ...any)
		Warning(msg string)
		Warningw(msg string, args ...any)
		Warningf(msg string, args ...any)
		Error(msg string)
		Errorw(msg string, args ...any)
		Errorf(msg string, args ...any)
		Debug(msg string)
		Debugw(msg string, args ...any)
		Debugf(msg string, args ...any)
	}

	IDatabase interface {
		GetValue(namespace, key string) (value []byte, err error)
		SetValue(namespace, key string, value []byte) error
		DelValue(namespace, key string) error
		HasValue(namespace, key string) (bool, error)
		CreateApplication(application string) error
		DeleteApplication(application string) error
		SetApplicationData(application string, item DiscoveryItem) error
		Backup(backup io.Writer) error
		Close() error
	}

	IHttpServer interface {
		Listen(useSsl bool)
		CreateServer(router *gin.Engine, port int) *http.Server
		SetServer(srv *http.Server) IHttpServer
		GetServer() *http.Server
		Shutdown(ctx context.Context) error
	}

	IDiscoveryService interface {
		GetService(service string) []DiscoveryItem
		SetService(service DiscoveryItem)
		DelService(service DiscoveryItem)
		GetSecret() string
		GetAllService() map[string][]DiscoveryItem
	}

	IConfigurationService interface {
		GetConfiguration(system, profile string) ConfigFile
	}

	IApplication interface {
		GetWorkDir() string
		GetLogger() ILogger
		GetDatabase() IDatabase
		GetDiscoveryService() IDiscoveryService
		GetConfiguration() *appconfig.AppConfig
		GetConfigurationService() IConfigurationService
		UseSSL() bool
		Run()
	}
)
