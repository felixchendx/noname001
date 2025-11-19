package embedding

import (
	"embed"

	"github.com/valyala/fasthttp"
)

var (
	//go:embed all:assets
	_assetFS embed.FS

	assetFS *fasthttp.FS

	assetFSHandler fasthttp.RequestHandler
)

func FakeSingletonAssetHandler() (fasthttp.RequestHandler) {
	if assetFS == nil {
		assetFS = &fasthttp.FS{
			GenerateIndexPages: true, // TODO: set to false for prod
			FS: _assetFS,

			// stripping "/wall" from "/wall/assets/"
			// to match the embedded fs paths "/assets/**"
			PathRewrite: fasthttp.NewPathSlashesStripper(1),
	
			// TODO: moar configs, especially cache dur
			// https://pkg.go.dev/github.com/valyala/fasthttp#FS
		}

		assetFSHandler = assetFS.NewRequestHandler()
	}

	return assetFSHandler
}
