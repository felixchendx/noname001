package sys

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"noname001/app/base/messaging"
)

const (
	SysSessionErr01001 = "SYS-SESS.ERR.01001"
	SysSessionErr01002 = "SYS-SESS.ERR.01002"
	SysSessionNtc01001 = "SYS-SESS.NTC.01002"
)

var (
	persistenceStoreError001 = messaging.NewMessageTemplate("STR.PRS.ERR.001", "Persistence Store has encountered internal error. Event ID: [%s].")

	msgSysSessionErr01001 = messaging.NewMessageTemplate(SysSessionErr01001, "Session with id '%s' not found.")
	msgSysSessionErr01002 = messaging.NewMessageTemplate(SysSessionErr01002, "Username and password combination not found.")
	msgSysSessionNtc01001 = messaging.NewMessageTemplate(SysSessionNtc01001, "Session authenticated!")
)

func (svc *SystemService) NewSession() (*SysSession) {
	sess := &SysSession{
		id: uuid.New().String(),

		Type: SessionTypeAnonymous,
		User: nil,

		ttl    : 60 * 60 * 8, // second
		lastAccess: time.Now().Unix(),
	}

	sessionRegistryMutex.Lock()
	sessionRegistry[sess.id] = sess
	sessionRegistryMutex.Unlock()

	return sess
}

func (svc *SystemService) FindSession(sessID string) (*SysSession) {
	sessionRegistryMutex.Lock()
	item, ok := sessionRegistry[sessID]
	sessionRegistryMutex.Unlock()
	if ok {
		item.lastAccess = time.Now().Unix()
		return item
	}

	return nil
}

func (svc *SystemService) AuthenticateSession(sessID, username, password string) (*messaging.Messages) {
	var messages *messaging.Messages = messaging.NewMessages()
	
	sess := svc.FindSession(sessID)
	if sess == nil {
		return messaging.OneLinerError(msgSysSessionErr01001.NewMessage(sessID))
	}

	sysUserPE, dbev := svc.store.DB.SysUser__GetByUsername(username)
	if dbev.IsError() {
		return messaging.OneLinerError(persistenceStoreError001.NewMessage(dbev.EventID()))
	}

	if sysUserPE == nil {
		return messaging.OneLinerError(msgSysSessionErr01002.NewMessage())
	}

	passCheckErr := svc.checkPasswordHash([]byte(sysUserPE.Password), []byte(password))
	if passCheckErr != nil {
		return messaging.OneLinerError(msgSysSessionErr01002.NewMessage())
	}


	sess.Type = SessionTypeCli
	sess.User = newSysUser(sysUserPE.ID, sysUserPE.Username, sysUserPE.RoleSimple)

	messages.AddNotice(msgSysSessionNtc01001.NewMessage())
	return messages
}

func (svc *SystemService) TerminateSession(sessID string) {
	sess := svc.FindSession(sessID)
	sessionRegistryMutex.Lock()
	if sess != nil {
		delete(sessionRegistry, sessID)
	}
	sessionRegistryMutex.Unlock()
}

func (svc *SystemService) WebLogin(username, password string) (*SysSession, *messaging.Messages) {
	var messages *messaging.Messages = messaging.NewMessages()

	sysUserPE, dbev := svc.store.DB.SysUser__GetByUsername(username)
	if dbev.IsError() {
		return nil, messaging.OneLinerError(persistenceStoreError001.NewMessage(dbev.EventID()))
	}

	if sysUserPE == nil {
		return nil, messaging.OneLinerError(msgSysSessionErr01002.NewMessage())
	}

	passCheckErr := svc.checkPasswordHash([]byte(sysUserPE.Password), []byte(password))
	if passCheckErr != nil {
		return nil, messaging.OneLinerError(msgSysSessionErr01002.NewMessage())
	}


	sess := svc.NewSession()
	sess.Type = SessionTypeWeb
	sess.User = newSysUser(sysUserPE.ID, sysUserPE.Username, sysUserPE.RoleSimple)

	messages.AddNotice(msgSysSessionNtc01001.NewMessage())
	return sess, messages
}
func (svc *SystemService) WebLogout(sessID string) {
	svc.TerminateSession(sessID)
}

func (svc *SystemService) WebChangePassword(sessID string, oldPass, newPass string) (*messaging.Messages) {
	messages := messaging.NewMessages()

	sysSess := svc.FindSession(sessID)
	if sysSess == nil {
		messages.AddError(SYSSESS_ERR_11501.NewMessage())
		return messages
	}

	sysUserPE, dbev := svc.store.DB.SysUser__Get(sysSess.User.ID)
	switch {
	case dbev.IsError():   messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
	case sysUserPE == nil: messages.AddError(SYSSESS_ERR_11502.NewMessage())
	}
	if messages.HasError() { return messages }

	passCheckErr := svc.checkPasswordHash([]byte(sysUserPE.Password), []byte(oldPass))
	if passCheckErr != nil {
		messages.AddError(SYSSESS_ERR_11503.NewMessage())
		return messages
	}

	hashBytes, hashErr := svc.hashPassword([]byte(newPass))
	if hashErr != nil {
		messages.AddError(SYSSESS_ERR_11504.NewMessage())
		return messages
	}

	dbev2 := svc.store.DB.SysUser__SetPassword(sysUserPE.ID, string(hashBytes))
	if dbev2.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return messages
	}

	messages.AddNotice(SYSSESS_NTC_11101.NewMessage())
	return messages
}

// TODO: max password length 72, limitation from bcrypt ?
func (svc *SystemService) hashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, 13)
}
func (svc *SystemService) checkPasswordHash(hash, password []byte) (error) {
	return bcrypt.CompareHashAndPassword(hash, password)
}
