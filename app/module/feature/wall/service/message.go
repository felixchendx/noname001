package service

import (
	"noname001/app/base/messaging"
)

// _00nnn ~ _09nnn is reserved
// _90nnn ~ _99nnn is reserved
var (
	STRDB_ERR_00001 = messaging.NewMessageTemplate("WALL.SVC.STRDB_ERR_00001", "Store - DB has encountered internal error. Event ID: [%s].")
)

// wall_layout__basic.go
var (
	WLLLYTBSC_ERR_13001 = messaging.NewMessageTemplate("WALL.SVC.WLLLYTBSC_ERR_13001", "Wall Layout with id '%s' does not exist.")
	WLLLYTBSC_ERR_14001 = messaging.NewMessageTemplate("WALL.SVC.WLLLYTBSC_ERR_14001", "Wall Layout with code '%s' does not exist.")
)

// wall__nonstd.go
var (
	WLLNONSTD_NTC_12001 = messaging.NewMessageTemplate("WALL.SVC.WLLNONSTD_NTC_12001", "New Wall '%s' added.")
	WLLNONSTD_ERR_12501 = messaging.NewMessageTemplate("WALL.SVC.WLLNONSTD_ERR_12501", "Failed to add Wall. Code required.")
	WLLNONSTD_ERR_12511 = messaging.NewMessageTemplate("WALL.SVC.WLLNONSTD_ERR_12511", "Failed to add Wall. Code contains illegal char '%s'. Legal chars: %s.")
	WLLNONSTD_ERR_13501 = messaging.NewMessageTemplate("WALL.SVC.WLLNONSTD_ERR_13501", "Wall with id '%s' does not exist.")
	WLLNONSTD_ERR_14501 = messaging.NewMessageTemplate("WALL.SVC.WLLNONSTD_ERR_14501", "Wall with code '%s' does not exist.")
	WLLNONSTD_NTC_15001 = messaging.NewMessageTemplate("WALL.SVC.WLLNONSTD_NTC_15001", "Wall '%s' updated.")
	WLLNONSTD_ERR_15501 = messaging.NewMessageTemplate("WALL.SVC.WLLNONSTD_ERR_15501", "Failed to update. Wall with id '%s' does not exist.")
	WLLNONSTD_NTC_16001 = messaging.NewMessageTemplate("WALL.SVC.WLLNONSTD_NTC_16001", "Wall '%s' deleted.")
	WLLNONSTD_ERR_16501 = messaging.NewMessageTemplate("WALL.SVC.WLLNONSTD_ERR_16501", "Failed to delete. Wall with id '%s' does not exist.")

	WLLNONSTD_NTC_17001 = messaging.NewMessageTemplate("WALL.SVC.WLLNONSTD_NTC_17001", "Wall item '#%d' updated.")
	WLLNONSTD_ERR_17501 = messaging.NewMessageTemplate("WALL.SVC.WLLNONSTD_ERR_17501", "Failed to update. Wall item with id '%s' does not exist.")
	WLLNONSTD_ERR_18501 = messaging.NewMessageTemplate("WALL.SVC.WLLNONSTD_ERR_17801", "Wall item with id '%s' does not exist.")
)