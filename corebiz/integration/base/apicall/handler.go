package apicall

import (
	"noname001/logging"
)

// type APICallThingy struct {
// 	FnCode string
// 	FnSignature func(string)(bool, APICallEventIntface)

// 	FnList []func(string)(bool, APICallEventIntface)
// }

type APICallHandler struct {
	logger *logging.WrappedLogger

	// Registry map[string]*APICallThingy

	succeedFunctions map[string]int

	dumper *APICallDumper
}

func NewHandler(logger *logging.WrappedLogger) (*APICallHandler) {
	handler := &APICallHandler{}
	handler.logger = logger
	handler.succeedFunctions = make(map[string]int)
	// handler.Registry = make(map[string]*APICallThingy)
	handler.dumper = newAPICallDumper(handler.logger)
	
	return handler
}

func (handler *APICallHandler) Call(fnCode string) {

}

func (handler *APICallHandler) MarkSucceedFunction(fnCode string, fnIdx int) {
	handler.succeedFunctions[fnCode] = fnIdx
}
func (handler *APICallHandler) RetrieveMarkedSucceedFunction(fnCode string) (fnIdx int, hasMarked bool) {
	fnIdx, hasMarked = handler.succeedFunctions[fnCode]
	return
}

// func (handler *APICallHandler) Register(fnCode string) {
// 	thingy := &APICallThingy{}
// 	thingy.FnCode = fnCode
// 	// thingy.FnSignature = func(bool)(bool)
// 	thingy.FnList = make([]func(string)(bool, APICallEventIntface), 0)
// }

func (handler *APICallHandler) SpawnCollector(fnCode string, fnCount int) (*APICallEventCollector) {
	return newCollector(fnCode, fnCount)
}

func (handler *APICallHandler) RetrieveCollector(collector *APICallEventCollector) {
	handler.dumper.Dump(collector)
}
