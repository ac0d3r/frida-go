package main

import (
	"fmt"

	"github.com/Buzz2d0/fridago"
)

func main() {
	manager := fridago.NewDeviceManager()
	defer manager.Close()

	// testRemoteDevice(manager)
	testGetProcessByName(manager)
}

func testGetProcessByName(manager *fridago.DeviceManager) {
	d, err := manager.GetUsbDevice()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(d.GetProcessByName("FridaGoDemo01"))
}

func testRemoteDevice(manager *fridago.DeviceManager) {
	devive, err := manager.AddRemoteDevice("10.0.2.16:1456")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(devive.Description())
	defer devive.Free()

	devices, err := manager.EnumerateDevices()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, d := range devices {
		fmt.Println(d.Description())
	}
}
