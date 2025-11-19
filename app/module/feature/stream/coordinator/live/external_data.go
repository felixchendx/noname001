package live

import (
	"fmt"
	"os"

	appConstant  "noname001/app/constant"
	baseTyping   "noname001/app/base/typing"
	streamTyping "noname001/app/base/typing/stream"

	mediasrvIntface "noname001/app/module/common/mediasrv/intface"

	deviceService "noname001/app/module/common/device/service"

	streamfs "noname001/app/module/feature/stream/filesystem"
)

// these stuffs are half-arsed, revisit later

func (liveStream *LiveStream) resetAllData() {
	liveStream.pDat = streamTyping.StreamPersistenceData{}

	liveStream.inputs = nil
	liveStream.outputs = nil
}

func (liveStream *LiveStream) initializeStreamData() (failCode string, err error) {
	streamItemPE, dbev1 := liveStream.store.DB.StreamItem__Get(liveStream.id)
	if dbev1.IsError() {
		return "db", fmt.Errorf("Store - DB has encountered internal error. Event ID: [%s].", dbev1.EventID())
	}
	if streamItemPE == nil {
		return "db", fmt.Errorf("Stream with id '%s' does not exist.", liveStream.id)
	}

	streamGroupPE, dbev2 := liveStream.store.DB.StreamGroup__Get(streamItemPE.StreamGroupID)
	if dbev2.IsError() {
		return "db", fmt.Errorf("Store - DB has encountered internal error. Event ID: [%s].", dbev2.EventID())
	}
	if streamGroupPE == nil {
		return "db", fmt.Errorf("Stream Group with id '%s' does not exist.", streamItemPE.StreamGroupID)
	}

	streamProfilePE, dbev3 := liveStream.store.DB.StreamProfile__Get(streamGroupPE.StreamProfileID)
	if dbev3.IsError() {
		return "db", fmt.Errorf("Store - DB has encountered internal error. Event ID: [%s].", dbev3.EventID())
	}
	if streamProfilePE == nil {
		return "dep:profile", fmt.Errorf("Stream Profile with id '%s' does not exist.", streamGroupPE.StreamProfileID)
	}

	liveStream.code = streamItemPE.Code

	liveStream.pDat.State            = streamItemPE.State
	liveStream.pDat.SourceType       = streamItemPE.SourceType
	liveStream.pDat.DeviceCode       = streamItemPE.DeviceCode
	liveStream.pDat.DeviceChannelID  = streamItemPE.DeviceChannelID
	liveStream.pDat.DeviceStreamType = streamItemPE.DeviceStreamType
	liveStream.pDat.ExternalURL      = streamItemPE.ExternalURL
	liveStream.pDat.Filepath         = streamItemPE.Filepath
	liveStream.pDat.EmbeddedFilepath = ""

	liveStream.pDat.GroupID    = streamGroupPE.ID
	liveStream.pDat.GroupState = streamGroupPE.State

	liveStream.pDat.ProfileID              = streamProfilePE.ID
	liveStream.pDat.ProfileCode            = streamProfilePE.Code
	liveStream.pDat.ProfileName            = streamProfilePE.Name
	liveStream.pDat.VideoEnabled           = true
	liveStream.pDat.TargetVideoCodec       = streamProfilePE.TargetVideoCodec
	liveStream.pDat.TargetVideoCompression = streamProfilePE.TargetVideoCompression
	liveStream.pDat.TargetVideoBitrate     = streamProfilePE.TargetVideoBitrate
	liveStream.pDat.AudioEnabled           = true
	liveStream.pDat.TargetAudioCodec       = streamProfilePE.TargetAudioCodec
	liveStream.pDat.TargetAudioCompression = streamProfilePE.TargetAudioCompression
	liveStream.pDat.TargetAudioBitrate     = streamProfilePE.TargetAudioBitrate

	return "", nil
}

func (liveStream *LiveStream) reloadStreamData() (failCode string, err error) {
	return liveStream.initializeStreamData()
}

