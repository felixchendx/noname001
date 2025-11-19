package filesystem

import (
	"os"
	"path/filepath"
	
	rootfs "noname001/filesystem"
)

const (
	WALL__DIR_NAME = "wall" // TODO: from moduledefinition
)

func WallRuntimeDir() (string) {
	return filepath.Join(rootfs.RuntimeDir, WALL__DIR_NAME)
}

func PrepareAll() (err error) {
	os.MkdirAll(WallRuntimeDir(), rootfs.DEFAULT_DIRECTORY_PERMISSION)

	return
}
