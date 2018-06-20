package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/mcollinge/chromecast/messages"
	"github.com/mcollinge/chromecast/net"
	"github.com/mcollinge/chromecast/protobuf"
)

const mediaNamespace = "urn:x-cast:com.google.cast.media"

type MediaController struct {
	channel        *net.Channel
	DestinationID  string
	MediaSessionID int
	Status         *messages.MediaStatus
}

func NewMediaController(connection *net.Connection, sourceID string, destinationID string) *MediaController {
	controller := &MediaController{
		channel:       net.NewChannel(connection, sourceID, destinationID, mediaNamespace),
		DestinationID: destinationID,
	}
	connection.Events.Subscribe("MEDIA_STATUS", controller.mediaStatus)
	return controller
}

func (c *MediaController) Start() error {
	_, err := c.GetStatus()
	return err
}

func (c *MediaController) GetStatus() (*messages.MediaStatusResponse, error) {
	var mediaStatus messages.MediaStatusResponse
	message, err := c.channel.Request(&messages.MediaStatusPayload)
	if err != nil {
		return &mediaStatus, err
	}
	err = json.Unmarshal([]byte(*message.CastMessage.PayloadUtf8), &mediaStatus)
	return &mediaStatus, err
}

func (c *MediaController) Play() (*protobuf.CastMessage, error) {
	var castMessage protobuf.CastMessage
	message, err := c.channel.Request(&messages.MediaCommand{
		PayloadHeader:  messages.PlayMediaPayload,
		MediaSessionID: c.MediaSessionID,
	})
	if err != nil {
		return &castMessage, err
	}
	castMessage = *message.CastMessage
	return &castMessage, nil
}

func (c *MediaController) Pause() (*protobuf.CastMessage, error) {
	var castMessage protobuf.CastMessage
	message, err := c.channel.Request(&messages.MediaCommand{
		PayloadHeader:  messages.PauseMediaPayload,
		MediaSessionID: c.MediaSessionID,
	})
	if err != nil {
		return &castMessage, err
	}
	castMessage = *message.CastMessage
	return &castMessage, nil
}

func (c *MediaController) Stop() error {
	if c.MediaSessionID == 0 {
		return nil
	}
	return c.channel.Send(&messages.MediaCommand{
		PayloadHeader:  messages.StopMediaPayload,
		MediaSessionID: c.MediaSessionID,
	})
}

func (c *MediaController) Load(media messages.MediaItem, currentTime int, autoplay bool, customData interface{}) (*protobuf.CastMessage, error) {
	message := &messages.LoadMediaCommand{
		PayloadHeader: messages.LoadMediaPayload,
		Media:         media,
		CurrentTime:   currentTime,
		Autoplay:      autoplay,
		CustomData:    customData,
	}
	var castMessage protobuf.CastMessage
	response, err := c.channel.Request(message)
	if err != nil {
		return &castMessage, err
	}
	if response.Payload.Type == "LOAD_FAILED" {
		return &castMessage, fmt.Errorf("Load media Failed")
	}
	return response.CastMessage, nil
}

func (c *MediaController) mediaStatus(message net.Message) {
	var response messages.MediaStatusResponse
	err := json.Unmarshal([]byte(*message.CastMessage.PayloadUtf8), &response)
	if err != nil {
		fmt.Println(err)
	}
	for _, status := range response.Status {
		c.Status = status
		c.MediaSessionID = status.MediaSessionID
	}
}
