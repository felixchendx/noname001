package coordinator

// tempe
// review after restructurig stream stuffs, and reviewing cache stuffs

func (coord *Coordinator) StreamingPorts() (map[string]string) {
	return coord.mediaServer.StreamingPorts()
}

func (coord *Coordinator) RelayAuthnPair() (string) {
	return coord.mediaServer.RelayAuthnPair()
}


func (coord *Coordinator) AddPathConfiguration(pathName, source string, onDemand bool) (error) {
	return coord.mediaServer.AddPathConfiguration(pathName, source, onDemand)
}

func (coord *Coordinator) ReplacePathConfiguration(pathName, source string, onDemand bool) (error) {
	return coord.mediaServer.ReplacePathConfiguration(pathName, source, onDemand)
}

func (coord *Coordinator) DeletePathConfiguration(pathName string) (error) {
	return coord.mediaServer.DeletePathConfiguration(pathName)
}
