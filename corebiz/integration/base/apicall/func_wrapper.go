package apicall

type FuncWrapper struct {
	FnCode string 
	FnSig  func(any)(any, APICallEventIntface)

	FnList []func(any)(any, APICallEventIntface)

	FnArgs []any
	FnRets []any
}

type FuncWrapperIntface interface {
	
}
