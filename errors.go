package fridago

//#include "frida-core.h"
import "C"
import (
	"errors"
)

var (
	ErrServerNotRunning       = errors.New("server not running")
	ErrExecutableNotFound     = errors.New("executable not found")
	ErrExecutableNotSupported = errors.New("executable not supported")
	ErrProcessNotFound        = errors.New("process not found")
	ErrProcessNotResponding   = errors.New("process not responding")
	ErrInvalidArgument        = errors.New("invalid argument")
	ErrInvalidOperation       = errors.New("invallid operation")
	ErrPermissionDenied       = errors.New("permission denied")
	ErrAddressInUse           = errors.New("address in use")
	ErrTimedOut               = errors.New("timeout error")
	ErrNotSupported           = errors.New("not supported")
	ErrProtocolViolation      = errors.New("protocol violation")
	ErrTransport              = errors.New("transport error")
)

type GError struct {
	Msg  string
	Code int
}

func NewGError(e *C.GError) error {
	return GError{
		Msg:  C.GoString(e.message),
		Code: int(e.code),
	}
}

func (err GError) Error() string {
	return err.Msg
}
