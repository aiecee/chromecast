package chromecast

import (
	"fmt"

	"github.com/mcollinge/chromecast/controllers"
	"github.com/mcollinge/chromecast/net"
)

const (
	chromecastAppID     = "CC1AD845"
	chromecastAppURL    = "5CB45E5A"
	port                = 34455
	connectionNamespace = "urn:x-cast:com.google.cast.tp.connection"
	recieverNamespace   = "urn:x-cast:com.google.cast.receiver"
	mediaNamespace      = "urn:x-cast:com.google.cast.media"
	defaultSender       = "sender-0"
	defaultReciever     = "receiver-0"
	transportSender     = "Tr@n$p0rt-0"
	transportReceiver   = "Tr@n$p0rt-0"
)

var contentTypes = []string{
	"video/webm",
	"video/mp4",
}

type Device struct {
	Name                 string
	dnsEntry             net.DNSEntry
	connection           *net.Connection
	connectionController *controllers.ConnectionController
	heartbeatController  *controllers.HeartbeatController
	receiverController   *controllers.ReceiverController
	mediaController      *controllers.MediaController
}

func GetAllDevices() []Device {
	dnsEntries := net.GetAllEntries()
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
	d.connection = net.NewConnection(d.dnsEntry)
	err := d.connection.Connect()
	if err != nil {
		return err
	}
	d.connectionController = controllers.NewConnectionController(d.connection, defaultSender, defaultReciever)
	if err = d.connectionController.Start(); err != nil {
		return err
	}

	d.heartbeatController = controllers.NewHeartbeatController(d.connection, defaultSender, defaultReciever)
	if err = d.heartbeatController.Start(); err != nil {
		return err
	}

	d.receiverController = controllers.NewReceiverController(d.connection, defaultSender, defaultReciever)
	if err = d.heartbeatController.Start(); err != nil {
		return err
	}
	return nil
}

func (d *Device) Reciever() *controllers.ReceiverController {
	return d.receiverController
}

func (d *Device) Media() (*controllers.MediaController, error) {
	if d.mediaController == nil {
		transportID, err := d.startApp(chromecastAppID)
		if err != nil {
			return nil, err
		}
		conn := controllers.NewConnectionController(d.connection, defaultSender, transportID)
		if err := conn.Start(); err != nil {
			return nil, err
		}
		d.mediaController = controllers.NewMediaController(d.connection, defaultSender, transportID)
		if err := d.mediaController.Start(); err != nil {
			return nil, err
		}
	}
	return d.mediaController, nil
}

func (d *Device) Close() error {
	if d.connection != nil {
		return d.connection.Close()
	}
	return nil
}

func (d *Device) startApp(appID string) (string, error) {
	status, err := d.receiverController.GetStatus()
	if err != nil {
		return "", err
	}
	app := status.GetSessionByAppId(appID)
	if app == nil {
		status, err = d.receiverController.Launch(appID)
		if err != nil {
			return "", err
		}
		app = status.GetSessionByAppId(appID)
	}
	if app == nil {
		return "", fmt.Errorf("Could not get transportId")
	}
	return *app.TransportID, nil
}
