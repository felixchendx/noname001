package live

const (
	misc_job__lc_full_scan miscJobCode = "lc:full-scan"
	misc_job__lc_cleanup miscJobCode = "lc:cleanup"

	misc_job__node_not_seen miscJobCode = "node:not-seen"
)

type miscJobCode string
type t_miscJob struct {
	code   miscJobCode
	params []string
}

func (lc *LiveCache) liveCacheWorker() {

	looper:
	for {
		select {
		case <- lc.context.Done():
			break looper

		case miscJob := <- lc.miscJobChan:
			switch miscJob.code {
			case misc_job__lc_full_scan:
				lc.scanRoutine()

			case misc_job__lc_cleanup:
				lc.cleanupRoutine()

			case misc_job__node_not_seen:
				lc.refreshNode(miscJob.params[0])
			}
		}
	}
}
