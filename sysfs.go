package hotplugd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	sysfsUsbDevices = "/sys/bus/usb/devices"
	// *bus*-*port* or *bus*-*port.*port*, ignore hubs and configurations/interfaces
	usbDeviceLinkRegex = "^[0-9]+-[0-9]+\\.?[0-9]*$"
)

func getDeviceProperty(node string, key string) (string, error) {
	raw, err := os.ReadFile(path.Join(node, key))
	if err != nil {
		return "", err
	}
	return strings.Trim(string(raw), "\n "), nil
}

type USBDevice struct {
	node         string
	idVendor     string
	idProduct    string
	manufacturer string
	product      string
}

func NewUSBDevice(node string) USBDevice {
	return USBDevice{node: node}
}

func (u *USBDevice) populate() error {
	idVendor, err := getDeviceProperty(u.node, "idVendor")
	if err != nil {
		return err
	}
	u.idVendor = idVendor
	idProduct, err := getDeviceProperty(u.node, "idProduct")
	if err != nil {
		return err
	}
	u.idProduct = idProduct
	manufacturer, _ := getDeviceProperty(u.node, "manufacturer")
	u.manufacturer = manufacturer
	product, _ := getDeviceProperty(u.node, "product")
	u.product = product
	return nil
}

func (u *USBDevice) String() string {
	desc := u.manufacturer
	if len(u.product) > 0 {
		desc += fmt.Sprintf(" %s", u.product)
	}
	desc = strings.Trim(desc, " ")
	if len(desc) == 0 {
		desc = "Unknown"
	}
	return fmt.Sprintf("%s:%s %s", u.idVendor, u.idProduct, desc)
}

type Sysfs struct {
	usbDevRe *regexp.Regexp
}

func NewSysfs() Sysfs {
	s := Sysfs{}
	s.usbDevRe, _ = regexp.Compile(usbDeviceLinkRegex)
	return s
}

func (s Sysfs) Discover() ([]USBDevice, error) {
	links, err := os.ReadDir(sysfsUsbDevices)
	if err != nil {
		return nil, err
	}
	devs := make([]USBDevice, 0)
	for _, link := range links {
		if !s.usbDevRe.MatchString(link.Name()) || link.Type() != os.ModeSymlink {
			continue
		}
		abs := path.Join(sysfsUsbDevices, link.Name())
		node, err := filepath.EvalSymlinks(abs)
		if err != nil {
			return nil, err
		}
		dev := NewUSBDevice(node)
		err = dev.populate()
		if err != nil {
			return nil, err
		}
		devs = append(devs, dev)
	}
	return devs, nil
}
