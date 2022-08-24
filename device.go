package fridago

//#include "frida-core.h"
import "C"
import (
	"fmt"
)

type DeviceType uint

func (dt DeviceType) String() string {
	switch dt {
	case DeviceTypeLocal:
		return "local"
	case DeviceTypeRemote:
		return "remote"
	case DeviceTypeUsb:
		return "usb"
	default:
		return "Unexpected frida device kind"
	}
}

const (
	DeviceTypeLocal  = DeviceType(C.FRIDA_DEVICE_TYPE_LOCAL)
	DeviceTypeRemote = DeviceType(C.FRIDA_DEVICE_TYPE_REMOTE)
	DeviceTypeUsb    = DeviceType(C.FRIDA_DEVICE_TYPE_USB)
)

type Device struct {
	ptr *C.FridaDevice

	ID   string
	Name string
	Kind DeviceType
}

func NewDevice(fd *C.FridaDevice) *Device {
	d := &Device{ptr: fd}
	d.fromFridaDevice()
	return d
}

func (d *Device) Free() {
	C.g_object_unref(C.gpointer(d.ptr))
	d.ptr = nil
}

func (d *Device) IsLost() bool {
	return GbooleanToBool(C.frida_device_is_lost(d.ptr))
}

func (d *Device) Description() string {
	return fmt.Sprintf("Frida.Device(id: \"%s\", name: \"%s\", kind: \"%s\")", d.ID, d.Name, d.Kind.String())
}

func (d *Device) Spawn(program string, opts ...*SpawnOptions) (uint, error) {
	var (
		gerr *C.GError
		opt  *SpawnOptions
	)
	if len(opts) == 0 {
		opt = NewSpawnOptions()
	} else {
		opt = opts[0]
	}
	defer opt.Free()

	pid := uint(C.frida_device_spawn_sync(d.ptr, C.CString(program), opt.ptr, nil, &gerr))
	if gerr != nil {
		return 0, NewGError(gerr)
	}
	return pid, nil
}

func (d *Device) Resume(pid uint) error {
	var gerr *C.GError

	C.frida_device_resume_sync(d.ptr, C.guint(pid), nil, &gerr)
	if gerr != nil {
		return NewGError(gerr)
	}
	return nil
}

func (d *Device) Attach(pid uint, opts ...*SessionOptions) (*Session, error) {
	var (
		gerr *C.GError
		opt  *SessionOptions
	)
	if len(opts) == 0 {
		opt = NewSessionOptions()
	} else {
		opt = opts[0]
	}
	defer opt.Free()

	session := C.frida_device_attach_sync(d.ptr, C.guint(pid), opt.ptr, nil, &gerr)
	if gerr != nil {
		return nil, NewGError(gerr)
	}
	return NewSession(d, session), nil
}

func (d *Device) fromFridaDevice() {
	d.Name = C.GoString(C.frida_device_get_name(d.ptr))
	d.ID = C.GoString(C.frida_device_get_id(d.ptr))
	d.Kind = DeviceType(C.frida_device_get_dtype(d.ptr))
	return
}
