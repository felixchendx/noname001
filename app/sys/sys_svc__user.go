package sys

import (
	"strings"
	
	"github.com/google/uuid"

	"noname001/app/base/messaging"

	"noname001/app/sys/sqlite"
)

type SearchPagination = sqlite.SearchPagination

type SysUserDE struct {
	ID         string
	Username   string
	Password   string
	RoleSimple string
}
func (de *SysUserDE) sanitize() {
	de.Username   = strings.TrimSpace(de.Username)
	de.Password   = strings.TrimSpace(de.Password)
	de.RoleSimple = strings.TrimSpace(de.RoleSimple)
}
func (de *SysUserDE) toPE() (*sqlite.SysUserPE) {
	pe := &sqlite.SysUserPE{}
	pe.ID         = de.ID
	pe.Username   = de.Username
	pe.Password   = de.Password
	pe.RoleSimple = de.RoleSimple

	return pe
}
func (de *SysUserDE) fromPE(pe *sqlite.SysUserPE) (*SysUserDE) {
	if pe == nil { return nil }

	de = &SysUserDE{}
	de.ID         = pe.ID
	de.Username   = pe.Username
	de.Password   = "" // redacted
	de.RoleSimple = pe.RoleSimple
	
	return de
}

// === search ===
type SysUser__SearchCriteria = sqlite.SysUser__SearchCriteria
type SysUser__SearchResult struct {
	Data       []*SysUserDE

	Pagination *SearchPagination
}

func (sr *SysUser__SearchResult) fromStore(_sr *sqlite.SysUser__SearchResult) (*SysUser__SearchResult) {
	sr.Data       = make([]*SysUserDE, 0, len(_sr.Data))
	sr.Pagination = _sr.Pagination

	for _, pe := range _sr.Data {
		sr.Data = append(sr.Data, (&SysUserDE{}).fromPE(pe))
	}

	return sr
}

func (svc *SystemService) SysUser__EmptyItem() (*SysUserDE) {
	return &SysUserDE{}
}
func (svc *SystemService) SysUser__EmptyList() ([]*SysUserDE) {
	return make([]*SysUserDE, 0)
}

func (svc *SystemService) SysUser__Search(sc *SysUser__SearchCriteria) (*SysUser__SearchResult, *messaging.Messages) {
	messages := messaging.NewMessages()

	sr, dbev := svc.store.DB.SysUser__Search(sc)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	return (&SysUser__SearchResult{}).fromStore(sr), messages
}

func (svc *SystemService) SysUser__Add(de *SysUserDE) (*SysUserDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	// sanitizing
	de.sanitize()

	// standard validating
	if de.Username == "" { messages.AddError(SYSUSER_ERR_11501.NewMessage()) }
	if de.Password == "" { messages.AddError(SYSUSER_ERR_11502.NewMessage()) }
	if messages.HasError() { return nil, messages }

	// TODO: valid role checks, length checks, illegal char checks
	// contextual validating
	existingSysUser, dbev1 := svc.store.DB.SysUser__GetByUsername(de.Username)
	switch {
	case dbev1.IsError()       : messages.AddError(STRDB_ERR_00001.NewMessage(dbev1.EventID()))
	case existingSysUser != nil: messages.AddError(SYSUSER_ERR_11551.NewMessage(de.Username))
	}
	if messages.HasError() { return nil, messages }

	hashBytes, hashErr := svc.hashPassword([]byte(de.Password))
	if hashErr != nil {
		messages.AddError(SYSUSER_ERR_11601.NewMessage())
		return nil, messages
	}
	de.Password = string(hashBytes)

	// defaulting
	de.ID = uuid.New().String()
	if de.RoleSimple == "" { de.RoleSimple = "viewer" }

	pe, dbev := svc.store.DB.SysUser__AtomicInsert(de.toPE())
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	messages.AddNotice(SYSUSER_NTC_11101.NewMessage(pe.Username))
	return de.fromPE(pe), messages
}

func (svc *SystemService) SysUser__Get(id string) (*SysUserDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.SysUser__Get(id)
	switch {
	case dbev.IsError(): messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
	case pe == nil     : messages.AddError(SYSUSER_ERR_12501.NewMessage(id))
	}
	if messages.HasError() { return nil, messages }

	return (&SysUserDE{}).fromPE(pe), messages
}

// func (svc *SystemService) SysUser__GetByUsername(username string) (*SysUserDE, *messaging.Messages) {
// 	messages := messaging.NewMessages()

// 	pe, dbev := svc.store.DB.SysUser__GetByUsername(username)
// 	switch {
// 	case dbev.IsError(): messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
// 	case pe == nil     : messages.AddError(SYSUSER_ERR_13501.NewMessage(username))
// 	}
// 	if messages.HasError() { return nil, messages }

// 	return (&SysUserDE{}).fromPE(pe), messages
// }

func (svc *SystemService) SysUser__Edit(id string, de *SysUserDE) (*SysUserDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.SysUser__Get(id)
	switch {
	case dbev.IsError()             : messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
	case pe == nil                  : messages.AddError(SYSUSER_ERR_14501.NewMessage(id))
	case pe.Username == "superadmin": messages.AddError(SYSUSER_ERR_14502.NewMessage(id))
	}
	if messages.HasError() { return nil, messages }

	// sanitizing
	de.sanitize()

	// standard validating
	// TODO: valid role checks, length checks
	if messages.HasError() { return nil, messages }

	if de.Password != "" {
		hashBytes, hashErr := svc.hashPassword([]byte(de.Password))
		if hashErr != nil {
			messages.AddError(SYSUSER_ERR_14601.NewMessage())
			return nil, messages
		}
		de.Password = string(hashBytes)
	}

	// patching editable fields
	if de.Password != ""   { pe.Password = de.Password } 
	if de.RoleSimple != "" { pe.RoleSimple = de.RoleSimple }

	pe, dbev = svc.store.DB.SysUser__AtomicUpdate(pe)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	messages.AddNotice(SYSUSER_NTC_14101.NewMessage(pe.Username))
	return de.fromPE(pe), messages
}

func (svc *SystemService) SysUser__Delete(id string) (*messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.SysUser__Get(id)
	switch {
	case dbev.IsError()             : messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
	case pe == nil                  : messages.AddError(SYSUSER_ERR_15501.NewMessage(id))
	case pe.Username == "superadmin": messages.AddError(SYSUSER_ERR_15502.NewMessage(id))
	}
	if messages.HasError() { return messages }

	// TODO: usual stuffs

	if messages.HasError() { return messages }

	dbev = svc.store.DB.SysUser__AtomicDelete(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return messages
	}

	messages.AddNotice(SYSUSER_NTC_15101.NewMessage(pe.Username))
	return messages
}
