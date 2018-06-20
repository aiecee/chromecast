package chromecast

import (
	"encoding/json"
	"fmt"

	"github.com/mcollinge/chromecast/protobuf"
)

const (
	chromecastAppID     = "CC1AD845"
	port                = 34455
	connectionNamespace = "urn:x-cast:com.google.cast.tp.connection"
	recieverNamespace   = "urn:x-cast:com.google.cast.receiver"
	mediaNamespace      = "urn:x-cast:com.google.cast.media"
	defaultSender       = "sender-0"
	defaultReciever     = "receiver-0"
)

var contentTypes = []string{
	"video/webm",
	"video/mp4",
}

type Device struct {
	Name       string
	dnsEntry   dnsEntry
	connection *connection

	defaultConnection *connectionInterface
	defaultReciever   *connectionInterface
	mediaConnection   *connectionInterface
	mediaReciever     *connectionInterface

	Application protobuf.Application
	Media       protobuf.Media
	Volume      protobuf.VolumeConfig
}

func GetAllDevices() []Device {
	dnsEntries := getAllEntries()
	devices := make([]Device, len(dnsEntries))
	for i, entry := range dnsEntries {
		devices[i] = Device{
			Name:     entry.DeviceName,
			dnsEntry: entry,
		}
	}
	return devices
}

func (d *Device) Start() error {
	d.connection = newConnection(d.dnsEntry)
	err := d.connection.Connect()
	if err != nil {
		return err
	}
	d.defaultConnection = newConnectionInterface(d.connection, defaultSender, defaultReciever, connectionNamespace)
	d.defaultReciever = newConnectionInterface(d.connection, defaultSender, defaultReciever, recieverNamespace)

	if err = d.defaultConnection.Send(&protobuf.Connect); err != nil {
		return err
	}
	return d.Update()
}

func (d *Device) Update() error {
	reciverStatus, err := d.recieverStatus()
	if err != nil {
		return err
	}
	for _, application := range reciverStatus.Status.Applications {
		d.Application = application
	}
	d.Volume = reciverStatus.Status.Volume
	if reciverStatus.Status.IsStandBy {
		return nil
	}
	d.mediaConnection = newConnectionInterface(d.connection, defaultSender, d.Application.TransportID, connectionNamespace)
	d.mediaReciever = newConnectionInterface(d.connection, defaultSender, d.Application.TransportID, mediaNamespace)
	return d.updateMediaStatus()
}

func (d *Device) Close() {
	if d.mediaConnection != nil {
		d.mediaConnection.Send(&protobuf.Close)
	}
	d.defaultConnection.Send(&protobuf.Close)
}

func (d *Device) Play(url string) error {
	if d.Application.AppID != chromecastAppID {
		_, err := d.defaultReciever.SendAndWait(&protobuf.LaunchRequest{
			PayloadHeader: protobuf.Launch,
			AppID:         chromecastAppID,
		})
		if err != nil {
			return err
		}
		if err = d.Update(); err != nil {
			return err
		}
	}
	return d.mediaReciever.Send(&protobuf.LoadMediaCommand{
		PayloadHeader: protobuf.Load,
		CurrentTime:   0,
		Autoplay:      true,
		Media: protobuf.MediaItem{
			ContentID:   url,
			StreamType:  "BUFFERED",
			ContentType: "video/mp4",
		},
	})
}

func (d *Device) updateMediaStatus() error {
	d.mediaConnection.Send(&protobuf.Connect)
	mediaStatus, err := d.mediaStatus()
	if err != nil {
		return err
	}
	for _, media := range mediaStatus.Status {
		d.Media = media
		d.Volume = media.Volume
	}
	return nil
}

func (d *Device) recieverStatus() (*protobuf.ReceiverStatusResponse, error) {
	message, err := d.defaultReciever.SendAndWait(&protobuf.GetStatus)
	if err != nil {
		return nil, err
	}
	var response protobuf.ReceiverStatusResponse

	if err = json.Unmarshal([]byte(*message.CastMessage.PayloadUtf8), &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (d *Device) mediaStatus() (*protobuf.MediaStatusResponse, error) {
	message, err := d.mediaReciever.SendAndWait(&protobuf.GetStatus)
	if err != nil {
		return nil, err
	}
	res := *message.CastMessage.PayloadUtf8
	fmt.Println(res)
	var response protobuf.MediaStatusResponse
	if err = json.Unmarshal([]byte(*message.CastMessage.PayloadUtf8), &response); err != nil {
		return nil, err
	}
	return &response, nil
}
