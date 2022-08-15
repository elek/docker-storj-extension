package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zeebo/errs"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
)

type Cli struct {
	endpoint *Endpoint
	service  *Service
}

func NewCliDispatcher() *Cli {
	logger := zap.NewExample()

	service := NewService(logger)
	endpoint := NewEndpoint(logger, service)
	return &Cli{
		service:  service,
		endpoint: endpoint,
	}
}
func (c *Cli) Dispatch(cmd string, args []string) error {
	switch cmd {
	case "status":
		return c.handleMessage("GET", c.endpoint.Status, nil)
	case "local-images":
		return c.handleMessage("GET", c.endpoint.LocalImages, nil)
	case "pull":
		return c.handleMessage("POST", c.endpoint.Pull, nil)
	case "start":
		return c.handleMessage("POST", c.endpoint.Start, nil)
	case "stop":
		return c.handleMessage("POST", c.endpoint.Stop, nil)
	case "config":
		return c.handleMessage("POST", c.endpoint.Configure, struct {
			Bucket string `json:"bucket"`
			Grant  string `json:"grant"`
		}{
			Bucket: args[0],
			Grant:  args[1],
		})
	}

	return errs.New("Unsupported cmd " + cmd)
}

func (c *Cli) handleMessage(method string, h func(w http.ResponseWriter, r *http.Request), input interface{}) (err error) {
	s := []byte{}
	if input != nil {
		s, err = json.Marshal(input)
		if err != nil {
			return errs.Wrap(err)
		}

	}
	r, err := http.NewRequest(method, "", bytes.NewReader(s))
	if err != nil {
		return errs.Wrap(err)
	}
	recorder := httptest.NewRecorder()
	h(recorder, r)
	if recorder.Code != 200 {
		return errs.New(recorder.Body.String())
	}
	fmt.Println(recorder.Body.String())
	return nil
}
