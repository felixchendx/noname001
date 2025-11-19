package service

import (
	"noname001/app/base/messaging"
)

// _00nnn ~ _09nnn is reserved
// _90nnn ~ _99nnn is reserved
var (
	STRDB_ERR_00001 = messaging.NewMessageTemplate("STRM.SVC.STRDB_ERR_00001", "Store - DB has encountered internal error. Event ID: [%s].")
)

// stream_profile__basic.go
var (
	STRMPRFLBSC_NTC_11101 = messaging.NewMessageTemplate("STRM.SVC.STRMPRFLBSC_NTC_11101", "New Stream Profile '%s' added.")
	STRMPRFLBSC_ERR_11501 = messaging.NewMessageTemplate("STRM.SVC.STRMPRFLBSC_ERR_11501", "Failed to add Stream Profile. Code required.")
	STRMPRFLBSC_ERR_11511 = messaging.NewMessageTemplate("STRM.SVC.STRMPRFLBSC_ERR_11511", "Failed to add Stream Profile. Code contains illegal char '%s'. Legal chars: %s.")
	STRMPRFLBSC_ERR_11551 = messaging.NewMessageTemplate("STRM.SVC.STRMPRFLBSC_ERR_11551", "Failed to add Stream Profile. Code '%s' is already used.")

	STRMPRFLBSC_NTC_14101 = messaging.NewMessageTemplate("STRM.SVC.STRMPRFLBSC_NTC_14101", "Stream Profile '%s' updated.")
	STRMPRFLBSC_ERR_14501 = messaging.NewMessageTemplate("STRM.SVC.STRMPRFLBSC_ERR_14501", "Failed to edit Stream Profile. Stream Profile '%s' not found.")

	STRMPRFLBSC_NTC_02001 = messaging.NewMessageTemplate("STRM.SVC.STRMPRFLBSC_NTC_02001", "New Stream Profile added.")
	STRMPRFLBSC_ERR_03001 = messaging.NewMessageTemplate("STRM.SVC.STRMPRFLBSC_ERR_03001", "Stream Profile with id '%s' does not exist.")
	STRMPRFLBSC_ERR_04001 = messaging.NewMessageTemplate("STRM.SVC.STRMPRFLBSC_ERR_04001", "Stream Profile with code '%s' does not exist.")
	STRMPRFLBSC_ERR_05001 = messaging.NewMessageTemplate("STRM.SVC.STRMPRFLBSC_ERR_05001", "Stream Profile with id '%s' does not exist.")
	STRMPRFLBSC_NTC_05002 = messaging.NewMessageTemplate("STRM.SVC.STRMPRFLBSC_NTC_05002", "Stream Profile '%s' updated.")
	STRMPRFLBSC_ERR_06001 = messaging.NewMessageTemplate("STRM.SVC.STRMPRFLBSC_ERR_06001", "Stream Profile with id '%s' does not exist.")
	STRMPRFLBSC_NTC_06002 = messaging.NewMessageTemplate("STRM.SVC.STRMPRFLBSC_NTC_06002", "Stream Profile '%s' deleted.")

	STRMPRFLBSC_NTC_06010 = messaging.NewMessageTemplate("STRM.SVC.STRMPRFLBSC_NTC_06010", "Cannot delete this stream profile. It's in use by stream group(s): %s")
)

// stream_group__basic.go
var (
	STRMGRPBSC_NTC_11101 = messaging.NewMessageTemplate("STRM.SVC.STRMGRPBSC_NTC_11101", "New Stream Group '%s' added.")
	STRMGRPBSC_ERR_11501 = messaging.NewMessageTemplate("STRM.SVC.STRMGRPBSC_ERR_11501", "Failed to add Stream Group. Code required.")
	STRMGRPBSC_ERR_11511 = messaging.NewMessageTemplate("STRM.SVC.STRMGRPBSC_ERR_11511", "Failed to add Stream Group. Code contains illegal char '%s'. Legal chars: %s.")
	STRMGRPBSC_ERR_11551 = messaging.NewMessageTemplate("STRM.SVC.STRMGRPBSC_ERR_11551", "Failed to add Stream Group. Code '%s' is already used.")

	STRMGRPBSC_NTC_14101 = messaging.NewMessageTemplate("STRM.SVC.STRMGRPBSC_NTC_14101", "Stream Group '%s' updated.")
	STRMGRPBSC_ERR_14501 = messaging.NewMessageTemplate("STRM.SVC.STRMGRPBSC_ERR_14501", "Failed to edit Stream Group. Stream Group '%s' not found.")


	STRMGRPBSC_NTC_02001 = messaging.NewMessageTemplate("STRM.SVC.STRMGRPBSC_NTC_02001", "New Stream Group added.")
	STRMGRPBSC_ERR_03001 = messaging.NewMessageTemplate("STRM.SVC.STRMGRPBSC_ERR_03001", "Stream Group with id '%s' does not exist.")
	STRMGRPBSC_ERR_04001 = messaging.NewMessageTemplate("STRM.SVC.STRMGRPBSC_ERR_04001", "Stream Group with code '%s' does not exist.")
	STRMGRPBSC_ERR_05001 = messaging.NewMessageTemplate("STRM.SVC.STRMGRPBSC_ERR_05001", "Stream Group with id '%s' does not exist.")
	STRMGRPBSC_NTC_05002 = messaging.NewMessageTemplate("STRM.SVC.STRMGRPBSC_NTC_05002", "Stream Group '%s' updated.")
	STRMGRPBSC_ERR_06001 = messaging.NewMessageTemplate("STRM.SVC.STRMGRPBSC_ERR_06001", "Stream Group with id '%s' does not exist.")
	STRMGRPBSC_NTC_06002 = messaging.NewMessageTemplate("STRM.SVC.STRMGRPBSC_NTC_06002", "Stream Group '%s' deleted.")
)

