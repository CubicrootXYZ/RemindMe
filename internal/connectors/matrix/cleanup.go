package matrix

func (service *service) Cleanup() error {
	return service.matrixDatabase.Cleanup()
}
