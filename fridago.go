package fridago

//#include "frida-core.h"
import "C"

func Version() string {
	return C.GoString(C.frida_version_string())
}

// frida stdio enum
const (
	StdioInherit = uint(C.FRIDA_STDIO_INHERIT)
	StdioPipe    = uint(C.FRIDA_STDIO_PIPE)
)

// frida realm enum
const (
	RealmNative   = uint(C.FRIDA_REALM_NATIVE)
	RealmEmulated = uint(C.FRIDA_REALM_EMULATED)
)

// frida scipt runtime enum
const (
	ScriptRuntimeDefault = uint(C.FRIDA_SCRIPT_RUNTIME_DEFAULT)
	ScriptRuntimeQJS     = uint(C.FRIDA_SCRIPT_RUNTIME_QJS)
	ScriptRuntimeV8      = uint(C.FRIDA_SCRIPT_RUNTIME_V8)
)
