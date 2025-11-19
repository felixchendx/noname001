package mediamtxserver

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-cmd/cmd"

	"noname001/logging"

	"noname001/app/module/common/mediasrv/config"

	"noname001/thirdparty/mediamtx/httpapi/v1"
)

type MediaMTXServerParams struct {
	Context   context.Context
	Logger    *logging.WrappedLogger
	LogPrefix string
	Config    *config.RawModuleConfig
}

type MediaMTXServer struct {
	context   context.Context
	cancel    context.CancelFunc
	logger    *logging.WrappedLogger
	logPrefix string

	hostname string

	apiClient *v1.APIClient
	
	streamingPorts map[string]string
}

func NewMediaMTXServer(params *MediaMTXServerParams) (*MediaMTXServer, error) {
	srv := &MediaMTXServer{}
	srv.context, srv.cancel = context.WithCancel(params.Context)
	srv.logger = params.Logger
	srv.logPrefix = params.LogPrefix + ".mediaserver"

	srv.hostname = "localhost"

	srv.apiClient = v1.NewAPIClient(&v1.APIClientParams{
		ParentCtx: srv.context,
		Logger: srv.logger,
		Host: fmt.Sprintf("%s:%s", srv.hostname, params.Config.MediaServerAPIPort),
	})

	srv.streamingPorts = make(map[string]string)

	return srv, nil
}

func (srv *MediaMTXServer) Start() (err error) {
	kommandOpts := cmd.Options{
		Buffered: false,
		CombinedOutput: false,
		Streaming: true,
	}
	kommandBin := srv.unpackedBinPath()
	kommandArgs := []string{srv.unpackedConfigPath()}

	go func() {
		LimboLoop:
		for {
			kommand := cmd.NewCmdOptions(kommandOpts, kommandBin, kommandArgs...)

			go func() {
				for kommand.Stdout != nil || kommand.Stderr != nil {
					select {
					case <- srv.context.Done():
						return
	
					case line, open := <- kommand.Stdout:
						if !open {
							kommand.Stdout = nil
							continue
						}
						srv.logger.Debugf("%s: %s", srv.logPrefix, line)

					case line, open := <- kommand.Stderr:
						if !open {
							kommand.Stderr = nil
							continue
						}
						srv.logger.Debugf("%s: %s", srv.logPrefix, line)
					}
				}
			}()

			kommandStatusChannel := kommand.Start()

			KommandLoop:
			for {
				select {
				case <- srv.context.Done():
					kommand.Stop()
					break LimboLoop

				case kommandStatus := <- kommandStatusChannel:
					if kommandStatus.Complete {
						srv.logger.Warn("mediaserver stopped unexpectedly.", kommandStatus)
					}
					srv.logger.Warn("mediaserver restarting in 1 sec.")
					time.Sleep(1 * time.Second)
					break KommandLoop
				}
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)
	attempt, maxRetry, retryBackoff := 0, 9, (333 * time.Millisecond)
	for {
		attempt++

		mediamtxGlobalConf, _err := srv.apiClient.GetGlobalConfiguration()
		if _err != nil {
			// too soon? wait a bit for server to properly up
		} else {
			if mediamtxGlobalConf.RTSP {
				addrPart := strings.Split(mediamtxGlobalConf.RTSPAddress, ":")
				if len(addrPart) == 2 {
					srv.streamingPorts["rtsp"] = addrPart[1]
				}
			}
		
			if mediamtxGlobalConf.HLS {
				addrPart := strings.Split(mediamtxGlobalConf.HLSAddress, ":")
				if len(addrPart) == 2 {
					srv.streamingPorts["hls"] = addrPart[1]
				}
			}

			break
		}

		if attempt >= maxRetry {
			err := fmt.Errorf("Unable to communicate with mediaserver")

			if _err != nil {
				srv.logger.Errorf("%s: mediaserver api error: %s", srv.logPrefix, _err.Error())
			}
			srv.logger.Errorf("%s: %s", srv.logPrefix, err.Error())

			return err
		}

		time.Sleep(retryBackoff)
	}

	return
}

func (srv *MediaMTXServer) Stop() (err error) {
	return
}

func (srv *MediaMTXServer) Status() {}

// this publisher url does not need authn as long as it's via loopback ip
func (srv *MediaMTXServer) LocalRTSPPublisherBaseURL() (string) {
	// TODO: cache or make the formatting with path here (+checking)
	return fmt.Sprintf("rtsp://%s:%s", srv.hostname, srv.streamingPorts["rtsp"])
}

func (srv *MediaMTXServer) StreamingPorts() (map[string]string) {
	return srv.streamingPorts
}

func (srv *MediaMTXServer) ViewerAuthnPair() (string) {
	// TODO: coordiante with mediamtx
	return "viewer001:biwer001"
}

func (srv *MediaMTXServer) RelayAuthnPair() (string) {
	// TODO: coordinate with mediamtx
	return "relay001:really001"
}
