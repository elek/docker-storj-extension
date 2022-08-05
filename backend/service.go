// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package backend

import (
	"github.com/zeebo/errs"
	"go.uber.org/zap"
)

// ErrService - pin service error class.
var ErrService = errs.Class("docker service")

type ServiceConfig struct {
}

// Service for querying ERC20 token information from ethereum chain.
//
// architecture: Service
type Service struct {
	log      *zap.Logger
	endpoint string
}

// NewService creates new token service instance.
func NewService(log *zap.Logger) *Service {
	return &Service{
		log: log,
	}
}
