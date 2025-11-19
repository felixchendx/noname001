package mediamtxserver

import (
	"path/filepath"

	"noname001/app/module/common/mediasrv/filesystem"
	
	mediamtxEmbedding "noname001/thirdparty/mediamtx/embedding"
)

func (srv *MediaMTXServer) UnpackMediamtx() (error) {
	unpackList := []*mediamtxEmbedding.UnpackParams{
		&mediamtxEmbedding.UnpackParams{
			mediamtxEmbedding.V1_8_3__ARCHIVE_FILE, mediamtxEmbedding.UNPACK_ROUTINE_TAR,
			filesystem.MediasrvRuntimeDir(), "mediamtx",
		},
		&mediamtxEmbedding.UnpackParams{
			mediamtxEmbedding.MOD_MEDIASRV__V1_8_3__CONFIG_FILE, mediamtxEmbedding.UNPACK_ROUTINE_FILE,
			filesystem.MediasrvRuntimeDir(), "mediamtx.yml",
		},
	}

	for _, item := range unpackList {
		err := mediamtxEmbedding.Unpack(item)
		if err != nil {
			return err
		}
	}

	return nil
}
func (srv *MediaMTXServer) unpackedBinPath() (string) {
	return filepath.Join(filesystem.MediasrvRuntimeDir(), "mediamtx")
}
func (srv *MediaMTXServer) unpackedConfigPath() (string) {
	return filepath.Join(filesystem.MediasrvRuntimeDir(), "mediamtx.yml")
}
