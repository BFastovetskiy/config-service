package main

import (
	"config-service/about"
	"config-service/appconfig"
	"config-service/configurationservice"
	"config-service/core"
	"config-service/database"
	"config-service/discovery"
	"config-service/httpserver"
	"config-service/utils"
	"context"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func Init(logger core.ILogger, cfg *appconfig.AppConfig) *Application {
	a := &Application{
		workDir: filepath.Dir(os.Args[0]),
		Config:  cfg,
		Logger:  logger,
	}
	a.Logger.Info("Start application")
	a.useSSL = utils.CheckExistSSL(a.workDir, cfg.PemFile, cfg.KeyFile)
	a.Logger.Infof("Database name [%s]", cfg.DatabaseName)

	if a.Database = database.InitDatabase(a.Logger, filepath.Join(a.workDir, about.Database_Directory), cfg.DatabaseName); a.Database == nil {
		a.Logger.Error("Application won't initialize")
		return nil
	}
	a.DiscoveryService = discovery.InitDiscoveryService(a, time.Duration(cfg.Discovery.Frequency)*time.Second)
	a.ConfigurationService = configurationservice.InitConfigurationService(a)
	return a
}

// GetDatabase implements IApplication.
func (a Application) GetDatabase() core.IDatabase {
	return a.Database
}

// GetLogger implements IApplication.
func (a Application) GetLogger() core.ILogger {
	return a.Logger
}

// GetWorkDir implements IApplication.
func (a Application) GetWorkDir() string {
	return a.workDir
}

// GetDiscoveryService implements IApplication.
func (a Application) GetDiscoveryService() core.IDiscoveryService {
	return a.DiscoveryService
}

// GetConfiguration implements IApplication.
func (a Application) GetConfiguration() *appconfig.AppConfig {
	return a.Config
}

func (a Application) GetConfigurationService() core.IConfigurationService {
	return a.ConfigurationService
}

// UseSSL implements IApplication.
func (a Application) UseSSL() bool {
	return a.useSSL
}

// Run implements IApplication.
func (a Application) Run() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTERM)
	isExit := false

	router := a.CreateRouter()
	if !a.Config.OnlySSL {
		a.http = httpserver.InitHttpServer(&a)
		a.Logger.Infof("Listener port [:%d]", a.Config.ServerPort)
		a.Logger.Info("Start Prometheus exporter on HTTP interface")
		go func() {
			srv := a.http.CreateServer(router, a.Config.ServerPort)
			a.http = a.http.SetServer(srv)
			a.http.Listen(false)
		}()
	}
	if a.useSSL {
		a.https = httpserver.InitHttpServer(&a)
		a.Logger.Infof("Listener use SSL [%v]", a.useSSL)
		a.Logger.Infof("Listener SSL port [:%d]", a.Config.ServerSSLPort)
		a.Logger.Info("Start Prometheus exporter on HTTPS interface")
		go func() {
			srv := a.https.CreateServer(router, a.Config.ServerSSLPort)
			a.https = a.https.SetServer(srv)
			a.https.Listen(true)
		}()
	}

	a.Logger.Infof("Listener cluster port [:%d]", cfg.ClusterPort)
	for {
		select {
		case <-sigs:
			a.Logger.Info("Received from OS the exit signal")
			isExit = true
		default:
			{
				time.Sleep(about.Application_Sleep_Timeout)
			}
		}
		if isExit {
			ctx, cancel := context.WithTimeout(context.Background(), about.Application_Cancel_Context_Timeout)
			defer cancel()
			if a.http != nil {
				if err := a.http.Shutdown(ctx); err != nil {
					a.Logger.Errorf("Problem stopping HTTP listener [%v]", err)
				}

			}
			if a.https != nil {
				if err := a.https.Shutdown(ctx); err != nil {
					a.Logger.Errorf("Problem stopping SSL listener [%v]", err)
				}

			}

			if err := a.Database.Close(); err != nil {
				a.Logger.Errorf("Problem closing DB [%v]", err)
			}

			break
		}
	}

}
