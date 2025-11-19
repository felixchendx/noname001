package sqlite

type StreamGroupPE struct {
	ID              string
	Code            string
	Name            string
	State           string
	Note            string

	StreamProfileID string

	CreatedTs       string
	UpdatedTs       string
}

func (db *DB) StreamGroup__EntityFullScan(entity *StreamGroupPE) ([]any) {
	return []any{
		&entity.ID,
		&entity.Code,
		&entity.Name,
		&entity.State,
		&entity.Note,

		&entity.StreamProfileID,

		&entity.CreatedTs, &entity.UpdatedTs,
	}
}

func (entity *StreamGroupPE) fullScan() ([]any) {
	return []any{
		&entity.ID,
		&entity.Code,
		&entity.Name,
		&entity.State,
		&entity.Note,

		&entity.StreamProfileID,

		&entity.CreatedTs, &entity.UpdatedTs,
	}
}
