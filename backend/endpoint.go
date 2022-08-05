// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package backend

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zeebo/errs"
	"go.uber.org/zap"
)

var ErrEndpoint = errs.Class("backend endpoint")

type Endpoint struct {
	log     *zap.Logger
	service *Service
}

func NewEndpoint(log *zap.Logger, service *Service) *Endpoint {
	return &Endpoint{
		log:     log,
		service: service,
	}
}

func (endpoint *Endpoint) Register(router *mux.Router) {
	router.HandleFunc("/test", endpoint.Test).Methods(http.MethodGet)
}

func (endpoint *Endpoint) Test(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer mon.Task()(&ctx)(&err)

	err = json.NewEncoder(w).Encode(struct {
		Status string
	}{
		Status: "running",
	})
	if err != nil {
		endpoint.log.Error("failed to write json pins response", zap.Error(ErrEndpoint.Wrap(err)))
		return
	}
}

// serveJSONError writes JSON error to response output stream.
func (endpoint *Endpoint) serveJSONError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)

	var response struct {
		Error string `json:"error"`
	}

	response.Error = err.Error()

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		endpoint.log.Error("failed to write json error response", zap.Error(ErrEndpoint.Wrap(err)))
		return
	}
}
