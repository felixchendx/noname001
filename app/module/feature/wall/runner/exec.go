package runner

import (
	"noname001/app/constant"
)

func (runna *Runner) exec() (err error) {
	if runna.runnerConfig == nil { return }
	if !runna.runnerConfig.Execute { return }

	runna.execSegmentWall(runna.runnerConfig.SegmentWall)

	return
}

func (runna *Runner) execSegmentWall(seg SegmentWall) {
	if seg.Mode == "" || seg.Mode == constant.RUNNER__MODE_NONE { return }

	for _, wall := range seg.Walls {
		runna.walkWall(seg.Mode, wall)
	}
}
func (runna *Runner) walkWall(mode string, item Wall) {
	if item.Code == "" { return }

	refDE, refMessages := runna.service.WallLayout__FindByCode(item.WallLayoutCode)
	if refDE == nil {
		runna.logger.Errorf("%s: %s", runna.logPrefix, refMessages)
		return
	}

	de, messages := runna.service.Wall__FindByCode(item.Code, false)

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
	default:
	}

	switch {
	case doAdd:
		de = runna.service.Wall__Empty()
		de.Code  = item.Code
		de.Name  = item.Name
		de.State = item.State
		de.Note  = item.Note

		de.WallLayoutID = refDE.ID

		for _, wallItem := range item.Items {
			itemDE := runna.service.WallItem__Empty()
			itemDE.Index        = wallItem.Index
			itemDE.SourceNodeID = wallItem.SourceNodeID
			itemDE.StreamCode   = wallItem.StreamCode
			de.Items = append(de.Items, itemDE)
		}

		de, messages = runna.service.Wall__Add(de)
		if messages.HasError() {
			runna.logger.Errorf("%s: %s", runna.logPrefix, messages)
			return
		}

		runna.logger.Infof("%s: %s", runna.logPrefix, messages)

	case doEdit:
		de.Code  = item.Code
		de.Name  = item.Name
		de.State = item.State
		de.Note  = item.Note

		de.WallLayoutID = refDE.ID
		de.Items        = runna.service.WallItem__EmptyList()

		for _, wallItem := range item.Items {
			itemDE := runna.service.WallItem__Empty()
			itemDE.Index        = wallItem.Index
			itemDE.SourceNodeID = wallItem.SourceNodeID
			itemDE.StreamCode   = wallItem.StreamCode
			de.Items = append(de.Items, itemDE)
		}

		de, messages = runna.service.Wall__Edit(de.ID, de)
		if messages.HasError() {
			runna.logger.Errorf("%s: %s", runna.logPrefix, messages)
			return
		}

		runna.logger.Infof("%s: %s", runna.logPrefix, messages)

	case doDelete:
		messages = runna.service.Wall__Delete(de.ID)
		if messages.HasError() {
			runna.logger.Errorf("%s: %s", runna.logPrefix, messages)
			return
		}

		runna.logger.Infof("%s: %s", runna.logPrefix, messages)

	default:
	}
}
