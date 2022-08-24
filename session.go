package fridago

//#include "frida-core.h"
import "C"

type SessionOptions struct {
	handle *C.FridaSessionOptions
}

func NewSessionOptions() *SessionOptions {
	return &SessionOptions{handle: C.frida_session_options_new()}
}

func (s *SessionOptions) Free() {
	C.g_object_unref(C.gpointer(s.handle))
	s.handle = nil
}

func (s *SessionOptions) SetRealm(realm uint) {
	C.frida_session_options_set_realm(s.handle, C.FridaRealm(realm))
}

func (s *SessionOptions) SetPersistTimeout(timeout uint) {
	C.frida_session_options_set_persist_timeout(s.handle, C.guint(timeout))
}

type Session struct {
	handle       *C.FridaSession
	pid, timeout uint
}

func NewSession(s *C.FridaSession) *Session {
	return &Session{
		handle: s,
	}
}

func (s *Session) Pid() uint {
	if s.pid == 0 {
		s.pid = uint(C.frida_session_get_pid(s.handle))
	}
	return s.pid
}

func (s *Session) PersistTimeout() uint {
	if s.timeout == 0 {
		s.timeout = uint(C.frida_session_get_persist_timeout(s.handle))
	}
	return s.timeout
}

func (s *Session) IsDetached() bool {
	return Gbool(C.frida_session_is_detached(s.handle))
}

func (s *Session) Detach() error {
	var gerr *C.GError

	C.frida_session_detach_sync(s.handle, nil, &gerr)
	if gerr != nil {
		return NewGError(gerr)
	}
	C.g_object_unref(C.gpointer(s.handle))
	s.handle = nil
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
	fscript = C.frida_session_create_script_sync(s.handle, C.CString(source), opts, nil, &gerr)
	if gerr != nil {
		return nil, NewGError(gerr)
	}
	return NewScript(fscript, name), nil
}
