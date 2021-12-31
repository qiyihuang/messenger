package ratelimit

import (
	"net/http"
	"strconv"
	"time"
)

// Wait analysis the response header from Discord to comply with their dynamic
// rate limit.
// IMPORTANT: the function cannot prevent "webhook message/channel/min" limit.
func Wait(header http.Header) (err error) {
	// x-ratelimit-remaining contains the number of remaining quota.
	remaining := header.Get("x-ratelimit-remaining")
	// x-ratelimit-reset-after indicate the time (in sec) after which the limit
	// will be reset.
	resetAfter := header.Get("x-ratelimit-reset-after")
	// Discord sometimes respond w/o those headers.
	// No headers, no limit.
	if remaining == "" && resetAfter == "" {
		return
	}

	r, err := strconv.Atoi(remaining)
	if err != nil || r > 0 {
		return
	}

	wait, err := strconv.ParseFloat(resetAfter, 64)
	if err != nil {
		return
	}

	time.Sleep(time.Duration(wait) * time.Second)
	return
}
