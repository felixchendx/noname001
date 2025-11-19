package service

import (
	"noname001/app/base/messaging"
)

// _00nnn ~ _09nnn is reserved
// _90nnn ~ _99nnn is reserved
var (
	STRDB_ERR_00001 = messaging.NewMessageTemplate("DVC.SVC.STRDB_ERR_00001", "Store - DB has encountered internal error. Event ID: [%s].")
)

// device__basic.go
var (
	DVCBSC_NTC_11101 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_NTC_11101", "New Device '%s' added.")
	DVCBSC_ERR_11501 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_ERR_11501", "Failed to add Device. Code required.")
	DVCBSC_ERR_11502 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_ERR_11502", "Failed to add Device. Hostname required.")
	DVCBSC_ERR_11503 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_ERR_11503", "Failed to add Device. Username required.")
	DVCBSC_ERR_11504 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_ERR_11504", "Failed to add Device. Password required.")
	DVCBSC_ERR_11505 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_ERR_11505", "Failed to add Device. Brand required.")
	DVCBSC_ERR_11511 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_ERR_11511", "Failed to add Device. Code contains illegal char '%s'. Legal chars: %s.")
	DVCBSC_ERR_11551 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_ERR_11551", "Failed to add Device. Code '%s' is already used.")

	DVCBSC_NTC_14101 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_NTC_14101", "Device '%s' updated.")
	DVCBSC_ERR_14501 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_ERR_14501", "Failed to edit Device. Device '%s' not found.")
	DVCBSC_ERR_14502 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_ERR_14502", "Failed to edit Device. Hostname required.")
	DVCBSC_ERR_14503 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_ERR_14503", "Failed to edit Device. Username required.")
	DVCBSC_ERR_14504 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_ERR_14504", "Failed to edit Device. Password required.")
	DVCBSC_ERR_14505 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_ERR_14505", "Failed to edit Device. Brand required.")
	DVCBSC_ERR_14551 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_ERR_14551", "Failed to edit Device. Code '%s' is already used.")

	DVCBSC_ERR_13001 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_ERR_13001", "Device with id '%s' does not exist.")
	DVCBSC_ERR_14001 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_ERR_14001", "Device with code '%s' does not exist.")
	DVCBSC_ERR_15001 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_ERR_15001", "Device with id '%s' does not exist.")
	DVCBSC_NTC_15002 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_NTC_15002", "Device '%s' updated.")
	DVCBSC_ERR_16001 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_ERR_16001", "Device with id '%s' does not exist.")
	DVCBSC_NTC_16002 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_NTC_16002", "Device '%s' deleted.")

	DVCBSC_WRN_90501 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_WRN_90501", "Encryption routine failed. Unable to encrypt password.")
	DVCBSC_WRN_90502 = messaging.NewMessageTemplate("DVC.SVC.DVCBSC_WRN_90502", "Decryption routine failed. Unable to decrypt password.")
)

// cdt__001.go
var (
	CDT_ERR_11501 = messaging.NewMessageTemplate("DVC.SVC.CDT_ERR_11501", "Device '%s' is currently inactive.")
	CDT_ERR_11502 = messaging.NewMessageTemplate("DVC.SVC.CDT_ERR_11502", "Device '%s', failed to Get Device Info. See next messages for detail.")

	CDT_ERR_14501 = messaging.NewMessageTemplate("DVC.SVC.CDT_ERR_14501", "Device '%s' is currently inactive.")
	CDT_ERR_14502 = messaging.NewMessageTemplate("DVC.SVC.CDT_ERR_14502", "Device '%s', failed to Get Stream Info. See next messages for detail.")

	CDT_ERR_90501 = messaging.NewMessageTemplate("DVC.SVC.CDT_ERR_90501", "Device Identifier ID or Code is required.")
	CDT_ERR_90502 = messaging.NewMessageTemplate("DVC.SVC.CDT_ERR_90502", "Device with id '%s' does not exist.")
	CDT_ERR_90503 = messaging.NewMessageTemplate("DVC.SVC.CDT_ERR_90503", "Device with code '%s' does not exist.")
)

func (svc *Service) registerDeviceMessageTemplates() {}
