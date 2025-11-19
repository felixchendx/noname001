package auth

import (
	"noname001/logging"

	"noname001/web/base/cookie"
	"noname001/web/base/flash"
)

type AuthProviderParams struct {
	Logger      *logging.WrappedLogger
	LogPrefix   string
	CookieStore *cookie.CookieStore
	FlashStore  *flash.FlashStore
}
type AuthProvider struct {
	logger    *logging.WrappedLogger
	logPrefix string
	cookieStore *cookie.CookieStore
	flashStore  *flash.FlashStore
}

func NewAuthProvider(params *AuthProviderParams) (*AuthProvider) {
	authProvider := &AuthProvider{}
	authProvider.logger = params.Logger
	authProvider.logPrefix = params.LogPrefix + ".auth"
	authProvider.cookieStore = params.CookieStore
	authProvider.flashStore = params.FlashStore

	return authProvider
}
