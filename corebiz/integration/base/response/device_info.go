package response

type DeviceInfoWrapper struct {
	DeviceInfo   *DeviceInfo

	OriginalData map[string]any
}

type DeviceInfo struct {
	DeviceName           string
	DeviceID             string
	Model                string
	// SerialNumber         string
	// MACAddress           string
	// FirmwareVersion      string
	// FirmwareReleasedDate string
	// EncoderVersion       string
	// EncoderReleasedDate  string
	DeviceType           string
	// TelecontrolID        int   
	// HardwareVersion      string
	// Manufacturer         string
	// CustomizedInfo       string
}
