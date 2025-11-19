package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	
	rootfs "noname001/filesystem"
)

const (
	STREAM__DIR_NAME = "stream" // TODO: from moduledefinition

	OPERATIONAL__DIR_PATH = "stream"
	LOCAL__DIR_PATH       = "stream/local"
)

// TODO: cache after root dir is set

func StreamRuntimeDir() (string) {
	return filepath.Join(rootfs.RuntimeDir, STREAM__DIR_NAME)
}

func StreamLocalDir() (string) {
	return filepath.Join(rootfs.OperationalDir, LOCAL__DIR_PATH)
}

func PrepareAll() (err error) {
	os.MkdirAll(StreamRuntimeDir(), rootfs.DEFAULT_DIRECTORY_PERMISSION)

	os.MkdirAll(StreamLocalDir(), rootfs.DEFAULT_DIRECTORY_PERMISSION)

	return
}


func StreamLocalDirPlaceholder() (string) {
	return fmt.Sprintf("[%s]/%s", rootfs.APP_ROOTDIR__PLACEHOLDER, LOCAL__DIR_PATH)
}
