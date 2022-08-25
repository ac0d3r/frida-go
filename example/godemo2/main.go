package main

import (
	"fmt"
	"log"

	"github.com/Buzz2d0/fridago"
)

func main() {
	manager := fridago.NewDeviceManager()
	defer manager.Close()

	devive, err := manager.AddRemoteDevice("10.0.2.16:1456")
	if err != nil {
		fmt.Println("errror: ", err)
		return
	}
	fmt.Println(devive.Description())
	defer devive.Free()

	devices, err := manager.EnumerateDevices()
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range devices {
		fmt.Println(d.Description())
	}
}
