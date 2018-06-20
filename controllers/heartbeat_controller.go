package controllers

import (
	"fmt"
	"time"

	"github.com/mcollinge/chromecast/messages"
	"github.com/mcollinge/chromecast/net"
)

const (
	heartbeatNamespace = "urn:x-cast:com.google.cast.tp.heartbeat"
	maxPongs           = 3
	timerInterval      = time.Second * 5
)

type HeartbeatController struct {
	channel *net.Channel
	ticker  *time.Ticker
	pongs   int64
}

func NewHeartbeatController(connection *net.Connection, sourceID string, destinationID string) *HeartbeatController {
	controller := &HeartbeatController{
		channel: net.NewChannel(connection, sourceID, destinationID, heartbeatNamespace),
		pongs:   0,
	}
	connection.Events.Subscribe("PING", controller.ping)
	connection.Events.Subscribe("PONG", controller.pong)
	return controller
}

func (c *HeartbeatController) Start() error {
	if c.ticker != nil {
		c.Stop()
	}
	c.ticker = time.NewTicker(timerInterval)
	go c.handleTicker()
	return nil
}

func (c *HeartbeatController) Stop() error {
	if c.ticker != nil {
		c.ticker.Stop()
		c.ticker = nil
	}
	return nil
}

func (c *HeartbeatController) ping(message net.Message) {
	err := c.channel.Send(&messages.PongPayload)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *HeartbeatController) pong(message net.Message) {
	c.pongs = 0
}

func (c *HeartbeatController) handleTicker() {
heartbeat:
	for {
		select {
		case <-c.ticker.C:
			if c.pongs > maxPongs {
				break heartbeat
			}
			err := c.channel.Send(&messages.PingPayload)
			if err != nil {
				fmt.Println(err)
				break heartbeat
			}
			c.pongs++
		}
	}
}
