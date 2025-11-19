package typing

import (
	"fmt"
)

type FailedResponse struct {
	Status  string
	Message string
}

// =============== VVV conform to apicall.APIErrorIntface VVV =============== //
func (resp *FailedResponse) SimpleError() string {
	return fmt.Sprintf("(%s) %s", resp.Status, resp.Message)
}

func (resp *FailedResponse) FullError() string {
	return resp.SimpleError()
}
// =============== ^^^ conform to apicall.APIErrorIntface ^^^ =============== //
