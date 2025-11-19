package dirty

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"noname001/logging"
	"time"

	"github.com/icholy/digest"
)

// this package is intended for when shit hit the fan
// and you need to code quick and dirty to test out what api works
// whatever goes

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
	var api *APIClient = &APIClient{}
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

	// TODO: https support
	notSecureBaseUrl := fmt.Sprintf("%s://%s:%s", api.protocol, api.hostname, api.port)
	baseUrl, err := url.Parse(notSecureBaseUrl)
	if err != nil {
		return nil, err
	}
	api.baseURL = baseUrl.String()
	return api, nil
}

func (api *APIClient) PingOK() (ret string, err error) {
	return "ping ok", nil
}

func (api *APIClient) PingNotOK() (ret string, err error) {
	return "", fmt.Errorf("ping not ok")
}