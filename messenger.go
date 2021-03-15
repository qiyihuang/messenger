package messenger

import "net/http"

// Send sends Message to Discord webhook
func Send(url string, msg Message) (resp *http.Response, err error) {
	err = validateURL(url)
	if err != nil {
		return
	}

	err = validateMessage(msg)
	if err != nil {
		return
	}

	resp, err = makeRequest(msg, url)
	return
}
