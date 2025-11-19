package util

import (
	"fmt"

	"github.com/valyala/fasthttp"

	// "noname001/app/base/messaging"

	webConstant "noname001/web/constant"
)

// this structure is intended only for simple internal usage
// where the programmer has access to codes of both api provider and api consumer
// is not suitable for setup where be and fe are separate entities

type LAPIDefaultResponse struct {
	Status  string `json:"status"` // 'ok' | 'error'
	Message string `json:"message"`
	Data    any    `json:"data"`

	Messages *LAPIResponseMessages `json:"messages"`
}

// TODO: CORS block, just in case...

func (util *Util) OnelinerOkLAPIResponse(data any, message string) (*LAPIDefaultResponse) {
	return &LAPIDefaultResponse{Status: "ok", Data: data, Message: message, Messages: nil}
}

func (util *Util) OnelinerErrorLAPIResponse(data any, message string) (*LAPIDefaultResponse) {
	return &LAPIDefaultResponse{Status: "error", Data: data, Message: message, Messages: nil}
}

func (util *Util) Oneliner500LAPIResponse(ctx *fasthttp.RequestCtx, err error) {
	tempMessage := fmt.Sprintf(`{"message": "TODO: internal server error occured! err: %s"}`, err.Error())
	
	ctx.SetContentType(webConstant.CONTENT_TYPE_JSON)
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	ctx.SetBody([]byte(tempMessage))
}


// func (lapi_df *LAPIDefaultResponse) SetMessages(messages *messaging.Messages) (*LAPIDefaultResponse) {
// 	lapi_df.Messages = LAPIResponseMessages{}

// 	notices := make([]LAPIResponseMessage, 0)
// 	for _, msg := range messages.Notices {
// 		notices = append(notices, LAPIResponseMessage{msg.Code, msg.Description()})
// 	}
// 	warnings := make([]LAPIResponseMessage, 0)
// 	for _, msg := range messages.Warnings {
// 		warnings = append(warnings, LAPIResponseMessage{msg.Code, msg.Description()})
// 	}
// 	errors := make([]LAPIResponseMessage, 0)
// 	for _, msg := range messages.Errors {
// 		errors = append(errors, LAPIResponseMessage{msg.Code, msg.Description()})
// 	}

// 	lapi_df.Messages.Notices = notices
// 	lapi_df.Messages.Warnings = warnings
// 	lapi_df.Messages.Errors = errors
	
// 	return lapi_df
// }

type LAPIResponseMessages struct {
	Notices  []LAPIResponseMessage
	Warnings []LAPIResponseMessage
	Errors   []LAPIResponseMessage
}
type LAPIResponseMessage struct {
	Code string
	Desc string
}
