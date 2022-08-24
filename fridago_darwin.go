package fridago

/*
#cgo CFLAGS: -I${SRCDIR}/cfrida
#cgo LDFLAGS: -L${SRCDIR}/cfrida -lfrida-core -lbsm -ldl -lm -lresolv -framework Foundation -framework AppKit

#include "frida-core.h"
*/
import "C"

func init() {
	C.frida_init()
}
