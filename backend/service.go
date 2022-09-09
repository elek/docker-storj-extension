// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeebo/errs"
	"go.uber.org/zap"
	"os/exec"
	"strings"
)

var base = "127.0.0.1:5000"

// ErrService - pin service error class.
var ErrService = errs.Class("docker service")

type ServiceConfig struct {
}

// Service for querying ERC20 token information from ethereum chain.
//
// architecture: Service
type Service struct {
	log        *zap.Logger
	endpoint   string
	supervisor *Supervisor
}

// NewService creates new token service instance.
func NewService(log *zap.Logger) *Service {
	return &Service{
		log: log,
		supervisor: &Supervisor{
			log: log,
		},
	}
}

func (s Service) Create(bucket string, grant string) error {
	c := exec.Command("docker", "rm", "storj-registry")
	_, _ = c.CombinedOutput()

	c = exec.Command("docker",
		"create",
		"--name", "storj-registry",
		"-p", "9999:5000",
		"-e", "REGISTRY_STORAGE_STORJ_BUCKET="+bucket,
		"-e", "REGISTRY_STORAGE_STORJ_ACCESSGRANT="+grant,
		"ghcr.io/elek/distribution:618d19fb")
	out, err := c.CombinedOutput()
	if err != nil {
		return errs.New(err.Error() + " " + string(out))
	}
	s.log.Info("Container has been created", zap.String("output", string(out)))
	return err
}

func (s Service) Status() (string, error) {
	return s.supervisor.Status()
}

func (s Service) Start() error {
	return s.supervisor.Start()
}

func (s Service) Stop() error {
	return s.supervisor.Stop()
}

type Image struct {
	Name string
	Tag  string
	Id   string
	Size string
}

func (s Service) RemoteImages() ([]Image, error) {
	images := make([]Image, 0)
	ctx := context.TODO()
	url := "http://" + base + "/v2/_catalog"
	catalog, err := httpCall(ctx, "GET", url, nil)
	if err != nil {
		return images, errs.Wrap(err)
	}

	k := struct {
		Repositories []string `json:"repositories"`
	}{}
	err = json.Unmarshal(catalog, &k)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	for _, name := range k.Repositories {
		url = fmt.Sprintf("http://%s/5000/v2/%s/tags/list", base, name)
		tagsResponse, err := httpCall(ctx, "GET", url, nil)
		if err != nil {
			return nil, errs.Wrap(err)
		}

		tags := struct {
			Name string   `json:"name"`
			Tags []string `json:"tags"`
		}{}
		err = json.Unmarshal(tagsResponse, &tags)
		if err != nil {
			return nil, errs.Wrap(err)
		}

		for _, tag := range tags.Tags {
			images = append(images, Image{
				Name: name,
				Tag:  tag,
			})
		}
	}

	return images, nil

}
func (s Service) LocalImages() ([]Image, error) {
	c := exec.Command("docker",
		"images",
		"--format",
		"{{.Repository}} {{.Tag}} {{.ID}} {{.Size}}")
	out, err := c.CombinedOutput()
	if err != nil {
		return nil, errs.New(string(out))
	}
	var res []Image
	for ix, line := range strings.Split(string(out), "\n") {
		if ix == 0 {
			continue
		}
		parts := strings.Split(line, " ")
		if len(parts) < 4 {
			continue
		}
		if parts[1] == "<none>" {
			continue
		}
		res = append(res, Image{
			Name: parts[0],
			Tag:  parts[1],
			Id:   parts[2],
			Size: parts[3],
		})
	}

	return res, nil
}

func (s Service) Push(name string, tag string) error {
	name = strings.Trim(name, "\n")
	parts := strings.Split(name, "/")
	ref := name
	if len(parts) > 2 {
		ref = name + ":" + tag
	}

	c := exec.Command("docker",
		"tag",
		name+":"+tag,
		base+"/"+ref)
	out, err := c.CombinedOutput()
	s.log.Info("Container is tagged", zap.String("output", string(out)))
	if err != nil {
		return errs.New(string(out))
	}

	c = exec.Command("docker",
		"push",
		base+"/"+ref,
	)
	out, err = c.CombinedOutput()
	s.log.Info("Container is pushed", zap.String("output", string(out)))
	if err != nil {
		return errs.New(string(out))
	}

	return nil
}

func (s Service) Pull(name string, tag string) error {
	c := exec.Command("docker",
		"pull",
		base+"/"+name+":"+tag)
	out, err := c.CombinedOutput()
	s.log.Info("Container is pulled", zap.String("output", string(out)))
	if err != nil {
		return errs.New(string(out))
	}
	return nil
}

func (s Service) Configure(bucket string, grant string) error {
	return s.supervisor.Configure(bucket, grant)
}
