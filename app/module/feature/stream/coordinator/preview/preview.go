package preview

import (
	"context"
	"time"

	"github.com/go-cmd/cmd"

	"noname001/logging"
	"noname001/dilemma"
)

// temp, POC related
// extract media server as commonality components
// so that all mod has stream capabilities
// and move this stuffs to mod_device

type DeviceChannelPreviewParams struct {
	ParentContext   context.Context
	ParentLogger    *logging.WrappedLogger
	ParentLogPrefix string

	DeviceCode        string
	DeviceChannelID   string
	DeviceStreamType  string
	StreamSource      string
	StreamDestination string
}

type DeviceChannelPreview struct {
	context   context.Context
	cancel    context.CancelFunc
	logger    *logging.WrappedLogger
	logPrefix string

	deviceCode        string
	deviceChannelID   string
	deviceStreamType  string
	streamSource      string
	streamDestination string

	lastRequested int64
}

func NewDeviceChannelPreview(params *DeviceChannelPreviewParams) (*DeviceChannelPreview) {
	dcp := &DeviceChannelPreview{}
	dcp.context, dcp.cancel = context.WithCancel(params.ParentContext)
	dcp.logger = params.ParentLogger
	dcp.logPrefix = params.ParentLogPrefix + ".dcp"

	dcp.deviceCode = params.DeviceCode
	dcp.deviceChannelID = params.DeviceChannelID
	dcp.deviceStreamType = params.DeviceStreamType
	dcp.streamSource = params.StreamSource
	dcp.streamDestination = params.StreamDestination

	dcp.UpdateLastRequested()

	return dcp
}

func (dcp *DeviceChannelPreview) UpdateLastRequested() {
	dcp.lastRequested = time.Now().Unix()
}

func (dcp *DeviceChannelPreview) LastRequested() (int64) {
	return dcp.lastRequested
}

func (dcp *DeviceChannelPreview) Start() {
	kommandOpts := cmd.Options{
		Buffered: false,
		CombinedOutput: false,
		Streaming: true,
	}
	kommandBin := dilemma.FFMPEG_BIN_PATH
	kommandArgs := []string{
		"-hide_banner",
		"-loglevel", "warning",
		"-stats_period", "10",

		"-rtsp_transport", "tcp",
		"-i", dcp.streamSource,

		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-tune", "zerolatency",
		"-profile:v", "baseline",
		"-force_key_frames:v", "expr:gte(t,n_forced*1)",
		"-b:v", "500000", "-minrate", "450000", "-maxrate", "550000",

		"-c:a", "aac",
		"-ac", "1",

		"-movflags", "faststart",
		"-movflags", "use_metadata_tags",
		"-map_metadata", "0",

		"-f", "rtsp",
		"-rtsp_transport", "tcp",
		"-rtsp_flags", "listen",
		"-ignore_io_errors", "true",

		dcp.streamDestination,
	}

	go func() {
		LimboLoop:
		for {
			kommand := cmd.NewCmdOptions(kommandOpts, kommandBin, kommandArgs...)

			go func() {
				for kommand.Stdout != nil || kommand.Stderr != nil {
					select {
					case <- dcp.context.Done():
						return

					case line, open := <- kommand.Stdout:
						if !open {
							kommand.Stdout = nil
							continue
						}
						dcp.logger.Debugf("%s: %s", dcp.logPrefix, line)

					case line, open := <- kommand.Stderr:
						if !open {
							kommand.Stderr = nil
							continue
						}
						dcp.logger.Debugf("%s: %s", dcp.logPrefix, line)
					}
				}
			}()

			kommandStatusChannel := kommand.Start()

			KommandLoop:
			for {
				select {
				case <- dcp.context.Done():
					kommand.Stop()
					break LimboLoop

				case kommandStatus := <- kommandStatusChannel:
					if kommandStatus.Complete {
						dcp.logger.Warnf("%s: preview-stream stopped unexpectedly. %s", dcp.logPrefix, kommandStatus)
					}
					dcp.logger.Warnf("%s: preview-stream restarting in 1 sec.", dcp.logPrefix)
					time.Sleep(1 *time.Second)
					break KommandLoop
				}
			}
		}
	}()

	dcp.UpdateLastRequested()
}
func (dcp *DeviceChannelPreview) Stop() {
	dcp.cancel()
}
