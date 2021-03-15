package messenger

import "net/http"

// Sender object that can send http requests. e.g. Request.
type Sender interface {
	send() (*http.Response, error)
}

// Send sends Message to Discord webhook
func Send(s Sender) (*http.Response, error) {
	return s.send()
}
