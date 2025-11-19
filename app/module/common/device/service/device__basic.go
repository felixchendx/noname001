package service 

import (
	"github.com/google/uuid"

	"noname001/app/base/messaging"
	"noname001/app/base/sanitation"
	"noname001/app/constant"

	deviceEv "noname001/app/module/common/device/event"
)

func (svc *Service) EmptyDevice() (*DeviceDE) {
	return &DeviceDE{}
}
func (svc *Service) EmptyDeviceList() ([]*DeviceDE) {
	return make([]*DeviceDE, 0)
}

func (svc *Service) AddDevice(de *DeviceDE) (*DeviceDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	// sanitizing
	de.sanitize()

	// standard validating
	if de.Code == "" {
		messages.AddError(DVCBSC_ERR_11501.NewMessage())
	} else {
		isIllegal, illegalChar := sanitation.Code_ContainsIllegalChar(de.Code)
		if isIllegal {
			messages.AddError(DVCBSC_ERR_11511.NewMessage(illegalChar, sanitation.CODE__LEGAL_CHARS))
		}
	}
	if de.Hostname == "" { messages.AddError(DVCBSC_ERR_11502.NewMessage()) }
	if de.Username == "" { messages.AddError(DVCBSC_ERR_11503.NewMessage()) }
	if de.Password == "" { messages.AddError(DVCBSC_ERR_11504.NewMessage()) }
	if de.Brand == ""    { messages.AddError(DVCBSC_ERR_11505.NewMessage()) }
	if messages.HasError() { return nil, messages }

	// TODO: other contextual validating such as valid hostname, illegal chars
	// contextual validating
	existingDevice, dbev1 := svc.store.DB.Device__GetByCode(de.Code)
	switch {
	case dbev1.IsError()      : messages.AddError(STRDB_ERR_00001.NewMessage(dbev1.EventID()))
	case existingDevice != nil: messages.AddError(DVCBSC_ERR_11551.NewMessage(de.Code))
	}
	if messages.HasError() { return nil, messages }

	if err := svc.deviceEncryptionRoutine(de); err != nil {
		messages.AddWarning(DVCBSC_WRN_90501.NewMessage())
	}

	// defaulting
	de.ID = uuid.New().String()
	if de.Name  == "" { de.Name = de.Code }
	if de.State == "" { de.State = constant.ENTITY__STATE_INACTIVE }

	if de.Protocol == "" { de.Protocol = "http" }
	if de.Port == ""     { de.Port = "80" }

	if de.FallbackRTSPPort == "" { de.FallbackRTSPPort = "554" }

	pe, dbev := svc.store.DB.Device__AtomicInsert(de.toPE())
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	de = de.fromPE(pe)

	if err := svc.deviceDecryptionRoutine(de); err != nil {
		messages.AddWarning(DVCBSC_WRN_90502.NewMessage())
	}

	svc.evHub.PublishDeviceEvent(
		deviceEv.DEVICE_EVENT_CODE__CREATE,
		de.ID, de.Code,
	)

	messages.AddNotice(DVCBSC_NTC_11101.NewMessage(de.Code))
	return de, messages
}

