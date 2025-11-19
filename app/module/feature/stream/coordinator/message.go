package coordinator

import (
	"noname001/app/base/messaging"
)

var (
	STRDB_ERR_00001 = messaging.NewMessageTemplate("STRM.CDT.STRDB_ERR_00001", "Store - DB has encountered internal error. Event ID: [%s].")
)

// post_start.go
var (
	// POST_ERR_01001 = messaging.NewMessageTemplate("STRM.CDT.POST_ERR_01001", "Failed to load autostart streams.")
)

// stream.go
var (
	// STRM_ERR_01001 = messaging.NewMessageTemplate("STRM.CDT.STRM_ERR_01001", "Failed to load Stream '%s'.")
	// STRM_ERR_01002 = messaging.NewMessageTemplate("STRM.CDT.STRM_ERR_01002", "Failed to load Stream '%s'. Referenced Stream Item does not exist.")
	// STRM_ERR_01003 = messaging.NewMessageTemplate("STRM.CDT.STRM_ERR_01003", "Failed to load Stream '%s'. Referenced Stream Group does not exist.")
	// STRM_ERR_01004 = messaging.NewMessageTemplate("STRM.CDT.STRM_ERR_01004", "Failed to load Stream '%s'. Referenced Stream Profile does not exist.")
	// STRM_ERR_01005 = messaging.NewMessageTemplate("STRM.CDT.STRM_ERR_01005", "Failed to load Stream '%s'. err: %s.")

	// STRM_ERR_99002 = messaging.NewMessageTemplate("STRM.CDT.STRM_ERR_99002", "Failed to load Stream '%s'. Unsupported SourceType '%s'.")
)

// svc__001.go
var (
	// SVC_ERR_01001 = messaging.NewMessageTemplate("STRM.CDT.SVC_ERR_01001", "Failed to load list Stream Item by Stream Group '%s'.")
	// SVC_ERR_01002 = messaging.NewMessageTemplate("STRM.CDT.SVC_ERR_01001", "Failed to load list Stream Item by Stream Profile '%s'.")
)
