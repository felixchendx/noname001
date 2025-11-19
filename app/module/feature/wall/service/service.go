package service

import (
	"context"

	"noname001/logging"
	"noname001/app/base/messaging"
	"noname001/app/module/feature/wall/event"
	"noname001/app/module/feature/wall/store"
	"noname001/app/module/feature/wall/comm"
	"noname001/app/module/feature/wall/coordinator"
)

type ServiceParams struct {
	Context     context.Context
	Logger      *logging.WrappedLogger
	LogPrefix   string
	EvHub       *event.EventHub
	Store       *store.Store
	CommBundle  *comm.CommBundle
	Coordinator *coordinator.Coordinator
}
type Service struct {
	context      context.Context
	cancel       context.CancelFunc

	logger       *logging.WrappedLogger
	logPrefix    string

	evHub        *event.EventHub
	store        *store.Store
	commBundle   *comm.CommBundle
	coordinator  *coordinator.Coordinator

	msgTemplates *messaging.MessageTemplateBundle
}

func NewService(params *ServiceParams) (*Service, error) {
	var err error

	svc := &Service{}
	svc.context, svc.cancel = context.WithCancel(params.Context)
	svc.logger = params.Logger
	svc.logPrefix = params.LogPrefix + ".svc"

	svc.evHub = params.EvHub
	svc.store = params.Store
	svc.commBundle = params.CommBundle
	svc.coordinator = params.Coordinator

	svc.msgTemplates = messaging.NewMessageTemplateBundle()

	err = svc.Init()

	if err != nil {
		return nil, err
	}

	return svc, nil
}

func (svc *Service) Init() (err error) {
	return
}

func (svc *Service) Start() (err error) {
	serviceInstance = svc
	return
}

func (svc *Service) PostStart() {
}

func (svc *Service) Stop() (err error) {
	serviceInstance = nil
	svc.cancel()
	return
}
