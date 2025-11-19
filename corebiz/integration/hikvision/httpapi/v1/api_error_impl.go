package v1

import (
	"fmt"
)

// =============== VVV conform to apicall.APIErrorIntface VVV =============== //
func (rs *XML_ResponseStatus) SimpleError() string {
	return fmt.Sprintf("(%s) %s", rs.ErrorCode, rs.ErrorMsg)
}

func (rs *XML_ResponseStatus) FullError() string {
	addStatList := rs.AdditionalErr.StatusList
	addStatString := "[]"
	if len(addStatList) > 0 {
		addStatString = "["
		for _, addStat := range addStatList {
			inArgs := []any{addStat.ID, addStat.StatusCode, addStat.StatusString, addStat.SubStatusCode}
			addStatString += fmt.Sprintf("\n%s|%d|%s|%s", inArgs)
		}
		addStatString += fmt.Sprintf("\n]")
	}

	stringTemplate := `%s|%d|%s|%d|%s|%d|%s|%s`
	args := []any{
		rs.RequestURL,
		rs.StatusCode,
		rs.StatusString,
		rs.ID,
		rs.SubStatusCode,
		rs.ErrorCode,
		rs.ErrorMsg,
		addStatString,
	}

	return fmt.Sprintf(stringTemplate, args...)
}
// =============== ^^^ conform to apicall.APIErrorIntface ^^^ =============== //
