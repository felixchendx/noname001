package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/icholy/digest"

	"noname001/logging"
)

// TODO: read moar https://pkg.go.dev/net/http

type APIClientParams struct {
	Context    context.Context
	Logger     *logging.WrappedLogger

	Protocol, Hostname, Port, Username, Password string

	APITimeout time.Duration
}
type APIClient struct {
	context     context.Context
	cancel      context.CancelFunc
	logger      *logging.WrappedLogger
	
	protocol, hostname, port, username, password string

	httpTimeout time.Duration 
	httpClient  http.Client
	baseURL     string
}

func NewAPIClient(params *APIClientParams) (*APIClient, error) {
	api := &APIClient{}
	api.context, api.cancel = context.WithCancel(params.Context)
	api.logger = params.Logger
	api.protocol = params.Protocol
	api.hostname = params.Hostname
	api.port = params.Port
	api.username = params.Username
	api.password = params.Password

	api.httpTimeout = params.APITimeout
	api.httpClient = http.Client{
		Transport: &digest.Transport{
			Username: api.username,
			Password: api.password,
		},
		Timeout: api.httpTimeout, // TODO: forego this, context timeout enough ?
	}

	// TODO: protocol validation
	baseURL, err := url.Parse(fmt.Sprintf("%s://%s:%s", api.protocol, api.hostname, api.port))
	if err != nil {
		return nil, err
	}
	api.baseURL = baseURL.String()

	return api, nil
}
