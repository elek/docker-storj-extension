// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"context"
	"github.com/elek/docker-storj-extension/backend"
	"github.com/spf13/pflag"
	"github.com/zeebo/errs"
	"go.uber.org/zap"
	"log"
	"storj.io/private/cfgstruct"
)

// Flags contains app configuration.
var Flags struct {
	Config backend.Config
}

func init() {
	cfgstruct.Bind(pflag.CommandLine, &Flags)
}

func main() {
	pflag.Parse()

	if err := run(context.Background(), Flags.Config); err != nil {
		log.Fatalf("%++v", err)
	}
}

func run(ctx context.Context, config backend.Config) error {
	logger := zap.NewExample()
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Println(err)
		}
	}()

	app, err := backend.NewApp(logger.Named("storj-extension"), config)
	if err != nil {
		return err
	}

	runErr := app.Run(ctx)
	closeErr := app.Close()
	return errs.Combine(runErr, closeErr)
}
