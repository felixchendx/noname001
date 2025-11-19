package apicall

type APICallEventBundle struct {
	name           string
	
	items          []*APICallEvent
	
	partialSuccess bool
}

func NewBundle(name string) (*APICallEventBundle) {
	return &APICallEventBundle{
		name: name,
		items: make([]*APICallEvent, 0),
	}
}

func (evBundle *APICallEventBundle) BundleName() (string) {
	return evBundle.name
}

func (evBundle *APICallEventBundle) AddItem(item *APICallEvent) {
	evBundle.items = append(evBundle.items, item)
}

func (evBundle *APICallEventBundle) Items() ([]*APICallEvent) {
	return evBundle.items
}

func (evBundle *APICallEventBundle) MarkAsPartialSuccess() {
	evBundle.partialSuccess = true
}

func (evBundle *APICallEventBundle) IsPartialSuccess() (bool) {
	return evBundle.partialSuccess
}

func (evBundle *APICallEventBundle) HasError() (bool) {
	for _, item := range evBundle.items {
		if item.IsConsideredError() {
			return true
		}
	}

	return false
}


// ================= VVV conform to APICallEventIntface VVV ================= //
func (evBundle *APICallEventBundle) IsConsideredError() (bool) {
	return evBundle.HasError()
}

func (evBundle *APICallEventBundle) IsGoError() (bool) {
	return evBundle.GoError() != nil
}
func (evBundle *APICallEventBundle) GoError() (error) {
	for _, item := range evBundle.items {
		goErr := item.GoError()
		if goErr != nil {
			return goErr
		}
	}
	
	return nil 
}

func (evBundle *APICallEventBundle) IsAPIError() (bool) {
	return evBundle.APIError() != nil
}
func (evBundle *APICallEventBundle) APIError() (APIErrorIntface) {
	for _, item := range evBundle.items {
		apiErr := item.APIError()
		if apiErr != nil {
			return apiErr
		}
	}
	
	return nil
}

func (evBundle *APICallEventBundle) Error() (string) {
	for _, item := range evBundle.items {
		errString := item.Error()
		if errString != "" {
			return errString
		}
	}
	
	return ""
}

func (evBundle *APICallEventBundle) HasSerializedData() (bool) {
	for _, item := range evBundle.items {
		if item.HasSerializedData() {
			return true
		}
	}

	return false
}
// ================= ^^^ conform to APICallEventIntface ^^^ ================= //
