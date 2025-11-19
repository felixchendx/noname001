package netcam

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"

	"noname001/logging"

	baseTyping "noname001/app/base/typing"

	deviceEv "noname001/app/module/common/device/event"
	"noname001/app/module/common/device/comm"

	camInterface "noname001/corebiz/integration/panasonic/netcam"
)

type PanasonicNetworkCameraParams struct {
	ParentContext context.Context
	Logger        *logging.WrappedLogger
	LogPrefix     string
	EvHub         *deviceEv.EventHub
	CommBundle    *comm.CommBundle
	Timezone      *time.Location

	ID    string
	Code  string
	Name  string
	State string
	Brand string

	Protocol, Hostname, Port, Username, Password string

	FallbackRTSPPort string
}

type PanasonicNetworkCamera struct {
	context    context.Context
	cancel     context.CancelFunc
	logger     *logging.WrappedLogger
	logPrefix  string
	evHub      *deviceEv.EventHub
	commBundle *comm.CommBundle

	timezone *time.Location
	cron     *cron.Cron
	cronJobs map[string]cron.EntryID

	id    string
	code  string
	name  string
	state string
	brand string

	protocol, hostname, port, username, password string

	fallbackRTSPPort string

	api *camInterface.APIClient

	cache *t_cache
}

func NewPanasonicNetworkCamera(params *PanasonicNetworkCameraParams) (*PanasonicNetworkCamera, error) {
	var err error

	dev := &PanasonicNetworkCamera{}
	dev.context, dev.cancel = context.WithCancel(params.ParentContext)
	dev.logger, dev.logPrefix = params.Logger, params.LogPrefix + ".pan-cam"
	dev.evHub = params.EvHub
	dev.commBundle = params.CommBundle

	dev.timezone = params.Timezone
	dev.cron     = cron.New(
		cron.WithLocation(dev.timezone),
		cron.WithSeconds(),
	)
	dev.cronJobs = make(map[string]cron.EntryID)

	dev.id    = params.ID
	dev.code  = params.Code
	dev.name  = params.Name
	dev.state = params.State
	dev.brand = params.Brand

	dev.protocol = params.Protocol
	dev.hostname = params.Hostname
	dev.port     = params.Port
	dev.username = params.Username
	dev.password = params.Password

	dev.fallbackRTSPPort = params.FallbackRTSPPort

	dev.api, err = camInterface.NewAPIClient(&camInterface.APIClientParams{
		dev.context,
		dev.logger,

		dev.protocol, dev.hostname, dev.port, dev.username, dev.password,
	})
	if err != nil {
		return nil, err
	}

	err = dev.setupCrons()
	if err != nil {
		return nil, err
	}

	dev.newCache()
	dev.cache.lDat.State = baseTyping.DEVICE_LIVE_STATE__NEW // TODO

	return dev, nil
}

type PanasonicNetworkCameraPatchParams struct {
	Name  string
	State string
	Brand string

	Protocol, Hostname, Port, Username, Password string

	FallbackRTSPPort string
}

// iffy
func (dev *PanasonicNetworkCamera) PatchAndReload(patchParams *PanasonicNetworkCameraPatchParams) (error) {
	var err error

	dev.name  = patchParams.Name
	dev.state = patchParams.State
	dev.brand = patchParams.Brand

	dev.protocol = patchParams.Protocol
	dev.hostname = patchParams.Hostname
	dev.port     = patchParams.Port
	dev.username = patchParams.Username
	dev.password = patchParams.Password

	dev.fallbackRTSPPort = patchParams.FallbackRTSPPort

	dev.api, err = camInterface.NewAPIClient(&camInterface.APIClientParams{
		dev.context,
		dev.logger,

		dev.protocol, dev.hostname, dev.port, dev.username, dev.password,
	})
	if err != nil {
		return nil
	}

	dev.newCache()

	dev.Reload()

	return nil
}
