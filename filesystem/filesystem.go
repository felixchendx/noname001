package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	DEFAULT_DIRECTORY_PERMISSION  = 0750
	DEFAULT_FILE_PERMISSION       = 0664
	DEFAULT_EXECUTABLE_PERMISSION = 0775

	APP_ROOTDIR__PLACEHOLDER = "app_rootdir"

	APPDATA_DIR_NAME     = ".appdata"
	LOG_DIR_NAME         = "log"
	OPERATIONAL_DIR_NAME = "operational"
)

var (
	RootDir        string

	AppdataDir     string
	DBDir          string
	TmpDir         string
	RuntimeDir     string

	WorkspaceDir   string

	LogDir         string
	OperationalDir string
)

func SetRootDir(rootDir string) {
	RootDir      = rootDir

	AppdataDir   = filepath.Join(RootDir, APPDATA_DIR_NAME)
	DBDir        = filepath.Join(RootDir, APPDATA_DIR_NAME, "db")
	TmpDir       = filepath.Join(RootDir, APPDATA_DIR_NAME, "tmp")
	RuntimeDir   = filepath.Join(RootDir, APPDATA_DIR_NAME, "runtime")

	WorkspaceDir = filepath.Join(RootDir, APPDATA_DIR_NAME, "workspace")

	LogDir         = filepath.Join(RootDir, LOG_DIR_NAME)
	OperationalDir = filepath.Join(RootDir, OPERATIONAL_DIR_NAME)
}

func PrepareAll() (err error) {
	err = os.MkdirAll(AppdataDir, DEFAULT_DIRECTORY_PERMISSION)
	if err != nil {
		return fmt.Errorf("filesystem: Unable to create directories. (%v)\n", err)
	}

	os.MkdirAll(DBDir, DEFAULT_DIRECTORY_PERMISSION)
	os.RemoveAll(TmpDir);os.MkdirAll(TmpDir, DEFAULT_DIRECTORY_PERMISSION)
	os.RemoveAll(RuntimeDir);os.MkdirAll(RuntimeDir, DEFAULT_DIRECTORY_PERMISSION)

	os.MkdirAll(WorkspaceDir, DEFAULT_DIRECTORY_PERMISSION)

	os.MkdirAll(LogDir, DEFAULT_DIRECTORY_PERMISSION)
	os.MkdirAll(OperationalDir, DEFAULT_DIRECTORY_PERMISSION)

	return
}

func PrepareTemporaryDir(dirName string) (string) {
	dir, _ := os.MkdirTemp(TmpDir, dirName + "_*")
	return dir
}
func CleanupTemporaryDir(dir string) {
	os.RemoveAll(dir)
}
