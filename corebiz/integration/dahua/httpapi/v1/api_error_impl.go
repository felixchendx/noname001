package v1

import "fmt"

// =============== VVV conform to apicall.APIErrorIntface VVV =============== //
func (rs *TXT_ResponseStatus) SimpleError() string {
	return fmt.Sprintf("(%s) %s", rs.StatusCode, rs.StatusMsg)
}

func (rs *TXT_ResponseStatus) FullError() string {
	stringTemplate := `%s|%d|%s`
	args := []any{
		rs.RequestURL,
		rs.StatusCode,
		rs.StatusMsg,
	}
	return fmt.Sprintf(stringTemplate, args...)
}
// =============== ^^^ conform to apicall.APIErrorIntface ^^^ =============== //
