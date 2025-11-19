package mediamtxserver

import (
	"noname001/thirdparty/mediamtx/httpapi/v1"
)

func (srv *MediaMTXServer) AddPathConfiguration(pathName, source string, onDemand bool) (error) {
	pathConfig := &v1.SimplePathConfiguration{
		Name: pathName,
		Source: source,
		SourceOnDemand: onDemand,
	}
	err := srv.apiClient.AddPathConfiguration(pathName, pathConfig)
	return err
}

func (srv *MediaMTXServer) ReplacePathConfiguration(pathName, source string, onDemand bool) (error) {
	pathConfig := &v1.SimplePathConfiguration{
		Name: pathName,
		Source: source,
		SourceOnDemand: onDemand,
	}
	err := srv.apiClient.ReplacePathConfiguration(pathName, pathConfig)
	return err
}

func (srv *MediaMTXServer) DeletePathConfiguration(pathName string) (error) {
	err := srv.apiClient.DeletePathConfiguration(pathName)
	return err
}
