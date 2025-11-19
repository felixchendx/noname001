package embedding

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"noname001/filesystem"
)

var (
	//go:embed embedroot/*
	EmbedFS embed.FS
)

type UnpackRoutine string

const (
	UNPACK_ROUTINE_FILE UnpackRoutine = "file"
	UNPACK_ROUTINE_TAR  UnpackRoutine = "tar"
)

type UnpackParams struct {
	EmbeddedFilepath EmbeddedFile
	UnpackRoutine    UnpackRoutine

	UnpackDirectory  string
	UnpackFilename   string
}
func Unpack(params *UnpackParams) (err error) {
	switch params.UnpackRoutine {
	case UNPACK_ROUTINE_FILE: err = unpackFile(params)
	case UNPACK_ROUTINE_TAR:  err = unpackTar(params)
	}

	if err != nil { return }

	return
}

func unpackFile(params *UnpackParams) (error) {
	fileContent, readErr := EmbedFS.ReadFile(string(params.EmbeddedFilepath))
	if readErr != nil { return readErr }

	mkdirErr := os.MkdirAll(params.UnpackDirectory, filesystem.DEFAULT_DIRECTORY_PERMISSION)
	if mkdirErr != nil { return mkdirErr }

	writeErr := os.WriteFile(
		filepath.Join(params.UnpackDirectory, params.UnpackFilename),
		fileContent,
		filesystem.DEFAULT_FILE_PERMISSION,
	)
	if writeErr != nil { return writeErr }

	return nil
}


// this routine is based on file 'mediamtx_v1.8.3_linux_amd64.tar.gz'
// with structure:
// - LICENSE
// - mediamtx
// - mediamtx.yml
func unpackTar(params *UnpackParams) (error) {
	fileContent, readErr := EmbedFS.ReadFile(string(params.EmbeddedFilepath))
	if readErr != nil { return readErr }

	mkdirErr := os.MkdirAll(params.UnpackDirectory, filesystem.DEFAULT_DIRECTORY_PERMISSION)
	if mkdirErr != nil { return mkdirErr }

	tmpUnpackDir := filesystem.PrepareTemporaryDir("mediamtx")
	defer filesystem.CleanupTemporaryDir(tmpUnpackDir)

	tmpArchivePath := filepath.Join(tmpUnpackDir, params.UnpackFilename)
	writeErr := os.WriteFile(tmpArchivePath, fileContent, filesystem.DEFAULT_FILE_PERMISSION)
	if writeErr != nil { return writeErr }

	// TODO: use pure go implementation ?
	cmdBash := exec.Command("/usr/bin/env", "/bin/bash")
	cmdBashIn, _ := cmdBash.StdinPipe()
	cmdBash.Stdout = os.Stdout
	cmdBash.Stderr = os.Stderr

	cmdBashErr := cmdBash.Start()
	if cmdBashErr != nil { return cmdBashErr }

	tarCommand       := fmt.Sprintf("/usr/bin/tar -xzf %s -C %s", tmpArchivePath, tmpUnpackDir)
	cpBinCommand     := fmt.Sprintf("/usr/bin/mv %s %s", filepath.Join(tmpUnpackDir, "mediamtx"), filepath.Join(params.UnpackDirectory, params.UnpackFilename)) 
	cpLicenseCommand := fmt.Sprintf("/usr/bin/mv %s %s", filepath.Join(tmpUnpackDir, "LICENSE"), filepath.Join(params.UnpackDirectory, "LICENSE"))
	cmdBashIn.Write([]byte(tarCommand + "\n"))
	cmdBashIn.Write([]byte(cpBinCommand + "\n"))
	cmdBashIn.Write([]byte(cpLicenseCommand + "\n"))
	cmdBashIn.Close()

	cmdBashWaitErr := cmdBash.Wait()
	if cmdBashWaitErr != nil { return cmdBashWaitErr }

	return nil
}
