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
	router.HandleFunc("/status", endpoint.Status).Methods(http.MethodGet)
	router.HandleFunc("/start", endpoint.Start).Methods(http.MethodPost)
	router.HandleFunc("/stop", endpoint.Stop).Methods(http.MethodPost)
	router.HandleFunc("/configure", endpoint.Configure).Methods(http.MethodPost)
	router.HandleFunc("/images/local", endpoint.LocalImages).Methods(http.MethodGet)
	router.HandleFunc("/images/remote", endpoint.RemoteImages).Methods(http.MethodGet)
	router.HandleFunc("/push", endpoint.Push).Methods(http.MethodPost)
	router.HandleFunc("/pull", endpoint.Pull).Methods(http.MethodPost)
}

func (endpoint *Endpoint) Status(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer mon.Task()(&ctx)(&err)

	st, err := endpoint.service.Status()
	if err != nil {
		endpoint.serveJSONError(w, http.StatusInternalServerError, err)
		return
	}

	resp := struct {
		Status string `json:"status"`
	}{
		Status: st,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		endpoint.serveJSONError(w, http.StatusInternalServerError, err)
		return
	}

	if err != nil {
		endpoint.log.Error("failed to write json pins response", zap.Error(ErrEndpoint.Wrap(err)))
		return
	}
}

func (endpoint *Endpoint) Configure(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer mon.Task()(&ctx)(&err)

	req := struct {
		Bucket string `json:"bucket"`
		Grant  string `json:"grant"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		endpoint.serveJSONError(w, http.StatusBadRequest, err)
		return
	}

	err = endpoint.service.Configure(req.Bucket, req.Grant)
	if err != nil {
		endpoint.serveJSONError(w, http.StatusInternalServerError, err)
		return
	}
}

func (endpoint *Endpoint) Push(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer mon.Task()(&ctx)(&err)

	req := struct {
		Image string `json:"image"`
		Tag   string `json:"tag"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		endpoint.serveJSONError(w, http.StatusBadRequest, err)
		return
	}

	err = endpoint.service.Push(req.Image, req.Tag)
	if err != nil {
		endpoint.serveJSONError(w, http.StatusBadRequest, err)
		return
	}
}

func (endpoint *Endpoint) Pull(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer mon.Task()(&ctx)(&err)

	req := struct {
		Image string `json:"image"`
		Tag   string `json:"tag"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		endpoint.serveJSONError(w, http.StatusBadRequest, err)
		return
	}

	err = endpoint.service.Pull(req.Image, req.Tag)
	if err != nil {
		endpoint.serveJSONError(w, http.StatusBadRequest, err)
		return
	}
}

func (endpoint *Endpoint) Start(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer mon.Task()(&ctx)(&err)

	err = endpoint.service.Start()
	if err != nil {
		endpoint.serveJSONError(w, http.StatusBadRequest, err)
		return
	}
}

func (endpoint *Endpoint) Stop(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer mon.Task()(&ctx)(&err)

	err = endpoint.service.Stop()
	if err != nil {
		endpoint.serveJSONError(w, http.StatusBadRequest, err)
		return
	}
}

func (endpoint *Endpoint) LocalImages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer mon.Task()(&ctx)(&err)

	images, err := endpoint.service.LocalImages()
	if err != nil {
		endpoint.serveJSONError(w, http.StatusInternalServerError, err)
		return
	}
	err = json.NewEncoder(w).Encode(&images)
	if err != nil {
		endpoint.serveJSONError(w, http.StatusInternalServerError, err)
		return
	}
}

func (endpoint *Endpoint) RemoteImages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer mon.Task()(&ctx)(&err)

	images, err := endpoint.service.RemoteImages()
	if err != nil {
		endpoint.serveJSONError(w, http.StatusInternalServerError, err)
		return
	}
	err = json.NewEncoder(w).Encode(&images)
	if err != nil {
		endpoint.serveJSONError(w, http.StatusInternalServerError, err)
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
