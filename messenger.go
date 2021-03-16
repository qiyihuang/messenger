package messenger

import "net/http"

// Sender object that can send http requests. e.g. Request.
type Sender interface {
	send(httpPoster) (*http.Response, error)
}

// Send sends Message to Discord webhook
func Send(s Sender) (*http.Response, error) {
	// In the future may let user customise client, just pass http.DefaultClient at the moment.
	return s.send(http.DefaultClient)
}
