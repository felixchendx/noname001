package netcam

import (
	nodeConstant "noname001/node/constant"
)

func (dev *PanasonicNetworkCamera) setupCrons() (error) {
	entryID, err := dev.cron.AddFunc(
		nodeConstant.CROSSNODE_CRON_TIMING__MOD_DEVICE__DECIDEONWHATTODO,
		dev.decideOnWhatToDo,
	)
	if err != nil {
		return err
	}

	dev.cronJobs["decideOnWhatToDo"] = entryID

	return nil
}
