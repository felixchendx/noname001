package filesystem

import (
	"os"
	"path/filepath"
	
	rootfs "noname001/filesystem"
)

const (
	MEDIASRV__DIR_NAME = "mediasrv" // TODO: from moduledefinition
)

func MediasrvRuntimeDir() (string) {
	return filepath.Join(rootfs.RuntimeDir, MEDIASRV__DIR_NAME)
}

func PrepareAll() (err error) {
	os.MkdirAll(MediasrvRuntimeDir(), rootfs.DEFAULT_DIRECTORY_PERMISSION)

	return
}
