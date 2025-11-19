package sys

import (
	"context"

	"noname001/logging"
)

type SystemService struct {
	context       context.Context
	cancel        context.CancelFunc

	store         *Store
}

func NewSystemService(ctx context.Context, store *Store) (*SystemService, error) {
	svc := &SystemService{}
	svc.context, svc.cancel = context.WithCancel(ctx)
	svc.store = store
	return svc, nil
}

func (svc *SystemService) Stop() {
	svc.cancel()
}

func (svc *SystemService) PingDB() {
	pingPE, sev := svc.store.DB.Ping()
	if !sev.IsError() {
		logging.Logger.Infof("sys_svc: pingdb rep %s", pingPE)
	}
}
