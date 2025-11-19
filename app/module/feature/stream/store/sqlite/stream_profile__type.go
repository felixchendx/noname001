package sqlite

type StreamProfilePE struct {
	ID                     string
	Code                   string
	Name                   string
	State                  string
	Note                   string

	TargetVideoFPS         float64
	TargetVideoWidth       int
	TargetVideoHeight      int
	TargetVideoCodec       string
	TargetVideoCompression int
	TargetVideoBitrate     int

	TargetAudioCodec       string
	TargetAudioCompression int
	TargetAudioBitrate     int

	ShowTimestamp          string
	ShowVideoInfo          string
	ShowAudioInfo          string
	ShowSiteInfo           string

	CreatedTs              string
	UpdatedTs              string
}

type StreamProfile__SearchCriteria struct {
	CodeLike   string
	NameLike   string
	State      []string

	Pagination *SearchPagination
}

type StreamProfile__SearchResult struct {
	Data       []*StreamProfilePE

	Pagination *SearchPagination
}

func (db *DB) StreamProfile__EntityFullScan(entity *StreamProfilePE) ([]any) {
	return []any{
		&entity.ID,
		&entity.Code,
		&entity.Name,
		&entity.State,
		&entity.Note,

		&entity.TargetVideoFPS,
		&entity.TargetVideoWidth,
		&entity.TargetVideoHeight,
		&entity.TargetVideoCodec,
		&entity.TargetVideoCompression,
		&entity.TargetVideoBitrate,

		&entity.TargetAudioCodec,
		&entity.TargetAudioCompression,
		&entity.TargetAudioBitrate,

		&entity.ShowTimestamp,
		&entity.ShowVideoInfo,
		&entity.ShowAudioInfo,
		&entity.ShowSiteInfo,

		&entity.CreatedTs, &entity.UpdatedTs,
	}
}
