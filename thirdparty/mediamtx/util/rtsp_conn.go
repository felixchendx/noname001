package util

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"time"
	"net"
)

func TestRTSPConnWithRetry(ip, port string) (bool, error) {
	attempt, maxRetry, retryBackoff := 0, 3, (1 * time.Second)
	for {
		attempt++

		reachable, err := TestRTSPConn(ip, port)
		if reachable {
			return true, nil
		} else if attempt >= maxRetry {
			return false, err // whatever last error it is
		} else {
			time.Sleep(retryBackoff)
		}
	}

	return false, nil // never reached this line
}

func TestRTSPConn(ip, port string) (bool, error) {
	host := fmt.Sprintf("%s:%s", ip, port)

	var sb strings.Builder
	fmt.Fprintf(&sb, "OPTIONS * RTSP/1.0\r\n")
	// fmt.Fprintf(&sb, "Authorization: Basic cmVsYXkwMDE6cmVhbGx5MDAx\r\n")
	fmt.Fprintf(&sb, "CSeq: 1\r\n")
	fmt.Fprintf(&sb, "User-Agent: noname-pingu\r\n")
	fmt.Fprintf(&sb, "\r\n")
	reqPayload := sb.String()

	// var sb strings.Builder
	// fmt.Fprintf(&sb, "DESCRIBE rtsp://%s@%s/%s RTSP/1.0\r\n", relayAuth, host, "pingu")
	// fmt.Fprintf(&sb, "Accept: application/sdp\r\n")
	// fmt.Fprintf(&sb, "CSeq: 2\r\n")
	// fmt.Fprintf(&sb, "\r\n")
	// reqPayload := sb.String()

	successPayloadLine01    := []byte("RTSP/1.0 200 OK\r\n")
	badRequestPayloadLine01 := []byte("RTSP/1.0 400 Bad Request\r\n")
	notFoundPayloadLine01   := []byte("RTSP/1.0 404 Not Found\r\n")

	conn, err := net.DialTimeout("tcp", host, 1 * time.Second)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(1 * time.Second))
	fmt.Fprintf(conn, reqPayload)

	respReader := bufio.NewReader(conn)
	respLine01, err := respReader.ReadSlice('\n')
	if err != nil {
		return false, err
	}

	switch {
	case bytes.Equal(respLine01, successPayloadLine01)   : // ok
	case bytes.Equal(respLine01, badRequestPayloadLine01): // also ok
	case bytes.Equal(respLine01, notFoundPayloadLine01)  : // olso ok
	default:
		return false, fmt.Errorf("_testRTSPConn respLine01 - %s", respLine01)
	}

	return true, nil
}
