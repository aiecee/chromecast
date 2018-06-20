package net

import (
	"encoding/json"

	"github.com/mcollinge/chromecast/messages"
	"github.com/mcollinge/chromecast/protobuf"
)

var requestID int

type Channel struct {
	connection    *Connection
	sourceID      string
	destinationID string
	namespace     string
}

func NewChannel(connection *Connection, sourceID string, destinationID string, namespace string) *Channel {
	return &Channel{
		connection:    connection,
		sourceID:      sourceID,
		destinationID: destinationID,
		namespace:     namespace,
	}
}

func (c *Channel) Send(payload messages.Payload) error {
	castMessage, err := c.buildMessage(payload)
	if err != nil {
		return err
	}
	_, err = c.connection.Send(requestID, castMessage)
	if err != nil {
		return err
	}
	return nil
}

func (c *Channel) Request(payload messages.Payload) (Message, error) {
	var message Message
	castMessage, err := c.buildMessage(payload)
	if err != nil {
		return message, err
	}
	responseChannel, err := c.connection.Send(requestID, castMessage)
	if err != nil {
		return message, err
	}
	return <-responseChannel, nil
}

func (c *Channel) buildMessage(payload messages.Payload) (*protobuf.CastMessage, error) {
	requestID++
	payload.SetRequestID(requestID)
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	payloadString := string(payloadJSON)
	return &protobuf.CastMessage{
		ProtocolVersion: protobuf.CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &c.sourceID,
		DestinationId:   &c.destinationID,
		Namespace:       &c.namespace,
		PayloadType:     protobuf.CastMessage_STRING.Enum(),
		PayloadUtf8:     &payloadString,
	}, nil
}
