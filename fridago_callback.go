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
	v, ok := cbs.Load(uintptr(C.gpointer(script)))
	if !ok {
		return
	}
	callback, ok := v.(ScriptMessageHandler)
	if !ok {
		return
	}
	callback(C.GoString(message))
}

//export onDetached
func onDetached(session *C.FridaSession) {
	v, ok := cbs.Load(uintptr(C.gpointer(session)))
	if !ok {
		return
	}
	ch, ok := v.(chan<- struct{})
	if !ok {
		return
	}
	ch <- struct{}{}
}
