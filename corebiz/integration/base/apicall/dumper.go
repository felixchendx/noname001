package apicall

import (
	"fmt"

	"noname001/logging"
)

type APICallDumper struct {
	logger *logging.WrappedLogger

	dumpTo string
}

func newAPICallDumper(logger *logging.WrappedLogger) (*APICallDumper) {
	dumper := &APICallDumper{}
	dumper.logger = logger
	dumper.dumpTo = "stdout" // TODO

	return dumper
}

func (dumper *APICallDumper) Dump(collector *APICallEventCollector) {
	hasDataToDump := false

	for _, aceI := range collector.items {
		if aceI == nil { continue }

		if aceI.HasSerializedData() {
			hasDataToDump = true
			break
		}
	}

	if !hasDataToDump { return }

	formattedDump := dumper.formatHeader(collector)

	// TODO: to file
	dumper.logger.Debug(formattedDump)
}

func (dumper *APICallDumper) formatHeader(collector *APICallEventCollector) (string) {
	formatTemplate := `
===
=== === APICall - dump begin === ===
FnCode     : %s
FnCount    : %v
Items      :

`
	formatArgs := []any{
		collector.name,
		len(collector.items),
		// collector.SucceedIdx,
	}

	formattedMain := fmt.Sprintf(formatTemplate, formatArgs...)

	for idx, item := range collector.items {
		if item == nil { continue }

		switch typedEv := item.(type) {
		case *APICallEvent: formattedMain += dumper.formatItem(idx, typedEv)
		case *APICallEventBundle: formattedMain += dumper.formatBundle(idx, typedEv)
		default:
			formattedMain += fmt.Sprintf(`
=== ev #%v ===
unimplemented ev type: %T
`, idx, item)
		}
	}
	
	formattedMain += `
=== === APICall - dump end === ===
===
`

	return formattedMain
}

func (dumper *APICallDumper) formatBundle(idx int, evBundle *APICallEventBundle) (string) {
	if evBundle == nil {
		return fmt.Sprintf(`
=== ev bundle #%v - begin ===
NO DATA
`, idx)
	}

	formatTemplate := `

=== ev bundle #%v ===
Bundle Name : %s
Event Items :
`
	formatArgs := []any{
		idx,
		evBundle.BundleName(), 
	}

	formattedBundle := fmt.Sprintf(formatTemplate, formatArgs...)

	evItems := evBundle.Items()
	if evItems != nil {
		for itemIdx, item := range evItems {
			formattedBundle += dumper.formatItem(itemIdx, item)
		}
	}

	formattedBundle += fmt.Sprintf(`
=== ev bundle #%v - end ===

`, idx)

	return formattedBundle
}

func (dumper *APICallDumper) formatItem(idx int, ev *APICallEvent) (string) {
	if ev == nil {
		return fmt.Sprintf(`
=== ev #%v ===
NO DATA
`, idx)
	}

	formatTemplate := `
=== ev #%v ===
ID             : %s
Name           : %s
Error          : %s
Begin          : %v
End            : %v
Elapsed        : %v
SerializedData :
`
	formatArgs := []any{
		idx,
		ev.EventID(),
		ev.EventName(),
		ev.GoError(),
		ev.BeginTimestamp(),
		ev.EndTimestamp(),
		ev.Elapsed(),
	}

	formattedItem := fmt.Sprintf(formatTemplate, formatArgs...)

	serDat := ev.SerializedData()
	if serDat != nil {
		for idx, dat := range serDat {
			formattedItem += fmt.Sprintf(`--- dat[%v] ---
%s
`, idx, dat)
		}
	}

	return formattedItem
}
