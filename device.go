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
	handle *C.FridaDevice

	id   string
	name string
	kind DeviceType
}

func NewDevice(fd *C.FridaDevice) *Device {
	return &Device{handle: fd}
}

func (d *Device) Free() {
	C.g_object_unref(C.gpointer(d.handle))
	d.handle = nil
}

func (d *Device) ID() string {
	if d.id == "" {
		d.id = C.GoString(C.frida_device_get_id(d.handle))
	}
	return d.id
}

func (d *Device) Kind() DeviceType {
	if d.kind == 0 {
		d.kind = DeviceType(C.frida_device_get_dtype(d.handle))
	}
	return d.kind
}

func (d *Device) Name() string {
	if d.name == "" {
		d.name = C.GoString(C.frida_device_get_name(d.handle))
	}
	return d.name
}

func (d *Device) IsLost() bool {
	return cbool(C.frida_device_is_lost(d.handle))
}

func (d *Device) Description() string {
	return fmt.Sprintf("Frida.Device(id: \"%s\", name: \"%s\", kind: \"%s\")", d.ID(), d.Name(), d.Kind().String())
}

type SpawnOptions struct {
	Args  []string
	Envp  map[string]string
	Env   map[string]string
	Cwd   string
	Stdio uint
}

func (so SpawnOptions) SetTo(handle *C.FridaSpawnOptions) {
	if len(so.Args) != 0 {
		gchar, length := slice2carray(so.Args)
		C.frida_spawn_options_set_argv(handle, gchar, length)
	}
	if len(so.Cwd) != 0 {
		C.frida_spawn_options_set_cwd(handle, C.CString(so.Cwd))
	}
	C.frida_spawn_options_set_stdio(handle, C.FridaStdio(so.Stdio))
	// TODO
	// frida_spawn_options_set_env
	// frida_spawn_options_set_envp
}

func (d *Device) Spawn(program string, so ...SpawnOptions) (uint, error) {
	var gerr *C.GError
	opts := C.frida_spawn_options_new()
	defer func() {
		C.g_object_unref(C.gpointer(opts))
		opts = nil
	}()
	if len(so) != 0 {
		so[0].SetTo(opts)
	}

	pid := uint(C.frida_device_spawn_sync(d.handle, C.CString(program), opts, nil, &gerr))
	if gerr != nil {
		return 0, NewGError(gerr)
	}
	return pid, nil
}

func (d *Device) GetProcessByName(name string) (uint, error) {
	var gerr *C.GError
	opts := C.frida_process_match_options_new()
	defer func() {
		C.g_object_unref(C.gpointer(opts))
		opts = nil
	}()

	process := C.frida_device_get_process_by_name_sync(d.handle, C.CString(name), opts, nil, &gerr)
	if gerr != nil {
		return 0, NewGError(gerr)
	}

	pid := uint(C.frida_process_get_pid(process))
	C.g_object_unref(C.gpointer(process))
	process = nil

	return pid, nil
}

func (d *Device) Resume(pid uint) error {
	var gerr *C.GError

	C.frida_device_resume_sync(d.handle, C.guint(pid), nil, &gerr)
	if gerr != nil {
		return NewGError(gerr)
	}
	return nil
}

type SessionOptions struct {
	Realm   uint
	Timeout uint
}

func (so SessionOptions) SetTo(handle *C.FridaSessionOptions) {
	C.frida_session_options_set_realm(handle, C.FridaRealm(so.Realm))
	if so.Timeout != 0 {
		C.frida_session_options_set_persist_timeout(handle, C.guint(so.Timeout))
	}
}

func (d *Device) Attach(pid uint, so ...SessionOptions) (*Session, error) {
	var gerr *C.GError
	opts := C.frida_session_options_new()
	defer func() {
		C.g_object_unref(C.gpointer(opts))
		opts = nil
	}()
	if len(so) != 0 {
		so[0].SetTo(opts)
	}
	session := C.frida_device_attach_sync(d.handle, C.guint(pid), opts, nil, &gerr)
	if gerr != nil {
		return nil, NewGError(gerr)
	}
	return NewSession(session), nil
}
