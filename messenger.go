package messenger

import (
	"github.com/qiyihuang/messenger/pkg/message"
	"github.com/qiyihuang/messenger/pkg/request"
)

// Send sends Message to Discord webhook
func Send(url string, msg message.Message) error {
	err := message.Validate(msg)
	if err != nil {
		return err
	}

	_, err = request.Send(msg, url)
	if err != nil {
		return err
	}

	return nil
}
