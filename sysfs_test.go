package hotplugd

import (
	"fmt"
	"testing"
)

func TestSysfsDiscover(t *testing.T) {
	sysfs := NewSysfs()
	devices, err := sysfs.Discover()
	if err != nil {
		t.Fatal(err)
	}
	if devices == nil {
		t.Fatal("devices are nil")
	}
	for _, device := range devices {
		fmt.Println(device.String())
	}
}
