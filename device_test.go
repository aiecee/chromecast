package chromecast_test

import (
	"fmt"
	"testing"

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

		err = device.Play("http://download.blender.org/peach/bigbuckbunny_movies/BigBuckBunny_320x180.mp4")
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		fmt.Println()
		device.Close()
	}
}
