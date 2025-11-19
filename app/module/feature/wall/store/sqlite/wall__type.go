package sqlite

// ============================= VVV header VVV ============================= //
type WallPE struct {
	ID            string
	Code          string
	Name          string
	State         string
	Note          string

	WallLayoutID  string

	CreatedTs     string
	UpdatedTs     string
}
func (entity *WallPE) fullScan() ([]any) {
	return []any{
		&entity.ID,
		&entity.Code,
		&entity.Name,
		&entity.State,
		&entity.Note,

		&entity.WallLayoutID,

		&entity.CreatedTs, &entity.UpdatedTs,
	}
}
// ============================= ^^^ header ^^^ ============================= //

// ============================== VVV item VVV ============================== //
type WallItemPE struct {
	ID           string
	WallID       string
	Index        int

	SourceNodeID string
	StreamCode   string

	CreatedTs    string
	UpdatedTs    string
}
func (entity *WallItemPE) fullScan() ([]any) {
	return []any{
		&entity.ID,
		&entity.WallID,
		&entity.Index,

		&entity.SourceNodeID,
		&entity.StreamCode,

		&entity.CreatedTs, &entity.UpdatedTs,
	}
}
// ============================== ^^^ item ^^^ ============================== //
