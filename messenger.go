package messenger

import (
	"github.com/qiyihuang/messenger/pkg/message"
	"github.com/qiyihuang/messenger/pkg/request"
)

const version = "0.1.0"

// Send sends Message to Discord webhook
func Send(url string, msg message.Message) (err error) {
	err = message.Validate(msg)
	if err != nil {
		return
	}

	_, err = request.Send(msg, url)
	if err != nil {
		return
	}

	return
}
