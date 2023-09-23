package discovery

import (
	"config-service/core"
	"time"
)

type DiscoveryService struct {
	app       core.IApplication
	log       core.ILogger
	db        core.IDatabase
	duration  time.Duration
	secret    string
	discovery map[string][]core.DiscoveryItem
}
