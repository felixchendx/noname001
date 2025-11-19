package web

import (
	"bytes"

	"github.com/valyala/fasthttp"

	"noname001/web/devmount"
)

func (svc *WebService) requestHandler(ctx *fasthttp.RequestCtx) {
	path := bytes.TrimRight(ctx.Path(), "/")

	switch {
	case bytes.HasPrefix(path, []byte("/assets")):
		svc.baseBundle.RouteHandler.RouteAsset(ctx)
		return
	}

	svc.baseBundle.Auth.InjectWebSession(ctx)

	// TODO: move standard guards here after separating private and public area

	// TODO: implement routing with adaptive radix trie, https://en.wikipedia.org/wiki/Radix_tree
	// prioritize non prefixed / specific paths
	switch {
	case bytes.Equal(path, []byte("")):
		redirectTo := "/home"
		if svc.baseBundle.Auth.IsLoggedOut(ctx) { redirectTo = "/login" }
		ctx.Redirect(redirectTo, fasthttp.StatusFound)

	case bytes.Equal(path, []byte("/home")):
		webSess := svc.baseBundle.Auth.RetrieveWebSession(ctx)

		redirectTo := ""
		if webSess.IsLoggedOut() {
			redirectTo = "/login"
		} else {
			switch {
			case webSess.IsSuperadmin(), webSess.IsAdmin():
				// redirectTo = "/admin/dashboard"
				redirectTo = "/monitoring/temp-dashboard"

			default:
				redirectTo = "/user/dashboard"
			}
		}
		ctx.Redirect(redirectTo, fasthttp.StatusFound)

	case bytes.Equal(path, []byte("/login")):
		svc.baseBundle.RouteHandler.RouteLogin(ctx)

	case bytes.Equal(path, []byte("/logout")):
		svc.baseBundle.RouteHandler.RouteLogout(ctx)

	case bytes.Equal(path, []byte("/change-password")):
		svc.baseBundle.RouteHandler.RouteChangePassword(ctx)


	case bytes.HasPrefix(path, []byte("/admin")):
		svc.adminBundle.RouteHandler.RootHandler(ctx)

	case bytes.HasPrefix(path, []byte("/user")):
		svc.userBundle.RouteHandler.RootHandler(ctx)


	case bytes.HasPrefix(path, []byte("/device")):
		svc.deviceBundle.RouteHandler.RootHandler(ctx)

	case bytes.HasPrefix(path, []byte("/stream")):
		svc.streamBundle.RouteHandler.RootHandler(ctx)

	case bytes.HasPrefix(path, []byte("/wall")):
		svc.wallBundle.RouteHandler.RootHandler(ctx)

	case bytes.HasPrefix(path, []byte("/monitoring")):
		svc.monitoringBundle.RouteHandler.RootHandler(ctx)


	case bytes.HasPrefix(path, []byte("/devmount")):
		isProd := false
		if isProd {
			svc.baseBundle.RouteHandler.Route404(ctx)
		} else {
			devmount.FakeSingletonDevMountHandler(ctx)
		}

	default:
		svc.baseBundle.RouteHandler.Route404(ctx)
	}

	// shenanigans
}
