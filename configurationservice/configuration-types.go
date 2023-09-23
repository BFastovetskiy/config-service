package configurationservice

import (
	"config-service/core"
)

type ConfigurationService struct {
	app     core.IApplication
	log     core.ILogger
	confDir string
	files   map[string]core.ConfigFile
}
