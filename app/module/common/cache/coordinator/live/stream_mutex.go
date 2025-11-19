package live

import (
	"slices"
)

func (streamService *t_streamService) registerStream(stream *t_stream) {
	streamService.streamsMutex.Lock()

	streamService.streams[stream.code] = stream
	streamService._generateSortedStreams()

	streamService.streamsMutex.Unlock()
}

func (streamService *t_streamService) deregisterStream(stream *t_stream) {
	streamService.streamsMutex.Lock()

	delete(streamService.streams, stream.code)
	streamService._generateSortedStreams()

	streamService.streamsMutex.Unlock()
}

func (streamService *t_streamService) _generateSortedStreams() {
	streamCodes := make([]string, len(streamService.streams))
	i := 0
	for _streamCode, _ := range streamService.streams {
		streamCodes[i] = _streamCode
		i++
	}

	slices.SortFunc(streamCodes, _caseInsensitiveSort)

	sortedStreams := make([]*t_stream, len(streamCodes))
	for _idx, _streamCode := range streamCodes {
		sortedStreams[_idx] = streamService.streams[_streamCode]
	}

	streamService.sortedStreams = sortedStreams
}
