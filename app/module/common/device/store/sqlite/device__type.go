package sqlite

type DevicePE struct {
	ID    string
	Code  string
	Name  string
	State string
	Note  string

	Protocol string
	Hostname string
	Port     string
	Username string
	Password string
	Brand    string

	FallbackRTSPPort string

	CreatedTs string
	UpdatedTs string
}

// TODO: does unimported / local method consume extra memory ?
//       when this object is still reffed outside this package ?
func (entity *DevicePE) fullScan() ([]any) {
	return []any{
		&entity.ID,
		&entity.Code,
		&entity.Name,
		&entity.State,
		&entity.Note,

		&entity.Protocol,
		&entity.Hostname,
		&entity.Port,
		&entity.Username,
		&entity.Password,
		&entity.Brand,

		&entity.FallbackRTSPPort,

		&entity.CreatedTs, &entity.UpdatedTs,
	}
}
