package service

var (
	serviceInstance *Service
)

func Instance() (*Service) {
	return serviceInstance
}
