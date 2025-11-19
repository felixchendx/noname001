package sys

import (
	"noname001/app/base/messaging"
)

// _00nnn ~ _09nnn is reserved
// _90nnn ~ _99nnn is reserved
// _xxnnn ~ _xxnnn where xx is running number of logical block, i.e. all msg from a function

// _nn000 ~ _nn099 is reserved
// _nn100 ~ _nn199 is reserved for notices
// _nn200 ~ _nn299 is reserved
// _nn300 ~ _nn399 is reserved for warnings
// _nn500 ~ _nn799 is reserved for errors

var (
	STRDB_ERR_00001 = messaging.NewMessageTemplate("SYS.SVC.STRDB_ERR_00001", "Store - DB has encountered internal error. Event ID: [%s].")
)

// sys_svc__user.go
var (
	SYSUSER_NTC_11101 = messaging.NewMessageTemplate("SYS.SVC.SYSUSER_NTC_11101", "New User '%s' added.")
	SYSUSER_ERR_11501 = messaging.NewMessageTemplate("SYS.SVC.SYSUSER_ERR_11501", "Failed to add User. Username required.")
	SYSUSER_ERR_11502 = messaging.NewMessageTemplate("SYS.SVC.SYSUSER_ERR_11502", "Failed to add User. Password required.")
	SYSUSER_ERR_11551 = messaging.NewMessageTemplate("SYS.SVC.SYSUSER_ERR_11551", "Failed to add User. Username '%s' is taken.")
	SYSUSER_ERR_11601 = messaging.NewMessageTemplate("SYS.SVC.SYSUSER_ERR_11601", "Failed to add User. Password hashing failed.") // TODO: useless message to end user

	SYSUSER_ERR_12501 = messaging.NewMessageTemplate("SYS.SVC.SYSUSER_ERR_12501", "User '%s' not found.")

	SYSUSER_ERR_13501 = messaging.NewMessageTemplate("SYS.SVC.SYSUSER_ERR_13501", "User '%s' not found.")

	SYSUSER_NTC_14101 = messaging.NewMessageTemplate("SYS.SVC.SYSUSER_NTC_14101", "User '%s' updated.")
	SYSUSER_ERR_14501 = messaging.NewMessageTemplate("SYS.SVC.SYSUSER_ERR_14501", "Failed to edit User. User '%s' not found.")
	SYSUSER_ERR_14502 = messaging.NewMessageTemplate("SYS.SVC.SYSUSER_ERR_14502", "Failed to edit User. User '%s' not found.") // fake message if is editing superadmin
	SYSUSER_ERR_14601 = messaging.NewMessageTemplate("SYS.SVC.SYSUSER_ERR_14601", "Failed to edit User. Password hashing failed.") // TODO: useless message to end user

	SYSUSER_NTC_15101 = messaging.NewMessageTemplate("SYS.SVC.SYSUSER_NTC_15101", "User '%s' deleted.")
	SYSUSER_ERR_15501 = messaging.NewMessageTemplate("SYS.SVC.SYSUSER_ERR_15501", "Failed to delete User. User '%s' not found.")
	SYSUSER_ERR_15502 = messaging.NewMessageTemplate("SYS.SVC.SYSUSER_ERR_15502", "Failed to delete User. User '%s' not found.") // fake message if is deleting superadmin
)

// sys_svc__session.go
var (
	SYSSESS_NTC_11101 = messaging.NewMessageTemplate("SYS.SVC.SYSSESS_NTC_11101", "Password updated.")
	SYSSESS_ERR_11501 = messaging.NewMessageTemplate("SYS.SVC.SYSSESS_ERR_11501", "Failed to update Password. Session already expired.")
	SYSSESS_ERR_11502 = messaging.NewMessageTemplate("SYS.SVC.SYSSESS_ERR_11502", "Failed to update Password. User not found.") // TODO: possible race condition. on del user, remove existing session 
	SYSSESS_ERR_11503 = messaging.NewMessageTemplate("SYS.SVC.SYSSESS_ERR_11503", "Failed to update Password. Old Password does not match.")
	SYSSESS_ERR_11504 = messaging.NewMessageTemplate("SYS.SVC.SYSSESS_ERR_11504", "Failed to update Password. Password hashing failed.")

)