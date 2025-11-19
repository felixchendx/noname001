package v1

import (
	"fmt"
	"reflect"
	"strings"

	"noname001/corebiz/integration/panasonic/netcam/v1/typing"
)

const (
	_DEFAULT_CUSTOM_MAPPER_TAG = "dcmt"
)

var (
	typeCache = make(map[reflect.Type][]reflect.StructField)
)

type (
	t_parsingStat struct {
		lineNum     int
		lineContent string
		reason      string
	}

	t_mappingStat struct {
		k, v      string
		fieldName string
		reason    string
	}

	t_decodingReport struct {
		err        error
		structType string

		parsingStats []t_parsingStat
		mappingStats []t_mappingStat
	}
)

// generic is nice, but costly...
// limited usage, for flat struct only
// yet to develop for nested / complex structs
func (api *APIClient) decodeSuccessResponse(respBody string, structPtr any) (*t_decodingReport) {
	var (
		decodingReport = &t_decodingReport{}
		parsingStats   = make([]t_parsingStat, 0)
		mappingStats   = make([]t_mappingStat, 0)
	)

	var (
		structPtrType       reflect.Type          = reflect.TypeOf(structPtr)
		structType          reflect.Type          = structPtrType.Elem()
		structVisibleFields []reflect.StructField

		valuePtr            reflect.Value         = reflect.ValueOf(structPtr)
		valueElem           reflect.Value         = valuePtr.Elem()
	)

	var (
		lines = strings.Split(respBody, "\r\n") // [CR][LF]
		kvs   = make(map[string]string)
	)

	cachedStructFields, inTypeCache := typeCache[structType]
	if !inTypeCache {
		cachedStructFields = reflect.VisibleFields(structType)
		typeCache[structType] = cachedStructFields
	}
	structVisibleFields = cachedStructFields


	for lineNum, line := range lines {
		if line == "" {
			parsingStats = append(parsingStats, t_parsingStat{
				lineNum: lineNum, lineContent: line,
				reason: "empty line",
			})
			continue
		}

		lineKV := strings.Split(line, "=")
		if len(lineKV) != 2 {
			// first line is non empty line with byte [239 187 191] ???
			parsingStats = append(parsingStats, t_parsingStat{
				lineNum: lineNum, lineContent: line,
				reason: "non-kv line",
			})
			continue
		}

		kvs[lineKV[0]] = lineKV[1]
	}

	for _, _structField := range structVisibleFields {
		tagMetadata, tagFound := _structField.Tag.Lookup(_DEFAULT_CUSTOM_MAPPER_TAG)
		if !tagFound {
			continue
		}

		mapperKey := tagMetadata
		respValue, respValueFound := kvs[mapperKey]
		if !respValueFound {
			mappingStats = append(mappingStats, t_mappingStat{
				k: mapperKey, v: "",
				fieldName: _structField.Name,
				reason: "key not in response",
			})
			continue
		}
		delete(kvs, mapperKey)

		valueElemField := valueElem.FieldByName(_structField.Name)
		if !valueElemField.IsValid() {
			mappingStats = append(mappingStats, t_mappingStat{
				k: mapperKey, v: respValue,
				fieldName: _structField.Name,
				reason: "field not valid",
			})
			continue
		}
		if !valueElemField.CanSet() {
			mappingStats = append(mappingStats, t_mappingStat{
				k: mapperKey, v: respValue,
				fieldName: _structField.Name,
				reason: "field not settable",
			})
			continue
		}

		switch valueElemField.Kind() {
		case reflect.String:
			valueElemField.SetString(respValue)

		default:
			mappingStats = append(mappingStats, t_mappingStat{
				k: mapperKey, v: respValue,
				fieldName: _structField.Name,
				reason: fmt.Sprintf("unimplemented type: %s", valueElemField.Kind().String()),
			})
			continue
		}
	}

	decodingReport.structType   = structType.String()
	decodingReport.parsingStats = parsingStats
	decodingReport.mappingStats = mappingStats

	return decodingReport
}

// TODO: parse html body
// https://stackoverflow.com/questions/30109061/golang-parse-html-extract-all-content-with-body-body-tags
// https://pkg.go.dev/golang.org/x/net/html#Parse
func (api *APIClient) decodeFailedResponse(respBody string, structPtr *typing.FailedResponse) {

}

func (api *APIClient) decodeUnknownResponse(respBody string, structPtr *typing.UnknownResponse) {

}