func (svc *Service) FindDevice(id string) (*DeviceDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.Device__Get(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if pe == nil {
		messages.AddError(DVCBSC_ERR_13001.NewMessage(id))
		return nil, messages
	}

	de := (&DeviceDE{}).fromPE(pe)

	if err := svc.deviceDecryptionRoutine(de); err != nil {
		messages.AddWarning(DVCBSC_WRN_90502.NewMessage())
	}

	return de, messages
}

func (svc *Service) FindDeviceByCode(code string) (*DeviceDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.Device__GetByCode(code)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if pe == nil {
		messages.AddError(DVCBSC_ERR_14001.NewMessage(code))
		return nil, messages
	}

	de := (&DeviceDE{}).fromPE(pe)

	if err := svc.deviceDecryptionRoutine(de); err != nil {
		messages.AddWarning(DVCBSC_WRN_90502.NewMessage())
	}

	return de, messages
}

// func (svc *Service) EditDevicePrecheck() (*messaging.Messages) {
// 	messages := messaging.NewMessages()
// 	return messages
// }

func (svc *Service) EditDevice(id string, de *DeviceDE) (*DeviceDE, *messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.Device__Get(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}
	if pe == nil {
		messages.AddError(DVCBSC_ERR_14501.NewMessage(id))
		return nil, messages
	}

	// sanitizing
	de.sanitize()

	// standard validating
	if de.Hostname == "" { messages.AddError(DVCBSC_ERR_14502.NewMessage()) }
	if de.Username == "" { messages.AddError(DVCBSC_ERR_14503.NewMessage()) }
	if de.Password == "" { messages.AddError(DVCBSC_ERR_14504.NewMessage()) }
	if de.Brand == ""    { messages.AddError(DVCBSC_ERR_14505.NewMessage()) }
	if messages.HasError() { return nil, messages }

	// TODO: contextual validating

	if err := svc.deviceEncryptionRoutine(de); err != nil {
		messages.AddWarning(DVCBSC_WRN_90501.NewMessage())
	}

	// merge editable fields
	pe.Name  = de.Name
	pe.State = de.State
	pe.Note  = de.Note

	pe.Protocol = de.Protocol
	pe.Hostname = de.Hostname
	pe.Port     = de.Port
	pe.Username = de.Username
	pe.Password = de.Password
	pe.Brand    = de.Brand

	pe.FallbackRTSPPort = de.FallbackRTSPPort

	// defaulting
	if pe.Name  == "" { pe.Name = pe.Code }
	if pe.State == "" { pe.State = constant.ENTITY__STATE_INACTIVE }

	if pe.Protocol == "" { pe.Protocol = "http" }
	if pe.Port == ""     { pe.Port = "80" }

	if pe.FallbackRTSPPort == "" { pe.FallbackRTSPPort = "554" }

	pe, dbev = svc.store.DB.Device__AtomicUpdate(pe)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return nil, messages
	}

	de = de.fromPE(pe)

	if err := svc.deviceDecryptionRoutine(de); err != nil {
		messages.AddWarning(DVCBSC_WRN_90502.NewMessage())
	}

	svc.evHub.PublishDeviceEvent(
		deviceEv.DEVICE_EVENT_CODE__UPDATE,
		de.ID, de.Code,
	)

	messages.AddNotice(DVCBSC_NTC_14101.NewMessage(pe.Code))
	return de, messages
}

// func (svc *Service) DeleteDevicePrecheck() (*messaging.Messages) {
// 	messages := messaging.NewMessages()
// 	return messages
// }

func (svc *Service) DeleteDevice(id string) (*messaging.Messages) {
	messages := messaging.NewMessages()

	pe, dbev := svc.store.DB.Device__Get(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return messages
	}
	if pe == nil {
		messages.AddError(DVCBSC_ERR_16001.NewMessage(id))
		return messages
	}

	// TODO: validation

	if messages.HasError() { return messages }

	dbev = svc.store.DB.Device__AtomicDelete(id)
	if dbev.IsError() {
		messages.AddError(STRDB_ERR_00001.NewMessage(dbev.EventID()))
		return messages
	}

	svc.evHub.PublishDeviceEvent(
		deviceEv.DEVICE_EVENT_CODE__DELETE,
		pe.ID, pe.Code,
	)

	messages.AddNotice(DVCBSC_NTC_16002.NewMessage(pe.Code))
	return messages
}

func (svc *Service) deviceEncryptionRoutine(de *DeviceDE) (error) {
	// encHostname, err := sec.Encrypt(de.Hostname, encryptionKey)
	// if err != nil { return err }

	// encUsername, err := sec.Encrypt(de.Username, encryptionKey)
	// if err != nil { return err }

	encPassword, err := svc.secBundle.Encrypt(de.Password, "device")
	if err != nil {
		svc.logger.Errorf("%s: deviceEncryptionRoutine - field password err, %w", svc.logPrefix, err)
		return err
	}

	// de.Hostname = encHostname
	// de.Username = encUsername
	de.Password = encPassword

	return nil
}

func (svc *Service) deviceDecryptionRoutine(de *DeviceDE) (error) {
	// decHostname, err := sec.Decrypt(de.Hostname, encryptionKey)
	// if err != nil { return err }

	// decUsername, err := sec.Decrypt(de.Username, encryptionKey)
	// if err != nil { return err }

	decPassword, err := svc.secBundle.Decrypt(de.Password, "device")
	if err != nil {
		svc.logger.Errorf("%s: deviceDecryptionRoutine - field password err, %w", svc.logPrefix, err)
		return err
	}

	// de.Hostname = decHostname
	// de.Username = decUsername
	de.Password = decPassword

	return nil
}
