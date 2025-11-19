package live

import (
	"context"
	"slices"
	"time"

	cmdlib "github.com/go-cmd/cmd"

	streamTyping "noname001/app/base/typing/stream"
)

// TODO: replace go-cmd with stdlib os.exec

const (
	// abspath, tested on:
	// - ubuntu: 22.04
	// - RHEL  : 9.6
	_FFMPEG_BIN_ABSPATH string = "/usr/bin/ffmpeg"
)

type ffmpegStreamer struct {
	context context.Context
	cancel  context.CancelFunc

	state streamTyping.LiveStreamerState

	ffmpegOpts *ffmpegOptions

	cmdOpts    cmdlib.Options
	cmdBinPath string
	cmdArgs    []string

	cmdInstance *cmdlib.Cmd

	lastStderr string
	lastStdout string
}

func (liveStream *LiveStream) newFfmpegStreamer() {
	if liveStream.streamerInstance != nil {
		liveStream.stopFfmpegStreamer()
	}

	streamer := &ffmpegStreamer{}
	streamer.context, streamer.cancel = context.WithCancel(liveStream.context)

	streamer.state = streamTyping.LIVE_STREAMER_STATE__NEW
	liveStream.lDat.StreamerState = streamer.state

	streamer.ffmpegOpts = liveStream.generateFfmpegOptions(liveStream.inputs, liveStream.outputs)

	streamer.cmdOpts    = cmdlib.Options{
		Buffered      : false,
		CombinedOutput: false,
		Streaming     : true,
	}
	streamer.cmdBinPath = _FFMPEG_BIN_ABSPATH // TODO: mod_stream config / app wide config for operator-defined ffmpeg path
	streamer.cmdArgs    = slices.Concat(
		streamer.ffmpegOpts.globalOpts.opts,
		streamer.ffmpegOpts.infileOpts.opts,
		streamer.ffmpegOpts.outfileOpts.video,
		streamer.ffmpegOpts.outfileOpts.videoFilter,
		streamer.ffmpegOpts.outfileOpts.audio,
		streamer.ffmpegOpts.outfileOpts.destination,
	)

	streamer.cmdInstance = cmdlib.NewCmdOptions(
		streamer.cmdOpts,
		streamer.cmdBinPath,
		streamer.cmdArgs...,
	)

	liveStream.streamerInstance = streamer
}

func (liveStream *LiveStream) startFfmpegStreamer() {
	streamer := liveStream.streamerInstance

	// todo: extract as optional logger
	go func() {
		for streamer.cmdInstance.Stdout != nil || streamer.cmdInstance.Stderr != nil {
			select {
			case <- streamer.context.Done():
				return

			case line, open := <- streamer.cmdInstance.Stdout:
				if !open {
					streamer.cmdInstance.Stdout = nil
					continue
				}
				streamer.lastStdout = line
				liveStream.logger.Debugf("%s: streamer(%s) stdout: %s", liveStream.logPrefix, liveStream.code, line)

			case line, open := <- streamer.cmdInstance.Stderr:
				if !open {
					streamer.cmdInstance.Stderr = nil
					continue
				}
				streamer.lastStderr = line
				liveStream.logger.Debugf("%s: streamer(%s) stderr: %s", liveStream.logPrefix, liveStream.code, line)
			}
		}
	}()


	statTicker := time.NewTicker(1 * time.Second)
	defer statTicker.Stop()

	streamer.state = streamTyping.LIVE_STREAMER_STATE__START
	liveStream.lDat.StreamerState = streamer.state
	cmdStatusChannel := streamer.cmdInstance.Start()

	kommandLoop:
	for {
		selectCase:
		select {
		case <- streamer.context.Done():
			liveStream.logger.Debugf("%s: streamer(%s) stopped-normal", liveStream.logPrefix, liveStream.code)
			streamer.state = streamTyping.LIVE_STREAMER_STATE__STOP_NORMAL
			liveStream.lDat.StreamerState = streamer.state
			streamer.cmdInstance.Stop()
			break kommandLoop

		// Exactly one Status is sent on the channel when the command ends.
		case cmdStatus := <- cmdStatusChannel:
			// go error, never encountered this yet...
			// if cmdStatus.Error != nil {
			// }

			// exit code differs depending on underlying program
			// in this case (ffmpeg)
			// switch cmdStatus.Exit {
			// case -1: // ffmpeg runnin
			// case 0: // ffmpeg exit normal
			// case 1: // ffmpeg exit with error
			// default: // (undocumented, not handled) exit code
			// }

			// true if not signalled / not stopped
			// meaning the command exit without interruption (hence called COMPLETE)
			// if cmdStatus.Complete {
			// 	// but current use case is supposed to NEVER EXITS (streaming forever)
			// 	// so, for this use case, COMPLETE indicating unexpected exit
			// 	return true
			// }
			
			liveStream.logger.Warnf(
				"%s: streamer(%s) stopped-unexpected, err: %s, exitCode: %s, complete: %s",
				liveStream.logPrefix,
				liveStream.code,
				cmdStatus.Error, cmdStatus.Exit, cmdStatus.Complete,
			)
			liveStream.logger.Warnf("%s: streamer(%s) stopped-unexpected, lastStdout: %s", liveStream.logPrefix, liveStream.code, streamer.lastStdout)
			liveStream.logger.Warnf("%s: streamer(%s) stopped-unexpected, lastStderr: %s", liveStream.logPrefix, liveStream.code, streamer.lastStderr)
			streamer.state = streamTyping.LIVE_STREAMER_STATE__STOP_UNEXPECTED
			liveStream.lDat.StreamerState = streamer.state

			// HRRMMMM....
			liveStream.execChan <- exec_mark_streamer_fail

			break kommandLoop

		// periodic check to see if still runnin
		case <- statTicker.C:
			cmdStatus := streamer.cmdInstance.Status()

			if cmdStatus.StartTs == 0 {
				// not started
				break selectCase
			} 

			if cmdStatus.StopTs == 0 {
				// is running
				streamer.state = streamTyping.LIVE_STREAMER_STATE__RUNNING
				// TODO: send event on state change 'others' > 'running' 
				liveStream.lDat.StreamerState = streamer.state
				break selectCase
			}

			// additional check is unnecessary, cmdStatusChannel will receive exactly one status, check there
		}
	}
}
// https://pkg.go.dev/github.com/go-cmd/cmd#Status
// type Status struct {
// 	Cmd      string
// 	PID      int
// 	Complete bool     // false if stopped or signaled
// 	Exit     int      // exit code of process
// 	Error    error    // Go error
// 	StartTs  int64    // Unix ts (nanoseconds), zero if Cmd not started
// 	StopTs   int64    // Unix ts (nanoseconds), zero if Cmd not started or running
// 	Runtime  float64  // seconds, zero if Cmd not started
// 	Stdout   []string // buffered STDOUT; see Cmd.Status for more info
// 	Stderr   []string // buffered STDERR; see Cmd.Status for more info
// }

func (liveStream *LiveStream) stopFfmpegStreamer() {
	if liveStream.streamerInstance == nil { return }

	liveStream.streamerInstance.cancel()
}
