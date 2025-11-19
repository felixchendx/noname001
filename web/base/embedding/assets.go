package embedding

import (
	"embed"

	"github.com/valyala/fasthttp"
)

var (
	//go:embed all:assets
	_assetFS embed.FS

	assetFS  *fasthttp.FS

	assetFSHandler fasthttp.RequestHandler
)

func FakeSingletonBaseAssetHandler() (fasthttp.RequestHandler) {
	if assetFS == nil {
		assetFS = &fasthttp.FS{
			GenerateIndexPages: true, // TODO: prod set to false

			FS: _assetFS,

			// TODO: moar configs, especially cache dur
			// https://pkg.go.dev/github.com/valyala/fasthttp#FS
		}

		assetFSHandler = assetFS.NewRequestHandler()
	}

	return assetFSHandler
}
