package service

// TODO: prohibit direct access, implement service locator
var (
	serviceInstance *Service
)

func Instance() (*Service) {
	return serviceInstance
}
