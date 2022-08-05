// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package backend

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zeebo/errs"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"storj.io/common/errs2"
)

// Error is an error class for API http server error.
var Error = errs.Class("api server")

// APIConfig holds API endpoint configuration.
type APIConfig struct {
	Address string `default:"127.0.0.1:7070" help:"Address to listen on."`
}

// Server represents the backend API server.
//
// architecture: Endpoint
type Server struct {
	log      *zap.Logger
	listener net.Listener
	router   *mux.Router
	http     http.Server
}

// NewServer creates new API server instance.
func NewServer(log *zap.Logger, listener net.Listener) *Server {
	router := mux.NewRouter()
	router.Name("api").PathPrefix("/api/v0")
	apiRouter := router.GetRoute("api").Subrouter()

	return &Server{
		log:      log,
		listener: listener,
		router:   apiRouter,
		http: http.Server{
			Handler: router,
		},
	}
}

// NewAPI creates new API route and register endpoint methods.
func (server *Server) NewAPI(path string, register func(*mux.Router)) {
	router := server.router.PathPrefix(path).Subrouter()
	router.StrictSlash(true)
	register(router)
}

func (server *Server) Register(register func(*mux.Router)) {
	register(server.router)
}

// Run runs the server that host api endpoints.
func (server *Server) Run(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)

	var group errgroup.Group
	group.Go(func() error {
		<-ctx.Done()
		return Error.Wrap(server.http.Shutdown(context.Background()))
	})
	group.Go(func() error {
		defer cancel()
		err := Error.Wrap(server.http.Serve(server.listener))
		if errs2.IsCanceled(err) || errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		return err
	})

	return Error.Wrap(group.Wait())
}

// Close closes server and underlying listener.
func (server *Server) Close() error {
	return Error.Wrap(server.http.Close())
}
