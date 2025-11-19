package navigation

type NavLink struct {
	Code        string
	URI         string
	DisplayName string
}

func (link NavLink) AsDataMap() (map[string]any) {
	return map[string]any{
		"code"        : link.Code,
		"uri"         : link.URI,
		"display_name": link.DisplayName,

		"is_shown"    : false,
	}
}
