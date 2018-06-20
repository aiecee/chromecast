package chromecast_test

import (
	"fmt"
	"testing"

	"github.com/mcollinge/chromecast/messages"

	"github.com/mcollinge/chromecast"
)

func TestGettingDevices(t *testing.T) {
	devices := chromecast.GetAllDevices()
	for i, device := range devices {
		err := device.Start()
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		fmt.Printf("%v: %s", i, device.Name)

		media, err := device.Media()
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		mediaItem := messages.MediaItem{
			ContentID:   "http://mirrors.standaloneinstaller.com/video-sample/jellyfish-25-mbps-hd-hevc.mp4",
			StreamType:  "BUFFERED",
			ContentType: "video/mp4",
		}
		_, err = media.Load(mediaItem, 0, true, map[string]interface{}{})
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
	}
}
