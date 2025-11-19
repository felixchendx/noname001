package sqlite

type SysUser__SearchCriteria struct {
	UsernameLike string
	RoleSimple   []string

	Pagination *SearchPagination
}

type SysUser__SearchResult struct {
	Data       []*SysUserPE

	Pagination *SearchPagination
}

// TODO: add state, especially to suspend a user
type SysUserPE struct {
	ID            string
	Username      string // find way to secure / hide this info
	Password      string
	RoleSimple    string // this also
}
func (entity *SysUserPE) fullScan() ([]any) {
	return []any{
		&entity.ID,
		&entity.Username,
		&entity.Password,
		&entity.RoleSimple,
	}
}
