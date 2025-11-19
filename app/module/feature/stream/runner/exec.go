package runner

import (
	"noname001/app/constant"
)

func (runna *Runner) exec() (err error) {
	if runna.runnerConfig == nil { return }
	if !runna.runnerConfig.Execute { return }

	runna.execSegmentStreamProfile(runna.runnerConfig.SegmentStreamProfile)
	runna.execSegmentStreamGroup(runna.runnerConfig.SegmentStreamGroup)

	return
}

func (runna *Runner) execSegmentStreamProfile(seg SegmentStreamProfile) {
	if seg.Mode == "" || seg.Mode == constant.RUNNER__MODE_NONE { return }

	for _, streamProfile := range seg.Profiles {
		runna.walkStreamProfile(seg.Mode, streamProfile)
	}
}
func (runna *Runner) walkStreamProfile(mode string, item StreamProfile) {
	if item.Code == "" { return }

	de, messages := runna.service.FindStreamProfileByCode(item.Code)

	if de == nil {
		de = runna.service.EmptyStreamProfile()
		de.Code = item.Code
		de.Name = item.Name
		de.State = item.State
		de.Note = item.Note

		de.TargetVideoCodec = item.TargetVideoCodec
		de.TargetVideoCompression = item.TargetVideoCompression
		de.TargetVideoBitrate = item.TargetVideoBitrate

		de.TargetAudioCodec = item.TargetAudioCodec
		de.TargetAudioCompression = item.TargetAudioCompression
		de.TargetAudioBitrate = item.TargetAudioBitrate

		de, messages = runna.service.AddStreamProfile(de)
		if messages.HasError() {
			runna.logger.Errorf("%s: %s", runna.logPrefix, messages)
			return
		}

		runna.logger.Infof("%s: %s", runna.logPrefix, messages)

	} else {

		de.Name = item.Name
		de.State = item.State
		de.Note = item.Note

		de.TargetVideoCodec = item.TargetVideoCodec
		de.TargetVideoCompression = item.TargetVideoCompression
		de.TargetVideoBitrate = item.TargetVideoBitrate

		de.TargetAudioCodec = item.TargetAudioCodec
		de.TargetAudioCompression = item.TargetAudioCompression
		de.TargetAudioBitrate = item.TargetAudioBitrate

		de, messages = runna.service.EditStreamProfile(de.ID, de)
		if messages.HasError() {
			runna.logger.Errorf("%s: %s", runna.logPrefix, messages)
			return
		}

		runna.logger.Infof("%s: %s", runna.logPrefix, messages)
	}
}

func (runna *Runner) execSegmentStreamGroup(seg SegmentStreamGroup) {
	if seg.Mode == "" || seg.Mode == constant.RUNNER__MODE_NONE { return }

	for _, streamGroup := range seg.Groups {
		runna.walkStreamGroup(seg.Mode, streamGroup)
	}
}
func (runna *Runner) walkStreamGroup(mode string, item StreamGroup) {
	if item.Code == "" { return }
	if item.StreamProfileCode == "" { return }

	refDE, refMessages := runna.service.FindStreamProfileByCode(item.StreamProfileCode)
	if refDE == nil {
		runna.logger.Errorf("%s: %s", runna.logPrefix, refMessages)
		return
	}

	de, messages := runna.service.FindStreamGroupByCode(item.Code)

	if de == nil {
		de = runna.service.EmptyStreamGroup()
		de.Code = item.Code
		de.Name = item.Name
		de.State = item.State
		de.Note = item.Note

		de.StreamProfileID = refDE.ID

		de, messages = runna.service.AddStreamGroup(de)
		if messages.HasError() {
			runna.logger.Errorf("%s: %s", runna.logPrefix, messages)
			return
		}

		runna.logger.Infof("%s: %s", runna.logPrefix, messages)

	} else {

		de.Name = item.Name
		de.State = item.State
		de.Note = item.Note

		de.StreamProfileID = refDE.ID

		de, messages = runna.service.EditStreamGroup(de.ID, de)
		if messages.HasError() {
			runna.logger.Errorf("%s: %s", runna.logPrefix, messages)
			return
		}

		runna.logger.Infof("%s: %s", runna.logPrefix, messages)
	}

	for _, streamItem := range item.Items {
		runna.walkStreamItem(mode, de.ID, streamItem)
	}
}

func (runna *Runner) walkStreamItem(mode string, groupID string, item StreamItem) {
	if item.Code == "" { return }

	de, messages := runna.service.FindStreamItemByCode(item.Code)

	if de == nil {
		de = runna.service.EmptyStreamItem()
		de.StreamGroupID = groupID
		de.Code = item.Code
		de.Name = item.Name
		de.State = item.State
		de.Note = item.Note

		de.SourceType = item.SourceType
		de.DeviceCode = item.DeviceCode
		de.DeviceChannelID = item.DeviceChannelID
		de.ExternalURL = item.ExternalURL
		de.Filepath = item.Filepath
		de.EmbeddedFilepath = item.EmbeddedFilepath

		de, messages = runna.service.AddStreamItem(de)
		if messages.HasError() {
			runna.logger.Errorf("%s: %s", runna.logPrefix, messages)
			return
		}

		runna.logger.Infof("%s: %s", runna.logPrefix, messages)

	} else {

		de.StreamGroupID = groupID
		de.Name = item.Name
		de.State = item.State
		de.Note = item.Note

		de.SourceType = item.SourceType
		de.DeviceCode = item.DeviceCode
		de.DeviceChannelID = item.DeviceChannelID
		de.ExternalURL = item.ExternalURL
		de.Filepath = item.Filepath
		de.EmbeddedFilepath = item.EmbeddedFilepath

		de, messages = runna.service.EditStreamItem(de.ID, de)
		if messages.HasError() {
			runna.logger.Errorf("%s: %s", runna.logPrefix, messages)
			return
		}

		runna.logger.Infof("%s: %s", runna.logPrefix, messages)
	}
}
