package fridago

/*
#include "frida-core.h"

extern void onMessage(FridaScript * script, const gchar * message);

void _script_on_message(FridaScript * script,
	const gchar * message,
	GBytes * data,
	gpointer user_data)
{
	JsonParser *parser;
	JsonObject *root;
	const gchar *type;

	parser = json_parser_new();
	json_parser_load_from_data(parser, message, -1, NULL);
	root = json_node_get_object(json_parser_get_root(parser));


	type = json_object_get_string_member(root, "type");
	if (strcmp(type, "log") == 0)
	{
		const gchar *log_message;
		log_message = json_object_get_string_member(root, "payload");
		onMessage(script, log_message);
	}else
	{
		onMessage(script, message);
	}
	g_object_unref(parser);
}
*/
import "C"
import (
	"unsafe"
)

type Script struct {
	handle           *C.FridaScript
	onMessageHandler C.gulong

	Name string
}

func NewScript(s *C.FridaScript, name string) *Script {
	return &Script{
		handle: s,
		Name:   name,
	}
}

func (s *Script) Free() {
	handlers := []C.gulong{s.onMessageHandler}
	for i := range handlers {
		if handlers[i] != 0 {
			C.g_signal_handler_disconnect(C.gpointer(s.handle), handlers[i])
			handlers[i] = 0
		}
	}
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

func (s *Script) SetOnMessageHandler(callback ScriptMessageHandler) {
	if s.onMessageHandler == 0 {
		signal := C.CString("message")
		defer C.free(unsafe.Pointer(signal))
		s.onMessageHandler = C.g_signal_connect_data(
			C.gpointer(s.handle),
			signal,
			C.GCallback(unsafe.Pointer(C._script_on_message)),
			nil, nil, 0)
	}
	cbs.Store(uintptr(C.gpointer(s.handle)), callback)
}
