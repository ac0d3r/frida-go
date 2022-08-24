package fridago

//#include "frida-core.h"
import "C"

type SpawnOptions struct {
	ptr *C.FridaSpawnOptions
}

func NewSpawnOptions() *SpawnOptions {
	return &SpawnOptions{ptr: C.frida_spawn_options_new()}
}

func (s *SpawnOptions) Free() {
	C.g_object_unref(C.gpointer(s.ptr))
	s.ptr = nil
}

// TODO
// frida_spawn_options_set_env
// frida_spawn_options_set_envp

func (s *SpawnOptions) SetArgv(args []string) {
	gchar, length := SliceToGStrings(args)
	C.frida_spawn_options_set_argv(s.ptr, gchar, length)
}

func (s *SpawnOptions) SetCwd(cwd string) {
	C.frida_spawn_options_set_cwd(s.ptr, C.CString(cwd))
}

func (s *SpawnOptions) SetStdio(stdio uint) {
	C.frida_spawn_options_set_stdio(s.ptr, C.FridaStdio(stdio))
}

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
