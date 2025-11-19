package coordinator

import (
	"noname001/app/base/messaging"
)

var (
	STRDB_ERR_00001 = messaging.NewMessageTemplate("DVC.CDT.STRDB_ERR_00001", "Store - DB has encountered internal error. Event ID: [%s].")
)

// post_start.go
var (
	POST_ERR_01001 = messaging.NewMessageTemplate("DVC.CDT.POST_ERR_01001", "Failed to load autostart devices.")
)

// device.go
var (
	DVC_ERR_11001 = messaging.NewMessageTemplate("DVC.CDT.DVC_ERR_11001", "Failed to load Device with id '%s'.")
	DVC_ERR_11002 = messaging.NewMessageTemplate("DVC.CDT.DVC_ERR_11002", "Failed to load Device. Device with id '%s' does not exist.")
	DVC_ERR_11201 = messaging.NewMessageTemplate("DVC.CDT.DVC_ERR_11201", "Failed to load Device with code '%s'. err: %s.")
	DVC_ERR_11202 = messaging.NewMessageTemplate("DVC.CDT.DVC_ERR_11202", "Failed to patch and reload Device with code '%s'. err: %s.")
	DVC_ERR_11299 = messaging.NewMessageTemplate("DVC.CDT.DVC_ERR_11299", "Cannot load Device with code '%s'. Unsupported brand '%s'.")

	DVC_ERR_12501 = messaging.NewMessageTemplate("DVC.CDT.DVC_ERR_12501", "Failed to connect to LiveDevice with id '%s'. err: %s.")

	DVC_ERR_13501 = messaging.NewMessageTemplate("DVC.CDT.DVC_ERR_13501", "Device '%s', failed to cache. err: %s.")

	DVC_ERR_90501 = messaging.NewMessageTemplate("DVC.CDT.DVC_ERR_90501", "Device '%s' has been removed.")

	DVC_ERR_99501 = messaging.NewMessageTemplate("DVC.CDT.DVC_ERR_99501", "Programmer error occured. Lost ability to periodically ping device.")

	// temp, will be swept later
	DVC_WRN_80501 = messaging.NewMessageTemplate("DVC.SVC.DVC_WRN_80501", "Decryption routine failed. Unable to decrypt password.")
)

// svc__001.go
var (
	// these errors are to be used with accompanying error that identify which device is error
	SVC_ERR_90501 = messaging.NewMessageTemplate("DVC.CDT.SVC_ERR_90501", "Device has been removed.")
	SVC_ERR_90502 = messaging.NewMessageTemplate("DVC.CDT.SVC_ERR_90502", "Device request err: %s")
	SVC_ERR_90503 = messaging.NewMessageTemplate("DVC.CDT.SVC_ERR_90503", "Device hardware returned err: %s")
)