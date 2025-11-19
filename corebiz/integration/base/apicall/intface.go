package apicall

type APICallEventIntface interface {
	IsConsideredError() (bool)

	IsGoError()         (bool)
	GoError()           (error)

	IsAPIError()        (bool)
	APIError()          (APIErrorIntface)

	Error()             (string)

	HasSerializedData() (bool)
}

type APIErrorIntface interface {
	SimpleError() (string)
	FullError()   (string)
}
