package fridago

//#include "frida-core.h"
import "C"

type Spawn struct {
	ptr *C.FridaSpawn

	Identifier string
	Pid        uint
}

func NewSpawn(fs *C.FridaSpawn) *Spawn {
	s := &Spawn{ptr: fs}
	s.fromFridaSpawn()
	return s
}

func (s *Spawn) fromFridaSpawn() {
	s.Identifier = C.GoString(C.frida_spawn_get_identifier(s.ptr))
	s.Pid = uint(C.frida_spawn_get_pid(s.ptr))
}
