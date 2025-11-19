package sqlite

type WallLayout__SearchCriteria struct {
	CodeLike   *string
	NameLike   *string
	State      []string

	Pagination *SearchPagination
}
type WallLayout__SearchResult struct {
	Data       []*WallLayoutPE

	Pagination *SearchPagination
}

type WallLayoutPE struct {
	ID              string
	Code            string
	Name            string
	State           string
	Note            string

	LayoutFormation string
	LayoutItemCount int

	DefinedBy       string
	CreatedTs       string
	UpdatedTs       string
}

func (db *DB) WallLayout__EntityFullScan(entity *WallLayoutPE) ([]any) {
	return []any{
		&entity.ID,
		&entity.Code,
		&entity.Name,
		&entity.State,
		&entity.Note,

		&entity.LayoutFormation,
		&entity.LayoutItemCount,

		&entity.DefinedBy,
		&entity.CreatedTs, &entity.UpdatedTs,
	}
}
