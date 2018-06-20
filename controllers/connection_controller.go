package controllers

import (
	"github.com/mcollinge/chromecast/messages"
	"github.com/mcollinge/chromecast/net"
)

const connectionNamespace = "urn:x-cast:com.google.cast.tp.connection"

type ConnectionController struct {
	channel *net.Channel
}

func NewConnectionController(connection *net.Connection, sourceID string, destinationID string) *ConnectionController {
	return &ConnectionController{
		channel: net.NewChannel(connection, sourceID, destinationID, connectionNamespace),
	}
}

func (c *ConnectionController) Start() error {
	return c.channel.Send(&messages.ConnectPayload)
}

func (c *ConnectionController) Close() error {
	return c.channel.Send(&messages.ClosePayload)
}
