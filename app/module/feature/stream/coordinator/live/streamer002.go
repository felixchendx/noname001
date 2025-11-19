package live

import (
	"fmt"
	"strconv"
	"time"
)

// TODO
// -fflags +genpts

type ffmpegOptions struct {
	globalOpts  *ffmpegGlobalOptions
	infileOpts  *ffmpegInfileOptions
	outfileOpts *ffmpegOutfileOptions
}
type ffmpegGlobalOptions struct {
	opts []string
}
type ffmpegInfileOptions struct {
	opts []string
}
type ffmpegOutfileOptions struct {
	video       []string
	videoFilter []string

	audio       []string

	destination []string
}

func (liveStream *LiveStream) generateFfmpegOptions(inputs []inputIntface, outputs []*outputDestinationInternalMediaServer) (*ffmpegOptions) {
	// for now, only support one in and one out
	in, out := inputs[0], outputs[0]

	// TODO: err - states
	ffmpegOpts := &ffmpegOptions{}
	liveStream._generateGlobalOptions(ffmpegOpts)
	liveStream._generateInfileOptions(ffmpegOpts, in)
	liveStream._generateOutfileOptions(ffmpegOpts, out)
	liveStream._generateDynamicOptions(ffmpegOpts, in, out)

	return ffmpegOpts
}

func (liveStream *LiveStream) _generateGlobalOptions(ffmpegOpts *ffmpegOptions) {
	globalOpts := &ffmpegGlobalOptions{}

	globalOpts.opts = []string{
		"-hide_banner",
		// "-report",
		"-loglevel", "warning",
		"-stats_period", "5",
	}

	ffmpegOpts.globalOpts = globalOpts
}

func (liveStream *LiveStream) _generateInfileOptions(
	ffmpegOpts *ffmpegOptions,
	in         inputIntface,
) {
	infileOpts := &ffmpegInfileOptions{}
	
	switch typedInput := in.(type) {
	case *inputSourceModDevice:
		infileOpts.opts = []string{
			"-rtsp_transport", "tcp",
			"-i", typedInput.StreamURL,
		}

	case *inputSourceExternal:
		// TODO: protocol / type
		infileOpts.opts = []string{
			"-rtsp_transport", "tcp",
			"-i", typedInput.URL,
		}

	case *inputSourceFile:
		infileOpts.opts = []string{
			"-re",
			"-stream_loop", "-1",
			"-i", typedInput.Filepath,
		}

	default:
		xerr := fmt.Errorf("livestream infile options: err unimplemented source type %T", in)
		liveStream.logger.Errorf("%s: %s", liveStream.logPrefix, xerr.Error())
		return
	}

	ffmpegOpts.infileOpts = infileOpts
}

func (liveStream *LiveStream) _generateOutfileOptions(
	ffmpegOpts *ffmpegOptions,
	out        *outputDestinationInternalMediaServer,
) {
	outfileOpts := &ffmpegOutfileOptions{}

	switch out.TargetOutput.VideoCodec {
	case "h264":
		outfileOpts.video = []string{
			"-c:v", "libx264",
			"-preset", "ultrafast",
			"-tune", "zerolatency",
			"-profile:v", "baseline",
		}

	case "h265":
		outfileOpts.video = []string{
			"-c:v", "libx265",
			"-preset", "ultrafast",
			"-tune", "zerolatency",
			"-profile:v", "main",
		}

	default:
		xerr := fmt.Errorf("livestream outfile options: err unimplemented video codec %s", out.TargetOutput.VideoCodec)
		liveStream.logger.Errorf("%s: %s", liveStream.logPrefix, xerr.Error())
		return
	}


	switch out.TargetOutput.AudioCodec {
	case "aac":
		outfileOpts.audio = []string{
			"-c:a", "aac",
			"-b:a", "128k", // TODO: bitrate
			"-ac", "1",
		}

	// case "opus": // TODO
	
	default:
		xerr := fmt.Errorf("livestream outfile options: err unimplemented audio codec %s", out.TargetOutput.AudioCodec)
		liveStream.logger.Errorf("%s: %s", liveStream.logPrefix, xerr.Error())
		return
	}


	// only outputing to internal media server for now
	outfileOpts.destination = []string{
		"-movflags", "faststart",
		"-movflags", "use_metadata_tags",
		"-map_metadata", "0",

		"-f", "rtsp",
		"-rtsp_transport", "tcp",
		"-rtsp_flags", "listen",
		"-ignore_io_errors", "true",

		out.PublishURL,
	}

	ffmpegOpts.outfileOpts = outfileOpts
}

