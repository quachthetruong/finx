package mock

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

type EventPublisher struct{}

func (e EventPublisher) Publish(message kafka.Message) error {
	fmt.Println(fmt.Sprintf("published message: %s, topic %s", message.Value, message.Topic))
	return nil
}
