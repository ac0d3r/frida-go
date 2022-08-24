package fridago

//#include "frida-core.h"
import "C"

import (
	"unsafe"
)

// frida stdio enum
const (
	FridaStdioInherit = uint(C.FRIDA_STDIO_INHERIT)
	FridaStdioPipe    = uint(C.FRIDA_STDIO_PIPE)
)

// frida realm enum
const (
	FridaRealmNative   = uint(C.FRIDA_REALM_NATIVE)
	FridaRealmEmulated = uint(C.FRIDA_REALM_EMULATED)
)

// frida script runtime enum
const (
	FridaScriptRuntimeDefault = uint(C.FRIDA_SCRIPT_RUNTIME_DEFAULT)
	FridaScriptRuntimeQJS     = uint(C.FRIDA_SCRIPT_RUNTIME_QJS)
	FridaScriptRuntimeV8      = uint(C.FRIDA_SCRIPT_RUNTIME_V8)
)

func GbooleanToBool(gb C.gboolean) bool {
	return int(gb) != 0
}

func GStringsToSlice(gs **C.gchar, length C.gint) []string {
	lens := int(length)
	if lens > 0 {
		// https://stackoverflow.com/questions/48756732/what-does-1-30c-yourtype-do-exactly-in-cgo
		arr := (*[1 << 30]*C.gchar)(unsafe.Pointer(gs))
		strs := make([]string, lens)
		for i := 0; i < lens; i++ {
			strs[i] = C.GoString(arr[i])
		}
		return strs
	}
	return nil
}

func SliceToGStrings(strs []string) (**C.gchar, C.gint) {
	buf := make([]*C.gchar, len(strs))
	for i := range strs {
		buf[i] = (*C.gchar)(unsafe.Pointer(C.CString(strs[i])))
	}
	gs := (**C.gchar)(unsafe.Pointer(&buf[0]))
	return gs, C.gint(len(strs))
}
