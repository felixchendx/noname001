package sqlite

type StreamItemPE struct {
	ID               string
	StreamGroupID    string
	Code             string
	Name             string
	State            string
	Note             string

	SourceType       string
	DeviceCode       string
	DeviceChannelID  string
	DeviceStreamType string
	ExternalURL      string
	Filepath         string
	EmbeddedFilepath string

	CreatedTs        string
	UpdatedTs        string
}

type StreamItem__SearchCriteria struct {
	Code  *string
	State *string
}

func (db *DB) StreamItem__EntityFullScan(entity *StreamItemPE) ([]any) {
	return []any{
		&entity.ID,
		&entity.StreamGroupID,
		&entity.Code,
		&entity.Name,
		&entity.State,
		&entity.Note,

		&entity.SourceType,
		&entity.DeviceCode, &entity.DeviceChannelID, &entity.DeviceStreamType,
		&entity.ExternalURL,
		&entity.Filepath,
		&entity.EmbeddedFilepath,

		&entity.CreatedTs, &entity.UpdatedTs,
	}
}

