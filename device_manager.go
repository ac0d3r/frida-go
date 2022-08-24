package fridago

//#include "frida-core.h"
import "C"

type DeviceManager struct {
	ptr *C.FridaDeviceManager
}

func NewDeviceManager() *DeviceManager {
	return &DeviceManager{ptr: C.frida_device_manager_new()}
}

// TODO
// addRemoteDevice
// removeRemoteDevice

func (dm *DeviceManager) Close() error {
	var gerr *C.GError

	C.frida_device_manager_close_sync(dm.ptr, nil, &gerr)
	if gerr != nil {
		return NewGError(gerr)
	}

	C.g_object_unref(C.gpointer(dm.ptr))
	dm.ptr = nil
	return nil
}

func (dm *DeviceManager) EnumerateDevices() ([]*Device, error) {
	var gerr *C.GError

	devices := C.frida_device_manager_enumerate_devices_sync(dm.ptr, nil, &gerr)
	if gerr != nil {
		return nil, NewGError(gerr)
	}

	defer func() {
		C.g_object_unref(C.gpointer(devices))
		devices = nil
	}()

	size := int(C.frida_device_list_size(devices))
	dl := make([]*Device, size)
	for i := 0; i < size; i++ {
		fd := C.frida_device_list_get(devices, C.int(i))
		dl[i] = NewDevice(fd)
	}
	return dl, nil
}
