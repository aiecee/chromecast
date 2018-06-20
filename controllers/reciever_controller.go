package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/mcollinge/chromecast/net"
	"github.com/mcollinge/chromecast/protobuf"

	"github.com/mcollinge/chromecast/messages"
)

const receiverNamespace = "urn:x-cast:com.google.cast.receiver"

type ReceiverController struct {
	channel *net.Channel
	Status  *messages.ReceiverStatus
}

func NewReceiverController(connection *net.Connection, sourceID string, destinationID string) *ReceiverController {
	controller := &ReceiverController{
		channel: net.NewChannel(connection, sourceID, destinationID, receiverNamespace),
	}
	connection.Events.Subscribe("RECEIVER_STATUS", controller.receiverStatus)
	return controller
}

func (c *ReceiverController) Start() error {
	return nil
}

func (c *ReceiverController) Stop() error {
	return nil
}

func (c *ReceiverController) GetStatus() (*messages.ReceiverStatus, error) {
	message, err := c.channel.Request(&messages.RecieverStatusPayload)
	if err != nil {
		return nil, err
	}
	var response messages.StatusResponse
	err = json.Unmarshal([]byte(*message.CastMessage.PayloadUtf8), &response)
	if err != nil {
		return nil, nil
	}
	return response.Status, nil
}

func (c *ReceiverController) SetVolume(volume *messages.Volume) (*protobuf.CastMessage, error) {
	message, err := c.channel.Request(&messages.ReceiverStatus{
		PayloadHeader: messages.SetVolumePayload,
		Volume:        volume,
	})
	if err != nil {
		return nil, err
	}
	return message.CastMessage, nil
}

func (c *ReceiverController) GetVolume() (*messages.Volume, error) {
	status, err := c.GetStatus()
	if err != nil {
		return nil, err
	}
	return status.Volume, nil
}

func (c *ReceiverController) Launch(appID string) (*messages.ReceiverStatus, error) {
	message, err := c.channel.Request(&messages.LaunchRequest{
		PayloadHeader: messages.LaunchRecieverPayload,
		AppID:         appID,
	})
	if err != nil {
		return nil, err
	}
	var response messages.StatusResponse
	err = json.Unmarshal([]byte(*message.CastMessage.PayloadUtf8), &response)
	if err != nil {
		return nil, err
	}
	return response.Status, err
}

func (c *ReceiverController) Quit(appID string) (*protobuf.CastMessage, error) {
	message, err := c.channel.Request(&messages.StopRecieverPayload)
	if err != nil {
		return nil, err
	}
	return message.CastMessage, nil
}

func (c *ReceiverController) receiverStatus(message net.Message) {
	fmt.Println(message.Payload)
}
