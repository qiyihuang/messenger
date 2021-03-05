package ratelimit

import (
	"net/http"
	"strconv"
	"time"
)

// Wait analysisthe response header from Discord  to comply with their dynamic
// rate limit.
// ! IMPORTANT: this cannot prevent "webhook message/channel/min" limit.
func Wait(header http.Header) (err error) {
	remaining := header.Get("x-ratelimit-remaining")
	resetAfter := header.Get("x-ratelimit-reset-after")
	// Discord sometimes respond w/o those headers.
	if remaining == "" && resetAfter == "" {
		return
	}

	r, err := strconv.Atoi(remaining)
	if err != nil {
		return
	}

	if r > 0 {
		return
	}

	// This header indicate the time (in sec) after which the limit will reset.
	wait, err := strconv.ParseFloat(resetAfter, 64)
	if err != nil {
		return
	}

	time.Sleep(time.Duration(wait) * time.Second)
	return
}
