package coordinator

import (
	"noname001/app/base/messaging"
)

var (
	STRDB_ERR_00001 = messaging.NewMessageTemplate("WALL.CDT.STRDB_ERR_00001", "Store - DB has encountered internal error. Event ID: [%s].")
)
