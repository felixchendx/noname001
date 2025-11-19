package navigation

import (
	"bytes"

	"github.com/valyala/fasthttp"

	"noname001/logging"

	modDef "noname001/app/module/definition"

	webAuth "noname001/web/base/auth"
)

type BaseNaviParams struct {
	Logger    *logging.WrappedLogger
	LogPrefix string

	AuthProvider *webAuth.AuthProvider
	TempModStates map[string]string
}
type BaseNavi struct {
	logger    *logging.WrappedLogger
	logPrefix string

	authProvider *webAuth.AuthProvider

	tempModStates map[string]string

	siteNavLinks       map[string]NavLink

	mappedSiteNavLinks map[string]map[string]string
}

func NewBaseNavi(params *BaseNaviParams) (*BaseNavi) {
	navi := &BaseNavi{}
	navi.logger = params.Logger
	navi.logPrefix = params.LogPrefix

	navi.authProvider = params.AuthProvider

	navi.tempModStates = params.TempModStates

	navi.generateNavLinks()

	return navi
}

func (navi *BaseNavi) Has(navCode string) (bool) {
	_, navOk := navi.siteNavLinks[navCode]
	return navOk
}
func (navi *BaseNavi) DoesNotHave(navCode string) (bool) {
	return !navi.Has(navCode)
}

func (navi *BaseNavi) MappedSiteNavLinks(ctx *fasthttp.RequestCtx) (map[string]map[string]any) {
	currPath := bytes.TrimRight(ctx.Path(), "/")

	mapped := make(map[string]map[string]any)

	for navCode, navLink := range navi.siteNavLinks {
		mappedNavLink := navLink.AsDataMap()

		// TODO: comprehensive location marking + bread crumbs
		if bytes.HasPrefix(currPath, []byte(navLink.URI)) {
			mappedNavLink["is_shown"] = true
		}

		mapped[navCode] = mappedNavLink
	}

	// TODO: temp, reorganize
	webSess := navi.authProvider.RetrieveWebSession(ctx)

	if webSess.HasAdminAuthorization() {
		// temp
		delete(mapped, "user")
		delete(mapped, "user_dashboard")

	} else {
		delete(mapped, "admin")
		delete(mapped, "admin_dashboard")
		delete(mapped, "user_listing")
	}

	return mapped
}

func (navi *BaseNavi) generateNavLinks() {
	navi.siteNavLinks = map[string]NavLink{
		"login"  : NavLink{"login",  "/login",  "Login"},
		"logout" : NavLink{"logout", "/logout", "Logout"},
	}

	// admin area
	navi.siteNavLinks["admin"]           = NavLink{"admin",           "/admin",              "Admin"}
	navi.siteNavLinks["admin_dashboard"] = NavLink{"admin_dashboard", "/admin/dashboard",    "Admin - Dashboard"}
	navi.siteNavLinks["user_listing"]    = NavLink{"user_listing",    "/admin/user/listing", "Admin - User Listing"}

	// user area
	navi.siteNavLinks["user"]           = NavLink{"user",           "/user",              "User"}
	navi.siteNavLinks["user_dashboard"] = NavLink{"user_dashboard", "/user/dashboard",    "User - Dashboard"}

	if modState, modOk := navi.tempModStates[string(modDef.COMMON_DEVICE)]; modOk {
		if modState == "start" {
			navi.siteNavLinks["device"]           = NavLink{"device",           "/device",                "Device"}
			navi.siteNavLinks["device_dashboard"] = NavLink{"device_dashboard", "/device/dashboard",      "Device - Dashboard"}
			navi.siteNavLinks["device_listing"]   = NavLink{"device_listing",   "/device/device/listing", "Device - Device Listing"}
		}
	}

	if modState, modOk := navi.tempModStates[string(modDef.FEATURE_STREAM)]; modOk {
		if modState == "start" {
			navi.siteNavLinks["stream"]                 = NavLink{"stream",                 "/stream",                        "Stream"}
			navi.siteNavLinks["stream_dashboard"]       = NavLink{"stream_dashboard",       "/stream/dashboard",              "Stream - Dashboard"}
			navi.siteNavLinks["stream_profile_listing"] = NavLink{"stream_profile_listing", "/stream/stream-profile/listing", "Stream - Stream Profile Listing"}
			navi.siteNavLinks["stream_group_listing"]   = NavLink{"stream_group_listing",   "/stream/stream-group/listing",   "Stream - Stream Group Listing"}
		}
	}

	if modState, modOk := navi.tempModStates[string(modDef.FEATURE_WALL)]; modOk {
		if modState == "start" {
			navi.siteNavLinks["wall"]           = NavLink{"wall",           "/wall",              "Wall"}
			navi.siteNavLinks["wall_dashboard"] = NavLink{"wall_dashboard", "/wall/dashboard",    "Wall - Dashboard"}
			navi.siteNavLinks["wall_listing"]   = NavLink{"wall_listing",   "/wall/wall/listing", "Wall - Wall Listing"}
		}
	}

	navi.siteNavLinks["monitoring"] = NavLink{"monitoring", "/monitoring", "Monitoring"}
	navi.siteNavLinks["monitoring_temp_dashboard"] = NavLink{
		"monitoring_temp_dashboard",
		"/monitoring/temp-dashboard",
		"Monitoring - Temp Dashboard",
	}
}

func (navi *BaseNavi) InjectTempModStates(modStates map[string]string) {
	navi.tempModStates = modStates
	navi.generateNavLinks()
}
