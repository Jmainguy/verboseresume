package main

import (
	"os"
	"strconv"
	"strings"
)

func maxUploadSize() int64 {
	if raw := strings.TrimSpace(os.Getenv("MAX_UPLOAD_BYTES")); raw != "" {
		if n, err := strconv.ParseInt(raw, 10, 64); err == nil && n > 0 {
			return n
		}
	}
	return maxUploadBytes
}
