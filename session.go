package fridago

//#include "frida-core.h"
import "C"

type SessionOptions struct {
	ptr *C.FridaSessionOptions
}

func NewSessionOptions() *SessionOptions {
	return &SessionOptions{ptr: C.frida_session_options_new()}
}

func (s *SessionOptions) Free() {
	C.g_object_unref(C.gpointer(s.ptr))
	s.ptr = nil
}

func (s *SessionOptions) SetRealm(realm uint) {
	C.frida_session_options_set_realm(s.ptr, C.FridaRealm(realm))
}

func (s *SessionOptions) SetPersistTimeout(timeout uint) {
	C.frida_session_options_set_persist_timeout(s.ptr, C.guint(timeout))
}

type Session struct {
	ptr *C.FridaSession

	Dev *Device
	Pid uint
}

func NewSession(dev *Device, fs *C.FridaSession) *Session {
	return &Session{
		ptr: fs,
		Dev: dev,
		Pid: uint(C.frida_session_get_pid(fs)),
	}
}

func (s *Session) Detach() error {
	var gerr *C.GError

	C.frida_session_detach_sync(s.ptr, nil, &gerr)
	if gerr != nil {
		return NewGError(gerr)
	}
	C.g_object_unref(C.gpointer(s.ptr))
	s.ptr = nil
	return nil
}

func (s *Session) CreateScript(name string, source string, runtime ...uint) (*Script, error) {
	var fruntime C.FridaScriptRuntime
	if len(runtime) == 0 {
		fruntime = C.FRIDA_SCRIPT_RUNTIME_DEFAULT
	} else {
		fruntime = C.FridaScriptRuntime(runtime[0])
	}

	opts := C.frida_script_options_new()
	defer func() {
		C.g_object_unref(C.gpointer(opts))
		opts = nil
	}()
	C.frida_script_options_set_name(opts, C.CString(name))
	C.frida_script_options_set_runtime(opts, fruntime)

	var (
		gerr    *C.GError
		fscript *C.FridaScript
	)
	fscript = C.frida_session_create_script_sync(s.ptr, C.CString(source), opts, nil, &gerr)
	if gerr != nil {
		return nil, NewGError(gerr)
	}
	return NewScript(fscript, name), nil
}
