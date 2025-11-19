package typing

import (
	"fmt"
)

type UnknownResponse struct {
	Status  string
	Message string
}

// =============== VVV conform to apicall.APIErrorIntface VVV =============== //
func (resp *UnknownResponse) SimpleError() string {
	return fmt.Sprintf("(%s) %s", resp.Status, resp.Message)
}

func (resp *UnknownResponse) FullError() string {
	return resp.SimpleError()
}
// =============== ^^^ conform to apicall.APIErrorIntface ^^^ =============== //
