package live

import (
	streamTyping "noname001/app/base/typing/stream"

	mediasrvIntface "noname001/app/module/common/mediasrv/intface"
)

func (lc *LiveCache) refreshStreamService(node *t_node, cascadingRefresh bool) {
	streamServiceInfo, err := lc.fetchStreamServiceInfo(node.id)
	if err != nil {
		// TODO: activity log
		return
	}

	_ = streamServiceInfo

	node.streamService.opFlag = true
	node.streamService.opStatus = "ok"

	if cascadingRefresh {
		lc.refreshStreams(node)
	}
}

func (lc *LiveCache) refreshStreams(node *t_node) {
	streamSnapshotList, err := lc.fetchStreamSnapshotList(node.id)
	if err != nil {
		// activity log
		return
	}

	for _, streamSnapshot := range streamSnapshotList {
		lc._streamRefreshRoutine(node, streamSnapshot)
	}
}

func (lc *LiveCache) refreshStream(node *t_node, streamCode string) {
	streamSnapshot, err := lc.fetchStreamSnapshot(node.id, streamCode)
	if err != nil {
		// activity log
		return
	}

	lc._streamRefreshRoutine(node, streamSnapshot)
}


func (lc *LiveCache) defunctStreamService(node *t_node, reason string) {
	node.streamService.opFlag = false
	node.streamService.opStatus = reason

	lc.defunctStreams(node, reason)
}

func (lc *LiveCache) defunctStreams(node *t_node, reason string) {
	for _, stream := range node.streamService.streams {
		lc.defunctStream(node, stream, reason)
	}
}

func (lc *LiveCache) defunctStream(node *t_node, stream *t_stream, defunctReason string) {
	node.streamService.markStreamAsDefunct(stream, defunctReason)

	mediasrvIntface.Provider().DeletePathConfiguration(stream.relayPathName)
}


func (lc *LiveCache) _streamRefreshRoutine(node *t_node, _streamSnapshot *streamTyping.StreamSnapshot) {
	var (
		seenStream, alreadySeen = node.streamService.streams[_streamSnapshot.Code]

		localMediasrvProvider = mediasrvIntface.Provider()
		hasStreamChanges bool = false

		prevLiveState    string = ""
		prevSourceStream string = ""
	)

	if alreadySeen {
		prevLiveState    = string(seenStream.streamSnapshot.Live.State)
		prevSourceStream = seenStream.sourceStream

		node.streamService.updateStreamData(seenStream, _streamSnapshot)

	} else {
		seenStream = node.streamService.addNewStream(_streamSnapshot)
	}

	lc.interpretStreamStateToStreamStatus(seenStream)
	seenStream.relayPathName = localMediasrvProvider.RelayedStreamURI(node.id, seenStream.code)

	if node.mediaServer.ipToUse == "" {
		seenStream.sourceStream = ""

	} else {
		seenStream.sourceStream = localMediasrvProvider.SourceStreamURL(
			node.mediaServer.authnPair,
			node.mediaServer.ipToUse,
			node.mediaServer.ports["rtsp"],
			seenStream.code,
		)
	}

	if prevLiveState != string(seenStream.streamSnapshot.Live.State) { hasStreamChanges = true }
	if prevSourceStream != seenStream.sourceStream { hasStreamChanges = true }

	if hasStreamChanges {
		if seenStream.sourceStream == "" {
			localMediasrvProvider.DeletePathConfiguration(seenStream.relayPathName)
		} else {
			localMediasrvProvider.ReplacePathConfiguration(seenStream.relayPathName, seenStream.sourceStream, true)
		}
	}
}
