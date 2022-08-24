package fridago

//#include "frida-core.h"
import "C"
import (
	"sync"
)

type ScriptMessageHandler func(string)

var (
	cbs sync.Map
)

//export onMessage
func onMessage(script *C.FridaScript, message *C.gchar) {
	callback, ok := cbs.Load(uintptr(C.gpointer(script)))
	if !ok {
		return
	}
	mh, ok := callback.(ScriptMessageHandler)
	if !ok {
		return
	}
	mh(C.GoString(message))
}