func (liveStream *LiveStream) loadInputData() (failCode string, err error) {
	switch liveStream.pDat.SourceType {
	case string(appConstant.STREAM__SOURCE_TYPE__MOD_DEVICE):
		deviceSnapshot, messages := deviceService.Instance().GetDeviceSnapshotByCode(liveStream.pDat.DeviceCode)
		if messages.HasError() {
			return "dep:device", fmt.Errorf("%s", messages.FirstErrorMessageString())
		}

		// TODO: read device state, more fleshed out err states

		streamInfoFromDevice, messages := deviceService.Instance().GetStreamInfo(
			&deviceService.DeviceIdentifier{
				Code: liveStream.pDat.DeviceCode,
			},
			liveStream.pDat.DeviceChannelID,
			appConstant.BrandStreamType(liveStream.pDat.DeviceStreamType), // TEMP, unsafe casting
		)
		if messages.HasError() {
			return "dep:device", fmt.Errorf("%s", messages.FirstErrorMessageString())
		}

		deviceStreamInfo := &baseTyping.BaseDeviceStreamInfo{
			ChannelID            : streamInfoFromDevice.ChannelID,
			ChannelName          : streamInfoFromDevice.ChannelName,
			Enabled              : streamInfoFromDevice.Enabled,

			StreamURL            : streamInfoFromDevice.StreamURL,

			// VideoEnabled         : ,
			VideoResolutionWidth : streamInfoFromDevice.VideoResolutionWidth,
			VideoResolutionHeight: streamInfoFromDevice.VideoResolutionHeight,
			VideoCodec           : streamInfoFromDevice.VideoCodecType,
			// VideoFPS             : 
			VideoBitrate         : streamInfoFromDevice.VideoBitrate,

			// AudioEnabled         : 
			AudioCodec           : streamInfoFromDevice.AudioCodecType,
			// AudioBitrate         : 
		}

		liveStream.inputs = []inputIntface{
			&inputSourceModDevice{
				StreamURL       : deviceStreamInfo.StreamURL,

				DeviceSnapshot  : deviceSnapshot,
				DeviceStreamInfo: deviceStreamInfo,
			},
		}

	case string(appConstant.STREAM__SOURCE_TYPE__EXTERNAL):
		liveStream.inputs = []inputIntface{
			&inputSourceExternal{
				URL: liveStream.pDat.ExternalURL,
			},
		}

	case string(appConstant.STREAM__SOURCE_TYPE__FILE):
		absFilepath := fmt.Sprintf("%s/%s", streamfs.StreamLocalDir(), liveStream.pDat.Filepath)

		_file, err := os.Open(absFilepath)
		defer _file.Close()
		if err != nil {
			return "dep:file", err
		}

		liveStream.inputs = []inputIntface{
			&inputSourceFile{
				Filepath: absFilepath,
			},
		}

	default:
		return "other", fmt.Errorf("unimplemented source type '%s'", liveStream.pDat.SourceType)
	}

	return "", nil
}

func (liveStream *LiveStream) loadOutputData() {
	var publishURL   string = mediasrvIntface.Provider().PublishStreamViewURL(liveStream.code)
	var targetOutput        = &targetOutput{
		StreamProfileCode: liveStream.pDat.ProfileCode,
		StreamProfileName: liveStream.pDat.ProfileName,

		VideoFPS        : 0,
		VideoWidth      : 0,
		VideoHeight     : 0,
		VideoCodec      : liveStream.pDat.TargetVideoCodec,
		VideoCompression: liveStream.pDat.TargetVideoCompression,
		VideoBitrate    : liveStream.pDat.TargetVideoBitrate,

		AudioCodec      : liveStream.pDat.TargetAudioCodec,
		AudioCompression: liveStream.pDat.TargetAudioCompression,
		AudioBitrate    : liveStream.pDat.TargetAudioBitrate,

		// ShowTimestamp   : liveStream.pDat.ShowTimestamp,
		// ShowVideoInfo   : liveStream.pDat.ShowVideoInfo,
		// ShowAudioInfo   : liveStream.pDat.ShowAudioInfo,
		// ShowSiteInfo    : liveStream.pDat.ShowSiteInfo,
	}

	liveStream.outputs = []*outputDestinationInternalMediaServer{
		&outputDestinationInternalMediaServer{
			PublishURL  : publishURL,
			TargetOutput: targetOutput,
		},
	}
}
