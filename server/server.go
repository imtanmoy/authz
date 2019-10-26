package server

import (
	"context"
	"github.com/imtanmoy/authz/config"
	"github.com/imtanmoy/authz/db"
	"github.com/imtanmoy/authz/logger"
	"net/http"
	"strconv"
	"time"
)

// Server provides an http.Server.
type Server struct {
	*http.Server
}

// NewServer creates and configures an APIServer serving all application routes.
func NewServer() (*Server, error) {
	logger.Info("configuring server...")
	api, err := New()
	if err != nil {
		return nil, err
	}

	host := config.Conf.SERVER.HOST
	port := strconv.Itoa(config.Conf.SERVER.PORT)
	addr := host + ":" + port

	srv := http.Server{
		Addr:    addr,
		Handler: api,
	}

	return &Server{&srv}, nil
}

// Start runs ListenAndServe on the http.Server with graceful shutdown.
func (srv *Server) Start(ctx context.Context) (err error) {
	logger.Info("starting server...")
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatalf("listen:%+s\n", err)
		}
	}()
	logger.Infof("Listening on %s\n", srv.Addr)

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		logger.Fatalf("server Shutdown Failed:%+s", err)
	}
	logger.Info("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}
	dbErr := db.Shutdown()
	if dbErr != nil {
		logger.Errorf("%s : %s", "Database shutdown failed", dbErr)
	}
	return
}
