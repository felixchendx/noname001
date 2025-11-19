package sys

import (
	"context"
	"sync"
	"time"
)

// TODO: implement sync.pool
// TODO: reorganize session related stuffs

type SysSessionType string
type SysSimpleRole  string

type SysUser struct {
	ID         string
	Username   string
	SimpleRole SysSimpleRole
}

type SysSession struct {
	id   string
	Type SysSessionType

	User *SysUser

	ttl        int64
	lastAccess int64

	// Inbox 
}

const (
	SessionTypeAnonymous SysSessionType = "anonymous"
	SessionTypeCli       SysSessionType = "cli"
	SessionTypeWeb       SysSessionType = "web"

	SimpleRoleSuperadmin SysSimpleRole = "superadmin"
	SimpleRoleAdmin      SysSimpleRole = "admin"
	SimpleRoleOperator   SysSimpleRole = "operator"
	SimpleRoleViewer     SysSimpleRole = "viewer"
)

var (
	sessionRegistryMutex sync.Mutex
	sessionRegistry      map[string]*SysSession = make(map[string]*SysSession)
)

func newSysUser(id, username, role string) (*SysUser) {
	simpleRole := SimpleRoleViewer

	switch role {
	case "superadmin": simpleRole = SimpleRoleSuperadmin
	case "admin": simpleRole = SimpleRoleAdmin
	case "operator": simpleRole = SimpleRoleOperator
	case "viewer": simpleRole = SimpleRoleViewer
	default: simpleRole = SimpleRoleViewer
	}
	
	return &SysUser{id, username, simpleRole}
}

func (sess *SysSession) ID() (string) {
	return sess.id
}

func (sess *SysSession) IsAnonymous() (bool) {
	return sess.Type == SessionTypeAnonymous
}
func (sess *SysSession) IsAuthenticated() (bool) {
	return sess.User != nil
}

func (sess *SysSession) IsSuperadmin() (bool) {
	return sess.User != nil && sess.User.SimpleRole == SimpleRoleSuperadmin
}
func (sess *SysSession) IsAdmin() (bool) {
	return sess.User != nil && sess.User.SimpleRole == SimpleRoleAdmin
}
func (sess *SysSession) HasAdminAuthorization() (bool) {
	switch {
	case sess.IsSuperadmin(): return true
	case sess.IsAdmin()     : return true
	}
	return false
}
func (sess *SysSession) DoesNotHaveAdminAuthorization() (bool) {
	return !sess.HasAdminAuthorization()
}

func (sess *SysSession) IsOperator() (bool) {
	return sess.User != nil && sess.User.SimpleRole == SimpleRoleOperator
}
func (sess *SysSession) IsViewer() (bool) {
	return sess.User != nil && sess.User.SimpleRole == SimpleRoleViewer
}

func sessionCleanupWorker(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer func() {
			ticker.Stop()
		}()

		WorkerLoop:
		for {
			select {
			case <- ctx.Done():
				break WorkerLoop
			case <- ticker.C:
				sessionCleanup()
			}
		}
	}()
}

func sessionCleanup() {
	naw := time.Now().Unix()
	cleanupList := make([]string, 0)

	for sessID, sess := range sessionRegistry {
		if naw - sess.lastAccess > sess.ttl {
			cleanupList = append(cleanupList, sessID)
		}
	}

	sessionRegistryMutex.Lock()
	for _, sessID := range cleanupList {
		delete(sessionRegistry, sessID)
	}
	sessionRegistryMutex.Unlock()
}
