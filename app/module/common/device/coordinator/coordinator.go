package coordinator

import (
	"context"
	"sync"
	"time"

	"noname001/logging"
	"noname001/app/base/sec"
	"noname001/app/module/common/device/event"
	"noname001/app/module/common/device/store"
	"noname001/app/module/common/device/comm"
	liveBase "noname001/app/module/common/device/coordinator/live/base"
)

type CoordinatorParams struct {
	Context   context.Context
	Logger    *logging.WrappedLogger
	LogPrefix string
	SecBundle *sec.DumbSecurityBundle
	EvHub     *event.EventHub
	Store     *store.Store
	CommBundle *comm.CommBundle
	Timezone  *time.Location
}
type Coordinator struct {
	context   context.Context
	cancel    context.CancelFunc

	logger    *logging.WrappedLogger
	logPrefix string
	secBundle *sec.DumbSecurityBundle

	evHub     *event.EventHub
	store     *store.Store
	commBundle *comm.CommBundle

	timezone  *time.Location

	// === liveDevicesMutex ===
	liveDevicesMutex    sync.Mutex
	liveDevices         map[string]liveBase.LiveDeviceIntface // k: deviceID
	liveDevicesCode     map[string]string // k: deviceCode, v: deviceID
	liveDevicesOrdered  []liveBase.LiveDeviceIntface
	// === liveDevicesMutex ===
}

func NewCoordinator(params *CoordinatorParams) (*Coordinator, error) {
	var err error

	coord := &Coordinator{}
	coord.context, coord.cancel = context.WithCancel(params.Context)
	coord.logger = params.Logger
	coord.logPrefix = params.LogPrefix + ".cdt"
	coord.secBundle = params.SecBundle

	coord.evHub = params.EvHub
	coord.store = params.Store
	coord.commBundle = params.CommBundle

	coord.timezone = params.Timezone

	coord.liveDevices        = make(map[string]liveBase.LiveDeviceIntface)
	coord.liveDevicesCode    = make(map[string]string)
	coord.liveDevicesOrdered = make([]liveBase.LiveDeviceIntface, 0)

	if err != nil {
		return nil, err
	}

	return coord, nil
}

func (coord *Coordinator) Start() (err error) {
	return
}

func (coord *Coordinator) PostStart() {
	coord.postStartRoutine()

	coord.eventListeners()
}

func (coord *Coordinator) Stop() (err error) {
	coord.cancel()
	return
}