// stream_item__basic.go
var (
	STRMITMBSC_NTC_11101 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_NTC_11101", "New Stream Item '%s' added.")
	STRMITMBSC_ERR_11501 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_11501", "Failed to add Stream Item. Code required.")
	STRMITMBSC_ERR_11502 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_11502", "Failed to add Stream Item. Source Type required.")
	STRMITMBSC_ERR_11503 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_11503", "Failed to add Stream Item. Device Code required.")
	STRMITMBSC_ERR_11504 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_11504", "Failed to add Stream Item. Device Channel ID required.")
	STRMITMBSC_ERR_11505 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_11505", "Failed to add Stream Item. External URL required.")
	STRMITMBSC_ERR_11506 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_11506", "Failed to add Stream Item. Filepath required.")
	STRMITMBSC_ERR_11511 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_11511", "Failed to add Stream Item. Code contains illegal char '%s'. Legal chars: %s.")
	STRMITMBSC_ERR_11551 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_11551", "Failed to add Stream Item. Code '%s' is already used.")

	STRMITMBSC_NTC_14101 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_NTC_14101", "Stream Item '%s' updated.")
	STRMITMBSC_ERR_14501 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_14501", "Failed to edit Stream Item. Stream Item '%s' not found.")
	STRMITMBSC_ERR_14502 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_14502", "Failed to edit Stream Item. Source Type required.")
	STRMITMBSC_ERR_14503 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_14503", "Failed to edit Stream Item. Device Code required.")
	STRMITMBSC_ERR_14504 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_14504", "Failed to edit Stream Item. Device Channel ID required.")
	STRMITMBSC_ERR_14505 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_14505", "Failed to edit Stream Item. External URL required.")
	STRMITMBSC_ERR_14506 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_14506", "Failed to edit Stream Item. Filepath required.")

	STRMITMBSC_NTC_02001 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_NTC_02001", "New Stream Item added.")
	STRMITMBSC_ERR_03001 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_03001", "Stream Item with id '%s' does not exist.")
	STRMITMBSC_ERR_04001 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_04001", "Stream Item with code '%s' does not exist.")
	STRMITMBSC_ERR_05001 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_05001", "Stream Item with id '%s' does not exist.")
	STRMITMBSC_NTC_05002 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_NTC_05002", "Stream Item '%s' updated.")
	STRMITMBSC_ERR_06001 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_ERR_06001", "Stream Item with id '%s' does not exist.")
	STRMITMBSC_NTC_06002 = messaging.NewMessageTemplate("STRM.SVC.STRMITMBSC_NTC_06002", "Stream Item '%s' deleted.")
)

// comm__strm.g0
var (
	COMMSTRM_ERR_11001 = messaging.NewMessageTemplate("STRM.SVC.COMMSTRM_ERR_11001", "Stream Item with code '%s' is currently inactive.")
	COMMSTRM_ERR_11002 = messaging.NewMessageTemplate("STRM.SVC.COMMSTRM_ERR_11002", "Stream Item with code '%s' is currently not streaming.")

	COMMSTRM_ERR_90001 = messaging.NewMessageTemplate("STRM.SVC.COMMSTRM_ERR_90001", "Stream Item Identifier ID or Code is required.")
	COMMSTRM_ERR_90002 = messaging.NewMessageTemplate("STRM.SVC.COMMSTRM_ERR_90002", "Stream Item with id '%s' does not exist.")
	COMMSTRM_ERR_90003 = messaging.NewMessageTemplate("STRM.SVC.COMMSTRM_ERR_90003", "Stream Item with code '%s' does not exist.")
)
