package httpserver

import (
	"config-service/core"
	"net/http"
)

type HTTPServer struct {
	app core.IApplication
	log core.ILogger
	srv *http.Server
}
