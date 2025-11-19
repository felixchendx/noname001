package logging

import (
	"bytes"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type CompactFormatter struct {}

func (f *CompactFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	timeFormat := "2006-01-02T15:04:05-0700" // time.RFC3339
	fmt.Fprintf(b, "[%s][%7s] %s\n", time.Now().In(loggingTz).Format(timeFormat), entry.Level, entry.Message)
	return b.Bytes(), nil
}
