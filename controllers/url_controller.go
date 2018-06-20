package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/mcollinge/chromecast/net"
	"github.com/mcollinge/chromecast/protobuf"

	"github.com/mcollinge/chromecast/messages"
)

const urlNamespace = "urn:x-cast:com.url.cast"

type URLController struct {
	channel       *net.Channel
	DestinationID string
	URLSessionID  int
}

func NewURLController(connection *net.Connection, sourceID string, destinationID string) *URLController {
	controller := &URLController{
		channel:       net.NewChannel(connection, sourceID, destinationID, urlNamespace),
		DestinationID: destinationID,
	}
	connection.Events.Subscribe("URL_STATUS", controller.urlStatus)
	return controller
}

func (c *URLController) Start() error {
	message, err := c.channel.Request(&messages.URLStatusPayload)
	if err != nil {
		return err
	}
	var response messages.URLStatusResponse
	err = json.Unmarshal([]byte(*message.CastMessage.PayloadUtf8), &response)
	if err != nil {
		return err
	}
	for _, status := range response.Status {
		c.URLSessionID = status.URLSessionID
	}
	return nil
}

func (c *URLController) Stop() error {
	return nil
}

func (c *URLController) GetStatus() (*messages.URLStatusResponse, error) {
	message, err := c.channel.Request(&messages.URLStatusPayload)
	if err != nil {
		return nil, err
	}
	var response messages.URLStatusResponse
	err = json.Unmarshal([]byte(*message.CastMessage.PayloadUtf8), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *URLController) Load(url string) (*protobuf.CastMessage, error) {
	message := messages.LoadURLCommand{
		PayloadHeader: messages.LoadURLPayload,
		URL:           url,
		Type:          "loc",
	}
	response, err := c.channel.Request(&message)
	if err != nil {
		return nil, err
	}
	if response.Payload.Type == "LOAD_FAILED" {
		return nil, fmt.Errorf("Failed to load url")
	}
	return response.CastMessage, nil
}

func (c *URLController) urlStatus(message net.Message) {
	var response messages.URLStatusResponse
	err := json.Unmarshal([]byte(*message.CastMessage.PayloadUtf8), &response)
	if err != nil {
		fmt.Println(err)
	}
	for _, status := range response.Status {
		c.URLSessionID = status.URLSessionID
	}
}
