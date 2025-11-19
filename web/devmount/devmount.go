package devmount

import (
	"os"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	devmountFS        *fasthttp.FS
	devmountFSHandler fasthttp.RequestHandler
)

func FakeSingletonDevMountHandler(ctx *fasthttp.RequestCtx) {
	if devmountFS == nil {
		devmountFS = &fasthttp.FS{
			GenerateIndexPages: true,

			FS:          os.DirFS("web/devmount"),
			PathRewrite: fasthttp.NewPathSlashesStripper(1),

			// hmm, no effect on cache ?
			CacheDuration: 0 * time.Second,
			SkipCache:     true,
		}

		devmountFSHandler = devmountFS.NewRequestHandler()
	}

	devmountFSHandler(ctx)
}
