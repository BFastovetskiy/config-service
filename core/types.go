package core

import "time"

type DiscoveryItem struct {
	Name           string    `json:"name" binding:"required"`
	Application    string    `json:"application" binding:"required"`
	Address        string    `json:"address" binding:"required"`
	Port           int       `json:"port" binding:"required"`
	LastAccessTime time.Time `json:"last_access_time" binding:"required"`
	Health         string    `json:"health"`
	Configuration  string    `json:"-"`
}

type ConfigFile struct {
	Modify time.Time
	Body   []byte
}
