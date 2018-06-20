package net

import (
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/asaskevich/EventBus"

	"github.com/golang/protobuf/proto"

	"github.com/mcollinge/chromecast/messages"
	"github.com/mcollinge/chromecast/protobuf"
)

type Message struct {
	CastMessage *protobuf.CastMessage
	Payload     *messages.PayloadHeader
}

type Connection struct {
	dnsEntry   DNSEntry
	connection *tls.Conn
	Responses  map[int]chan Message
	Events     EventBus.Bus
}

func NewConnection(dnsEntry DNSEntry) *Connection {
	return &Connection{
		dnsEntry:  dnsEntry,
		Responses: make(map[int]chan Message),
		Events:    EventBus.New(),
	}
}

func (c *Connection) Connect() error {
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
	go c.recieve()
	return nil
}

func (c *Connection) Close() error {
	if c.connection != nil {
		return c.connection.Close()
	}
	return nil
}

func (c *Connection) Send(requestID int, message *protobuf.CastMessage) (chan Message, error) {
	proto.SetDefaults(message)
	data, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	responseChannel := make(chan Message, 1)
	c.Responses[requestID] = responseChannel
	err = binary.Write(c.connection, binary.BigEndian, uint32(len(data)))
	if err != nil {
		return nil, err
	}
	written, err := c.connection.Write(data)
	fmt.Println(written)
	return responseChannel, err
}

func (c *Connection) recieve() {
	for {
		var length uint32
		err := binary.Read(c.connection, binary.BigEndian, &length)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if length == 0 {
			continue
		}
		packet := make([]byte, length)
		i, err := io.ReadFull(c.connection, packet)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if i != int(length) {
			fmt.Println(err)
			continue
		}

		message := &protobuf.CastMessage{}
		err = proto.Unmarshal(packet, message)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var headers messages.PayloadHeader
		fmt.Println(*message.PayloadUtf8)
		err = json.Unmarshal([]byte(*message.PayloadUtf8), &headers)
		if err != nil {
			fmt.Println(err)
			continue
		}
		messageWrapper := Message{
			CastMessage: message,
			Payload:     &headers,
		}
		if val, ok := c.Responses[headers.RequestID]; ok {
			val <- messageWrapper
			close(val)
			delete(c.Responses, headers.RequestID)
		}
		c.Events.Publish(headers.Type, messageWrapper)
	}
}
