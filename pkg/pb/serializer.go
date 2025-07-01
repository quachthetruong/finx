package pb

import (
	"fmt"

	"gitlab.com/enCapital/models"
	"google.golang.org/protobuf/proto"
)

func MarshalEncapMessage(payload proto.Message, messageType models.EncapMessage_MessageType) (marshaled []byte, err error) {
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("MarshalEncapMessage %w", err)
	}
	encapMessage := &models.EncapMessage{
		Type:    messageType,
		Payload: payloadBytes,
	}
	encapMessageBytes, err := proto.Marshal(encapMessage)
	if err != nil {
		return nil, fmt.Errorf("MarshalEncapMessage %w", err)
	}
	return encapMessageBytes, nil
}
