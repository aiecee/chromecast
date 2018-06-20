package chromecast

import (
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/mcollinge/chromecast/protobuf"
)

type Message struct {
	CastMessage *protobuf.CastMessage
	Payload     *protobuf.PayloadHeader
}

type connection struct {
	dnsEntry   dnsEntry
	connection *tls.Conn
	Send       chan *protobuf.CastMessage
	Recieve    chan *Message
	Error      chan error
}

func newConnection(dnsEntry dnsEntry) *connection {
	return &connection{
		dnsEntry: dnsEntry,
		Send:     make(chan *protobuf.CastMessage),
		Recieve:  make(chan *Message),
		Error:    make(chan error),
	}
}

func (c *connection) Connect() error {
	dialer := &net.Dialer{
		Timeout:   time.Second * 30,
		KeepAlive: time.Second * 30,
	}
	var err error
	c.connection, err = tls.DialWithDialer(
		dialer,
		"tcp",
		fmt.Sprintf("%s:%d", c.dnsEntry.AddressV4, c.dnsEntry.Port),
		&tls.Config{
			InsecureSkipVerify: true,
		},
	)
	if err != nil {
		return fmt.Errorf("Failed to connect to device: %s", err)
	}
	go c.send()
	go c.recieve()
	return nil
}

func (c *connection) send() {
	for {
		message := <-c.Send
		proto.SetDefaults(message)
		data, err := proto.Marshal(message)
		if err != nil {
			c.Error <- err
		}
		err = binary.Write(c.connection, binary.BigEndian, uint32(len(data)))
		if err != nil {
			c.Error <- err
		}
		_, err = c.connection.Write(data)
		if err != nil {
			c.Error <- err
		}
		c.Error <- nil
	}
}

func (c *connection) recieve() {
	for {
		var length uint32
		err := binary.Read(c.connection, binary.BigEndian, &length)
		if err != nil {
			c.Error <- err
			continue
		}
		if length == 0 {
			continue
		}
		packet := make([]byte, length)
		i, err := io.ReadFull(c.connection, packet)
		if err != nil {
			c.Error <- err
			continue
		}
		if i != int(length) {
			c.Error <- err
			continue
		}

		message := &protobuf.CastMessage{}
		err = proto.Unmarshal(packet, message)
		if err != nil {
			c.Error <- err
			continue
		}

		var headers protobuf.PayloadHeader
		err = json.Unmarshal([]byte(*message.PayloadUtf8), &headers)
		if err != nil {
			c.Error <- err
			continue
		}
		switch headers.Type {
		case "PING":
			payloadJSON, err := json.Marshal(protobuf.Pong)
			if err != nil {
				c.Error <- err
				continue
			}
			payload := string(payloadJSON)
			pongMessage := &protobuf.CastMessage{
				ProtocolVersion: protobuf.CastMessage_CASTV2_1_0.Enum(),
				SourceId:        message.SourceId,
				DestinationId:   message.DestinationId,
				Namespace:       message.Namespace,
				PayloadType:     protobuf.CastMessage_STRING.Enum(),
				PayloadUtf8:     &payload,
			}
			c.Send <- pongMessage
		default:
			messageWrapper := &Message{
				CastMessage: message,
				Payload:     &headers,
			}
			c.Recieve <- messageWrapper
		}
	}
}
