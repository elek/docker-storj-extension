package backend

import (
	"context"
	"github.com/zeebo/errs"
	"io"
	"io/ioutil"
	"net/http"
)

type httpRequestOpt func(*http.Request)

func httpCall(ctx context.Context, method string, url string, body io.Reader, opts ...httpRequestOpt) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	for _, o := range opts {
		o(req)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	if resp.StatusCode > 299 {
		return nil, errs.Combine(errs.New("%s url is failed (%s): %s", method, resp.Status, url), resp.Body.Close())
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	return responseBody, errs.Combine(err, resp.Body.Close())
}
