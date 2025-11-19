package sanitation

const (
	CODE__LEGAL_CHARS string = "[a-z, A-Z, 0-9, hypen(-), underscore(_)]"
)

func Code_StripIllegalChars(s string) (string) {
	bs := []byte(s)
	n := 0

	// TODO: strip string starting with -
	for _, b := range bs {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == '-' ||
			b == '_' {
			bs[n] = b
			n++
		}
	}

	s = string(bs[:n])

	return s
}

func Code_ContainsIllegalChar(s string) (bool, string) {
	bs := []byte(s)
	
	// TODO: disallow string starting with -
	for _, b := range bs {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == '-' ||
			b == '_' {
		} else {
			return true, string(b)
		}
	}

	return false, ""
}

// credits:
// https://stackoverflow.com/questions/54461423/efficient-way-to-remove-all-non-alphanumeric-characters-from-large-text
func _strip(bs []byte) ([]byte) {
	n := 0

	for _, b := range bs {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			bs[n] = b
			n++
		}
	}

	return bs[:n]
}