func (liveStream *LiveStream) _generateDynamicOptions(
	ffmpegOpts *ffmpegOptions,
	in         inputIntface,
	out        *outputDestinationInternalMediaServer,
) {
	// notes:
	// ffmpeg video filter chars that need escaping: [%, :]

	// TODO: embed monospace fonts, and use that font here

	// TODO: device n node information
	// TODO: codec n bitrate information
	// TODO: probe if no stream info is provided
	// TODO: dynamic font size, relative to video resolution
	// TODO: yada yada yada
	// TODO: dynamic to source timezone ? or compressor timezone ?

	var drawText, fontSize, fontColor string = "", "0.02*h", "yellow"
	epochLocalTz := time.Now().Unix() + (7 * 60 * 60)

	var targetBitrate, targetMinrate, targetMaxrate int = 0, 0, 0
	var displayInputBitrate, displayOutputBitrate, displayCompressionRatio string = "", "", ""

	switch typedInput := in.(type) {
	case *inputSourceModDevice:
		if typedInput.DeviceStreamInfo.VideoBitrate > 0 {
			switch {
			case out.TargetOutput.VideoCompression > 0 && int(out.TargetOutput.VideoBitrate) > 0:
				var compressionRatio int = int(out.TargetOutput.VideoCompression)

				targetBitrate, targetMinrate, targetMaxrate = liveStream._calculateTargetBitrate(
					typedInput.DeviceStreamInfo.VideoBitrate,
					compressionRatio,
				)

				if targetBitrate < int(out.TargetOutput.VideoBitrate) {
					compressionRatio = liveStream._calculateCompressionRatio(
						typedInput.DeviceStreamInfo.VideoBitrate,
						int(out.TargetOutput.VideoBitrate),
					)

					targetBitrate, targetMinrate, targetMaxrate = liveStream._calculateTargetBitrate(int(out.TargetOutput.VideoBitrate), 0)
				}

				displayInputBitrate     = fmt.Sprintf("%v bit/s", typedInput.DeviceStreamInfo.VideoBitrate)
				displayOutputBitrate    = fmt.Sprintf("%v bit/s", targetBitrate)
				displayCompressionRatio = fmt.Sprintf("%v pct", compressionRatio)

			case out.TargetOutput.VideoCompression > 0 && int(out.TargetOutput.VideoBitrate) <= 0:
				var compressionRatio int = int(out.TargetOutput.VideoCompression)

				targetBitrate, targetMinrate, targetMaxrate = liveStream._calculateTargetBitrate(
					typedInput.DeviceStreamInfo.VideoBitrate,
					compressionRatio,
				)

				displayInputBitrate     = fmt.Sprintf("%v bit/s", typedInput.DeviceStreamInfo.VideoBitrate)
				displayOutputBitrate    = fmt.Sprintf("%v bit/s", targetBitrate)
				displayCompressionRatio = fmt.Sprintf("%v pct", compressionRatio)

			case out.TargetOutput.VideoCompression <= 0 && int(out.TargetOutput.VideoBitrate) > 0:
				var compressionRatio int = 0

				compressionRatio = liveStream._calculateCompressionRatio(
					typedInput.DeviceStreamInfo.VideoBitrate,
					int(out.TargetOutput.VideoBitrate),
				)

				targetBitrate, targetMinrate, targetMaxrate = liveStream._calculateTargetBitrate(int(out.TargetOutput.VideoBitrate), 0)

				displayInputBitrate     = fmt.Sprintf("%v bit/s", typedInput.DeviceStreamInfo.VideoBitrate)
				displayOutputBitrate    = fmt.Sprintf("%v bit/s", targetBitrate)
				displayCompressionRatio = fmt.Sprintf("%v pct", compressionRatio)

			case out.TargetOutput.VideoCompression <= 0 && int(out.TargetOutput.VideoBitrate) <= 0:
				targetBitrate, targetMinrate, targetMaxrate = 0, 0, 0

				displayInputBitrate     = fmt.Sprintf("%v bit/s", typedInput.DeviceStreamInfo.VideoBitrate)
				displayOutputBitrate    = fmt.Sprintf("%v bit/s", typedInput.DeviceStreamInfo.VideoBitrate)
				displayCompressionRatio = "no compression"
			}
		}

		if typedInput.DeviceStreamInfo.VideoBitrate == 0 {
			if int(out.TargetOutput.VideoBitrate) > 0 {
				targetBitrate, targetMinrate, targetMaxrate = liveStream._calculateTargetBitrate(int(out.TargetOutput.VideoBitrate), 0)
			} else {
				targetBitrate, targetMinrate, targetMaxrate = 0, 0, 0
			}

			displayInputBitrate     = "n/a"
			displayOutputBitrate    = "n/a"
			displayCompressionRatio = "n/a"

			if targetBitrate > 0 {
				displayOutputBitrate = fmt.Sprintf("%v bit/s", targetBitrate)
			}
		}

		var datetimeText []string = []string{
			fmt.Sprintf(`[timestamp] = %%{pts\:gmtime\:%d}`, epochLocalTz),
		}

		var listDeviceText []string = []string{
			"[Node Info]",
			fmt.Sprintf("node-id = %s", "TODO"),
			"",
			"[Device Info]",
			fmt.Sprintf("device-name   = %s", typedInput.DeviceSnapshot.Hardware.DeviceName),
			"",
			"[Channel Info]",
			fmt.Sprintf("channel-id     = %s", typedInput.DeviceStreamInfo.ChannelID),
			fmt.Sprintf("channel-name   = %s", typedInput.DeviceStreamInfo.ChannelName),
			"",
			"[Stream Info]",
			fmt.Sprintf("input-video-bitrate  = %s", displayInputBitrate),
			fmt.Sprintf("output-video-bitrate = %s", displayOutputBitrate),
			fmt.Sprintf("compression-ratio    = %s", displayCompressionRatio),
		}

		drawText = fmt.Sprintf("%s,%s",
			prepTextLeftTop(datetimeText, fontColor, fontSize),
			prepTextLeftBottom(listDeviceText, fontColor, fontSize),
			// prepTextLeftMiddle(listStreamText, fontColor, fontSize),
		)

	case *inputSourceExternal:
		targetBitrate, targetMinrate, targetMaxrate = liveStream._calculateTargetBitrate(int(out.TargetOutput.VideoBitrate), 0)

		displayInputBitrate     = "n/a"
		displayOutputBitrate    = "n/a"
		displayCompressionRatio = "n/a"

		if targetBitrate > 0 {
			displayOutputBitrate = fmt.Sprintf("%v bit/s", targetBitrate)
		}

		var listTextFrame []string = []string{
			"[General Info]",
			fmt.Sprintf(`timestamp = %%{pts\:gmtime\:%d}`, epochLocalTz),
			"",
			"[Stream Info]",
			fmt.Sprintf("source (external) = %s", "redacted"),
			fmt.Sprintf("bitrate video-in  = %s", displayInputBitrate),
			fmt.Sprintf("bitrate video-out = %s", displayOutputBitrate),
			fmt.Sprintf("compression ratio = %s", displayCompressionRatio),
		}
		drawText = prepTextLeftBottom(listTextFrame, fontColor, fontSize)

	case *inputSourceFile:
		targetBitrate, targetMinrate, targetMaxrate = liveStream._calculateTargetBitrate(int(out.TargetOutput.VideoBitrate), 0)

		displayInputBitrate     = "n/a"
		displayOutputBitrate    = "n/a"
		displayCompressionRatio = "n/a"

		if targetBitrate > 0 {
			displayOutputBitrate = fmt.Sprintf("%v bit/s", targetBitrate)
		}

		var listTextFrame []string = []string{
			"[General Info]",
			fmt.Sprintf(`timestamp = %%{pts\:gmtime\:%d}`, epochLocalTz),
			"",
			"[Stream Info]",
			fmt.Sprintf("source (file)     = %s", typedInput.Filepath),
			fmt.Sprintf("bitrate video-in  = %s", displayInputBitrate),
			fmt.Sprintf("bitrate video-out = %s", displayOutputBitrate),
			fmt.Sprintf("compression ratio = %s", displayCompressionRatio),
		}
		drawText = prepTextLeftBottom(listTextFrame, fontColor, fontSize)

	default:
		drawText = ""
	}

	_vidFilter := []string{}
	_vidFilter = append(_vidFilter, "-force_key_frames:v", "expr:gte(t,n_forced*1)")

	if targetBitrate > 0 {
		_vidFilter = append(
			_vidFilter,
			"-b:v", strconv.Itoa(targetBitrate),
			"-minrate", strconv.Itoa(targetMinrate),
			"-maxrate", strconv.Itoa(targetMaxrate),
		)
	}

	_vidFilter = append(_vidFilter, "-vf", drawText)

	ffmpegOpts.outfileOpts.videoFilter = _vidFilter

	// temp
	liveStream.lDat.EstimatedOutputVideoBitrate = targetBitrate
}

func (liveStream *LiveStream) _calculateTargetBitrate(bitrate, compressionRatio int) (int, int, int) {
	var targetBitrate int = bitrate

	if 0 < compressionRatio && compressionRatio < 100 {
		targetBitrate     = (bitrate * (100 - compressionRatio)) / 100
	}
	
	var fluctuationPercentage int = 10
	var fluctuationThreshold  int = targetBitrate / fluctuationPercentage
	var targetMinrate         int = targetBitrate - fluctuationThreshold
	var targetMaxrate         int = targetBitrate + fluctuationThreshold

	return targetBitrate, targetMinrate, targetMaxrate
}

func (liveStream *LiveStream) _calculateCompressionRatio(inputBitrate, outputBitrate int) (int) {
	if inputBitrate > 0 {
		return 100 - ((outputBitrate * 100) / inputBitrate)
	}

	return 0
}
