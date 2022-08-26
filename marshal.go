package fridago

//#include "frida-core.h"
import "C"

import (
	"unsafe"
)

func cbool(gb C.gboolean) bool {
	return int(gb) != 0
}

func carray2slice(gs **C.gchar, length C.gint) []string {
	lens := int(length)
	arr := unsafe.Slice(gs, length)
	strs := make([]string, lens)
	for i := 0; i < lens; i++ {
		strs[i] = C.GoString(arr[i])
	}
	return strs
}

func slice2carray(strs []string) (**C.gchar, C.gint, []unsafe.Pointer) {
	buf := make([]*C.gchar, len(strs))
	ps := make([]unsafe.Pointer, len(strs))
	for i := range strs {
		ps[i] = unsafe.Pointer(C.CString(strs[i]))
		buf[i] = (*C.gchar)(ps[i])
	}
	gs := (**C.gchar)(unsafe.Pointer(&buf[0]))
	return gs, C.gint(len(strs)), ps
}
