package runner

import (
	"noname001/app/constant"
)

func (runna *Runner) exec() (err error) {
	if runna.runnerConfig == nil { return }
	if !runna.runnerConfig.Execute { return }

	runna.execSegmentDevice(runna.runnerConfig.SegmentDevice)

	return
}

func (runna *Runner) execSegmentDevice(seg SegmentDevice) {
	if seg.Mode == "" || seg.Mode == constant.RUNNER__MODE_NONE { return }

	for _, device := range seg.Devices {
		runna.walkDevice(seg.Mode, device)
	}
}
func (runna *Runner) walkDevice(mode string, item Device) {
	if item.Code == "" { return }

	de, messages := runna.service.FindDeviceByCode(item.Code)

	var (
		hasExistingData bool = (de != nil)

		doAdd    bool = false
		doEdit   bool = false
		doDelete bool = false
	)

	switch {
		case !hasExistingData && (mode == constant.RUNNER__MODE_ADD)        : doAdd = true
		case !hasExistingData && (mode == constant.RUNNER__MODE_ADD_OR_EDIT): doAdd = true
		case hasExistingData && (mode == constant.RUNNER__MODE_EDIT)        : doEdit = true
		case hasExistingData && (mode == constant.RUNNER__MODE_ADD_OR_EDIT) : doEdit = true
		case hasExistingData && (mode == constant.RUNNER__MODE_DELETE)      : doDelete = true
	}

	switch {
	case doAdd:
		de = runna.service.EmptyDevice()
		de.Code  = item.Code
		de.Name  = item.Name
		de.State = item.State
		de.Note  = item.Note

		de.Protocol = item.Protocol
		de.Hostname = item.Hostname
		de.Port     = item.Port
		de.Username = item.Username
		de.Password = item.Password
		de.Brand    = item.Brand

		de.FallbackRTSPPort = item.FallbackRTSPPort

		de, messages = runna.service.AddDevice(de)
		if messages.HasError() {
			runna.logger.Errorf("%s: %s", runna.logPrefix, messages)
			return
		}

		runna.logger.Infof("%s: %s", runna.logPrefix, messages)

	case doEdit:
		de.Name  = item.Name
		de.State = item.State
		de.Note  = item.Note

		de.Protocol = item.Protocol
		de.Hostname = item.Hostname
		de.Port     = item.Port
		de.Username = item.Username
		de.Password = item.Password
		de.Brand    = item.Brand

		de.FallbackRTSPPort = item.FallbackRTSPPort

		de, messages = runna.service.EditDevice(de.ID, de)
		if messages.HasError() {
			runna.logger.Errorf("%s: %s", runna.logPrefix, messages)
			return
		}

		runna.logger.Infof("%s: %s", runna.logPrefix, messages)

	case doDelete:
		messages = runna.service.DeleteDevice(de.ID)
		if messages.HasError() {
			runna.logger.Errorf("%s: %s", runna.logPrefix, messages)
			return
		}

		runna.logger.Infof("%s: %s", runna.logPrefix, messages)
	}
}
