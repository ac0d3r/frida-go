package main

// 参照例子：https://github.com/r0ysue/AndroidSecurityStudy/blob/master/FRIDA/A02/README.md#基本能力ⅰhook参数修改结果
import (
	"fmt"
	"log"

	"github.com/Buzz2d0/fridago"
)

func main() {
	manager := fridago.NewDeviceManager()
	defer manager.Close()

	devices, err := manager.EnumerateDevices()
	if err != nil {
		log.Fatal(err)
	}

	var device *fridago.Device
	for _, d := range devices {
		if d.Kind == fridago.DeviceTypeUsb {
			device = d
			continue
		}
		d.Free()
	}
	if device == nil {
		return
	}

	log.Println(device.Description())
	defer device.Free()
	pid, err := device.Spawn("com.zznq.demo01")
	if err != nil {
		log.Fatal(err)
	}
	device.Resume(pid)
	session, err := device.Attach(pid)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Detach()

	var javascript = `
	console.log("Script loaded successfully ");
	Java.perform(function x() {
		console.log("Inside java perform function");
		var main_activity = Java.use("com.zznq.demo01.MainActivity");
		main_activity.add.implementation = function(x,y){
			console.log("original call: add("+ x + ", " + y + ")");
			return this.add(500, 20);
		}
	});
	`
	script, err := session.CreateScript("test", javascript)
	if err != nil {
		log.Fatal(err)
	}
	defer script.UnLoad()
	script.Load()

	// hang
	var s string
	fmt.Print("hang...(Enter cancel)")
	fmt.Scanln(&s)
}
