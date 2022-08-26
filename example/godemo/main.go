package main

// 参照例子：https://github.com/r0ysue/AndroidSecurityStudy/blob/master/FRIDA/A02/README.md#基本能力ⅰhook参数修改结果
import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Buzz2d0/fridago"
)

func main() {
	fmt.Println(fridago.Version())

	manager := fridago.NewDeviceManager()
	defer manager.Close()

	usbDevice, err := manager.GetUsbDevice()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(usbDevice.Description())
	defer usbDevice.Free()

	app := "com.zznq.demo01"
	pid, err := usbDevice.Spawn(app)
	if err != nil {
		log.Fatalf("Spawn %s pid: %d error:%v", app, pid, err)
	}
	usbDevice.Resume(pid)
	session, err := usbDevice.Attach(pid)
	if err != nil {
		log.Fatalf("Attach pid: %d error:%v", pid, err)
	}
	defer session.Detach()

	quit := make(chan struct{}, 1)
	session.SetOnDetachedHandler(quit)

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
	if err := script.Load(); err != nil {
		log.Println("load error: ", err)
	}
	script.SetOnMessageHandler(func(s string) {
		fmt.Println("[message]->", s)
	})

	var sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
	case <-sigs:
	}
	fmt.Println("fridago bye!")
}
