package mock

import (
	"fmt"
)

type NotifyWebhookRepository struct{}

func (n *NotifyWebhookRepository) Send(channel, message string) error {
	fmt.Println("Send to channel: ", channel, " message: ", message)
	return nil
}
