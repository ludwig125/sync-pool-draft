package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"time"
)

// timeNow is a fake version of time.Now for tests.
func timeNow() time.Time {
	return time.Unix(1136214245, 0) // 2006-01-02T15:04:05Z
}

// Log関数のPoolを使わない版
func LogWithoutPool(w io.Writer, key, val string) {
	b := &bytes.Buffer{}
	// Replace this with time.Now() in a real logger.
	b.WriteString(timeNow().UTC().Format(time.RFC3339))
	b.WriteByte(' ')
	b.WriteString(key)
	b.WriteByte('=')
	b.WriteString(val)
	if _, err := w.Write(b.Bytes()); err != nil {
		log.Fatal(err)
	}
}

func main() {
	LogWithoutPool(os.Stdout, "path", "/search?q=flowers")
}
