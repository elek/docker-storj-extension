// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package backend

import (
	"context"
	"github.com/spacemonkeygo/monkit/v3"
	"github.com/zeebo/errs"
	"net"
	"os"
	"path"
	"storj.io/private/debug"
	"storj.io/storj/private/lifecycle"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var mon = monkit.Package()

type Config struct {
	Debug debug.Config
	API   APIConfig
}

// App is the storjscan process that runs API endpoint.
//
// architecture: Peer
type App struct {
	Log *zap.Logger

	Debug struct {
		Listener net.Listener
		Server   *debug.Server
	}

	API struct {
		Listener net.Listener
		Server   *Server
	}

	Servers  *lifecycle.Group
	Services *lifecycle.Group
}

// NewApp creates new storjscan application instance.
func NewApp(log *zap.Logger, config Config) (*App, error) {
	app := &App{
		Log:      log,
		Services: lifecycle.NewGroup(log.Named("services")),
		Servers:  lifecycle.NewGroup(log.Named("servers")),
	}

	{ // API
		var err error

		service := NewService(log)
		e := NewEndpoint(log, service)

		address := config.API.Address
		if strings.HasPrefix(address, "socket:") {
			socketFile := strings.Split(address, ":")[1]
			_ = os.MkdirAll(path.Dir(socketFile), 0755)
			_ = os.Remove(socketFile)
			log.Info("Starting socket listener", zap.String("socket", socketFile))
			app.API.Listener, err = net.Listen("unix", socketFile)
		} else {
			log.Info("Starting TCP listener", zap.String("address", address))
			app.API.Listener, err = net.Listen("tcp", address)
		}
		if err != nil {
			return nil, err
		}

		app.API.Server = NewServer(log.Named("api:server"), app.API.Listener)
		app.API.Server.Register(e.Register)

		app.Servers.Add(lifecycle.Item{
			Name:  "api",
			Run:   app.API.Server.Run,
			Close: app.API.Server.Close,
		})
	}

	return app, nil
}

// Run runs storjscan until it's either closed or it errors.
func (app *App) Run(ctx context.Context) (err error) {
	defer mon.Task()(&ctx)(&err)

	group, ctx := errgroup.WithContext(ctx)

	app.Servers.Run(ctx, group)
	app.Services.Run(ctx, group)

	return group.Wait()
}

// Close closes all the resources.
func (app *App) Close() error {
	return errs.Combine(app.Services.Close(), app.Servers.Close())
}
