package apicall

type APICallEventCollector struct {
	name  string
	
	items []APICallEventIntface
}

func newCollector(name string, itemCount int) (*APICallEventCollector) {
	collector := &APICallEventCollector{
		name: name,
		items: make([]APICallEventIntface, 0, itemCount),
	}

	for i := 0; i < itemCount; i++ {
		collector.items = append(collector.items, nil)
	}

	return collector
}

func (collector *APICallEventCollector) Collect(itemIdx int, aceI APICallEventIntface) {
	collector.items[itemIdx] = aceI
}
