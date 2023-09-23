package httpserver

import (
	"context"
	"net/http"
	"path/filepath"
	"strconv"

	"config-service/about"
	"config-service/core"

	"github.com/gin-gonic/gin"

	_ "net/http/pprof"
)

func InitHttpServer(app core.IApplication) *HTTPServer {
	h := &HTTPServer{
		app: app,
		log: app.GetLogger(),
	}
	return h
}

func (h HTTPServer) Listen(useSsl bool) {
	if !useSsl {
		h.log.Info("Starting HTTP service with out SSL")
		if err := h.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			h.log.Errorf("Error: %s", err)
		}
	} else {
		path := filepath.Join(h.app.GetWorkDir(), about.SSL_Directory)
		crt := filepath.Join(path, h.app.GetConfiguration().PemFile)
		key := filepath.Join(path, h.app.GetConfiguration().KeyFile)

		h.log.Info("Starting HTTP service with SSL")
		if err := h.srv.ListenAndServeTLS(crt, key); err != nil {
			h.log.Errorf("Error: %s", err)
		}
	}
}

func (h HTTPServer) CreateServer(router *gin.Engine, port int) *http.Server {
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: router,
	}
	return server
}

func (h HTTPServer) SetServer(srv *http.Server) core.IHttpServer {
	h.srv = srv
	return h
}

func (h HTTPServer) GetServer() *http.Server {
	return h.srv
}

func (h HTTPServer) Shutdown(ctx context.Context) error {
	err := h.srv.Shutdown(ctx)
	return err
}
