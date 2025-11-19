package v1

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"noname001/logging"
)

// note: this is a locally testable api, need not paranoid handling

type APIClientParams struct {
	ParentCtx context.Context
	Logger    *logging.WrappedLogger

	Host string
}

type APIClient struct {
	context     context.Context
	cancel      context.CancelFunc
	logger      *logging.WrappedLogger

	httpTimeout time.Duration
	httpClient  http.Client
	baseURL     string
}

func NewAPIClient(params *APIClientParams) (*APIClient) {
	api := &APIClient{}
	api.context, api.cancel = context.WithCancel(params.ParentCtx)
	api.logger = params.Logger

	api.httpTimeout = 3 * time.Second
	api.httpClient = http.Client{}
	api.baseURL = fmt.Sprintf("http://%s", params.Host)

	return api
}
