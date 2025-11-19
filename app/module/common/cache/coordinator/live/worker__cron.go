package live

import (
	nodeConst "noname001/node/constant"
)

func (lc *LiveCache) setupCrons() {
	lc.cronJobs["FULL_SCAN"], _ = lc.cron.AddFunc(
		nodeConst.CROSSNODE_CRON_TIMING__MOD_CACHE__FULL_SCAN,
		func() {
			lc.miscJobChan <- &t_miscJob{misc_job__lc_full_scan, nil}
		},
	)

	lc.cronJobs["CLEANUP"], _ = lc.cron.AddFunc(
		nodeConst.CROSSNODE_CRON_TIMING__MOD_CACHE__CLEANUP,
		func() {
			lc.miscJobChan <- &t_miscJob{misc_job__lc_cleanup, nil}
		},
	)

	lc.cronJobs["debug"], _ = lc.cron.AddFunc(
		nodeConst.CROSSNODE_CRON_TIMING__MOD_CACHE__DEBUG_PRINT,
		func() {
			lc.dump()
		},
	)
}
