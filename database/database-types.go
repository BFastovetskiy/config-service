package database

import (
	"config-service/core"

	"github.com/boltdb/bolt"
)

type Database struct {
	Name string
	path string
	log  core.ILogger
	db   *bolt.DB
	root *bolt.Bucket
}
