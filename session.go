package fridago

/*
#include "frida-core.h"

extern void onDetached(FridaSession *session);

void _session_on_detached(FridaSession *session,
            FridaSessionDetachReason reason,
            FridaCrash *crash,
            gpointer user_data)
{
	onDetached(session);
}
*/
import "C"
import "unsafe"

type Session struct {
	handle            *C.FridaSession
	onDetachedHandler C.gulong
	pid, timeout      uint
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
	return cbool(C.frida_session_is_detached(s.handle))
}

func (s *Session) Detach() error {
	var gerr *C.GError

	C.frida_session_detach_sync(s.handle, nil, &gerr)
	if gerr != nil {
		return NewGError(gerr)
	}

	if s.onDetachedHandler != 0 {
		C.g_signal_handler_disconnect(C.gpointer(s.handle), s.onDetachedHandler)
		s.onDetachedHandler = 0
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
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.frida_script_options_set_name(opts, cname)
	C.frida_script_options_set_runtime(opts, fruntime)

	var (
		gerr    *C.GError
		fscript *C.FridaScript
	)
	csource := C.CString(source)
	defer C.free(unsafe.Pointer(csource))
	fscript = C.frida_session_create_script_sync(s.handle, csource, opts, nil, &gerr)
	if gerr != nil {
		return nil, NewGError(gerr)
	}
	return NewScript(fscript, name), nil
}

func (s *Session) SetOnDetachedHandler(ch chan<- struct{}) {
	if s.onDetachedHandler == 0 {
		signal := C.CString("detached")
		defer C.free(unsafe.Pointer(signal))
		s.onDetachedHandler = C.g_signal_connect_data(
			C.gpointer(s.handle),
			signal,
			C.GCallback(unsafe.Pointer(C._session_on_detached)),
			nil, nil, 0)
	}
	cbs.Store(uintptr(C.gpointer(s.handle)), ch)
}
