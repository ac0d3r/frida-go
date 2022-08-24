package fridago

//#include "frida-core.h"
import "C"

type Script struct {
	handle *C.FridaScript

	Name string
}

func NewScript(s *C.FridaScript, name string) *Script {
	return &Script{
		handle: s,
		Name:   name,
	}
}

func (s *Script) Free() {
	C.g_object_unref(C.gpointer(s.handle))
	s.handle = nil
}

func (s *Script) Load() error {
	var gerr *C.GError
	C.frida_script_load_sync(s.handle, nil, &gerr)
	if gerr != nil {
		return NewGError(gerr)
	}
	return nil
}

func (s *Script) UnLoad() error {
	var gerr *C.GError
	C.frida_script_unload_sync(s.handle, nil, &gerr)
	if gerr != nil {
		return NewGError(gerr)
	}
	s.Free()
	return nil
}
