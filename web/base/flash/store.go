package flash

import (
	"github.com/valyala/fasthttp"

	"noname001/logging"
	"noname001/app/base/messaging"

	"noname001/web/base/cookie"
)

// TODO: data race review

type FlashStoreParams struct {
	Logger      *logging.WrappedLogger
	LogPrefix   string

	CookieStore *cookie.CookieStore
}
type FlashStore struct {
	logger      *logging.WrappedLogger
	logPrefix   string

	cookieStore *cookie.CookieStore

	items       map[string]*FlashItem // map[wsid]*FlashItem
}

func NewFlashStore(params *FlashStoreParams) (*FlashStore) {
	store := &FlashStore{}
	store.logger = params.Logger
	store.logPrefix = params.LogPrefix + ".flash"
	store.cookieStore = params.CookieStore

	store.items = make(map[string]*FlashItem)

	return store
}

func (store *FlashStore) GetFlashBundle(ctx *fasthttp.RequestCtx) (*FlashBundle) {
	flashItem := store.getFlashItem(ctx)
	
	flashBundle := &FlashBundle{}
	flashBundle.Prev = *flashItem
	flashBundle.Next = flashItem.reset()

	return flashBundle
}

func (store *FlashStore) getFlashItem(ctx *fasthttp.RequestCtx) (*FlashItem) {
	wsid := store.cookieStore.GetWebSessionID(ctx)

	flashItem, ok := store.items[wsid]
	if !ok {
		flashItem = &FlashItem{
			wsid: wsid,
			Data: make(map[string]map[string]any),
			Messages: messaging.NewMessages(),
			PlainMessages: make(map[string]string),
		}
		store.items[wsid] = flashItem
	}

	return flashItem
}

type FlashBundle struct {
	Prev FlashItem
	Next *FlashItem
}

type FlashItem struct {
	wsid          string

	Data          FlashData
	Messages      *messaging.Messages
	PlainMessages FlashPlainMessages

	// lastAccessed time.Time // TODO: delete stale item
}
func (item *FlashItem) reset() (*FlashItem) {
	item.Data = make(map[string]map[string]any)
	item.Messages = messaging.NewMessages()
	item.PlainMessages = make(map[string]string) 

	return item
}
func (item *FlashItem) HasData() (bool) {
	// TODO: temp, the right way to check map
	hasItem := false
	for _, v := range item.Data {
		_ = v
		hasItem = true
		break
	}
	return hasItem
}
func (item *FlashItem) HasMessage() (bool) {
	return item.Messages.HasMessage()
}

type FlashData          map[string]map[string]any
type FlashPlainMessages map[string]string

func (data FlashData) Has(key string) (bool) {
	_, ok := data[key]
	return ok
}
func (data FlashData) Set(key string, val map[string]any) {
	data[key] = val
}
func (data FlashData) Get(key string) (val map[string]any) {
	dat, ok := data[key]
	if ok {
		return dat
	}
	return map[string]any{}
}
