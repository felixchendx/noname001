package version

import (
	"fmt"
)

func FullVersion() (string) {
	return fmt.Sprintf("v%s_build-%s", VERSION, COMPILED_AT)
}
