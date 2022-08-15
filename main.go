// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"context"
	"github.com/elek/docker-storj-extension/backend"
	"github.com/spf13/cobra"
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

	c := cobra.Command{}

	run := cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(context.Background(), Flags.Config)
		},
	}
	c.AddCommand(&run)

	dp := cobra.Command{
		Use: "dispatch",
		RunE: func(cmd *cobra.Command, args []string) error {
			n := backend.NewCliDispatcher()
			return n.Dispatch(args[0], args[1:])
		},
	}
	c.AddCommand(&dp)

	err := c.Execute()
	if err != nil {
		log.Fatalf("%++v", err)
	}

}

func runServer(ctx context.Context, config backend.Config) error {
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
