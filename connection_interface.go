package chromecast

import (
	"encoding/json"

	"github.com/mcollinge/chromecast/protobuf"
)

var requestID int

type connectionInterface struct {
	connection    *connection
	sourceID      string
	destinationID string
	namespace     string
}

func newConnectionInterface(connection *connection, sourcID string, destinationID string, namespace string) *connectionInterface {
	return &connectionInterface{
		connection:    connection,
		sourceID:      sourcID,
		destinationID: destinationID,
		namespace:     namespace,
	}
}

func (c *connectionInterface) Send(payload protobuf.Payload) error {
	requestID++
	payload.SetRequestID(requestID)
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	payloadString := string(payloadJSON)
	castMessage := &protobuf.CastMessage{
		ProtocolVersion: protobuf.CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &c.sourceID,
		DestinationId:   &c.destinationID,
		Namespace:       &c.namespace,
		PayloadType:     protobuf.CastMessage_STRING.Enum(),
		PayloadUtf8:     &payloadString,
	}
	c.connection.Send <- castMessage
	return <-c.connection.Error
}

func (c *connectionInterface) SendAndWait(payload protobuf.Payload) (*Message, error) {
	var message *Message
	err := c.Send(payload)
	if err != nil {
		return message, err
	}
	return <-c.connection.Recieve, nil
}
