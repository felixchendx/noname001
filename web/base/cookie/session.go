package cookie

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

// TODO: encrypt decrypt ONCE!!!!
//       and store in custom context
//       current implementation enc dec whenever cookie value is accessed

// TODO: there's specialized cookie for session ?

const (
	SESSION_COOKIE_NAME     string = "_sess"
	SESSION_COOKIE_DOMAIN   string = "" 
	SESSION_COOKIE_PATH     string = "/"
	SESSION_COOKIE_SAMESITE        = fasthttp.CookieSameSiteStrictMode
	SESSION_COOKIE_MAXAGE   int    = 60 * 60 * 8
	SESSION_COOKIE_SECURE   bool   = false
	SESSION_COOKIE_HTTPONLY bool = true
)

func (store *CookieStore) newSessionCookie() (*fasthttp.Cookie) {
	sessCookie := &fasthttp.Cookie{}
	sessCookie.SetKey(SESSION_COOKIE_NAME)
	sessCookie.SetDomain(SESSION_COOKIE_DOMAIN)
	sessCookie.SetPath(SESSION_COOKIE_PATH)
	sessCookie.SetSameSite(SESSION_COOKIE_SAMESITE)
	sessCookie.SetMaxAge(SESSION_COOKIE_MAXAGE)
	sessCookie.SetSecure(SESSION_COOKIE_SECURE)
	sessCookie.SetHTTPOnly(SESSION_COOKIE_HTTPONLY)

	scv := &sessionCookieValue{
		wsid: uuid.New().String(),
		isid: "",
	}

	sessCookie.SetValueBytes(scv.serialize())

	return sessCookie
}

func (store *CookieStore) setSessionCookie(ctx *fasthttp.RequestCtx, scv *sessionCookieValue) {
	sessCookie := store.newSessionCookie()

	encryptedVal, err := store.EncryptCookie(scv.serialize())
	if err != nil {
		store.logger.Errorf("%s: encrypt cookie err - %s", store.logPrefix, err)
		encryptedVal = []byte("")
	}

	sessCookie.SetValueBytes(encryptedVal)

	ctx.Response.Header.SetCookie(sessCookie)
}

func (store *CookieStore) getSessionCookie(ctx *fasthttp.RequestCtx) (*sessionCookieValue) {
	cookieVal := ctx.Request.Header.Cookie(SESSION_COOKIE_NAME)

	var err error
	var decryptedVal []byte = []byte("")

	if string(cookieVal) == "" {
		// hmm, must be new visitor
	} else {
		decryptedVal, err = store.DecryptCookie(cookieVal)
		if err != nil {
			store.logger.Errorf("%s: decrypt cookie err - %s", store.logPrefix, err)
		}
	}

	scv := (&sessionCookieValue{}).deserialize(decryptedVal)

	if scv.wsid == "" {
		scv.wsid, scv.isid = uuid.New().String(), ""
		store.setSessionCookie(ctx, scv)
	}

	return scv
}

func (store *CookieStore) DestroySessionCookie(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.DelClientCookie(SESSION_COOKIE_NAME)
}

func (store *CookieStore) GetWebSessionID(ctx *fasthttp.RequestCtx) (string) {
	scv := store.getSessionCookie(ctx)
	return scv.wsid
}

func (store *CookieStore) GetInternalSessionID(ctx *fasthttp.RequestCtx) (string) {
	scv := store.getSessionCookie(ctx)
	return scv.isid
}

func (store *CookieStore) SetInternalSessionID(ctx *fasthttp.RequestCtx, isid string) {
	scv := store.getSessionCookie(ctx)
	scv.isid = isid
	store.setSessionCookie(ctx, scv)
}

type sessionCookieValue struct {
	wsid string
	isid string
}
func (scv *sessionCookieValue) serialize() ([]byte) {
	return []byte(fmt.Sprintf("%s::%s", scv.wsid, scv.isid))
}
func (scv *sessionCookieValue) deserialize(sessVal []byte) (*sessionCookieValue) {
	sessValParts := strings.Split(string(sessVal), "::")

	if len(sessValParts) != 2 {
		return scv
	}

	scv.wsid = sessValParts[0]
	scv.isid = sessValParts[1]

	return scv
}
