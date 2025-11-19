package netcam

import (
	"context"
	"time"

	"noname001/logging"

	"noname001/corebiz/integration/base/apicall"
	"noname001/corebiz/integration/panasonic/netcam/v1"
)

const (
	DEFAULT_API_TIMEOUT = 10 * time.Second
	QUICK_API_TIMEOUT   = 3  * time.Second
)

type APIClientParams struct {
	Context context.Context
	Logger  *logging.WrappedLogger

	Protocol, Hostname, Port, Username, Password string
}
type APIClient struct {
	context context.Context
	cancel  context.CancelFunc
	logger  *logging.WrappedLogger

	protocol, hostname, port, username, password string

	apicallHandler *apicall.APICallHandler

	APIV1 *v1.APIClient
}

func NewAPIClient(params *APIClientParams) (*APIClient, error) {
	var err error

	api := &APIClient{}
	api.context, api.cancel = context.WithCancel(params.Context)
	api.logger = params.Logger

	api.protocol = params.Protocol
	api.hostname = params.Hostname
	api.port = params.Port
	api.username = params.Username
	api.password = params.Password

	api.apicallHandler = apicall.NewHandler(api.logger)

	api.APIV1, err = v1.NewAPIClient(&v1.APIClientParams{
		api.context,
		api.logger,
		api.protocol, api.hostname, api.port, api.username, api.password,
		DEFAULT_API_TIMEOUT,
	})
	if err != nil {
		return nil, err
	}

	return api, nil
}
