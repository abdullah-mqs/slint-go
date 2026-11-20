package main

/*
#include <stdlib.h>
// bridge.h is intentionally NOT included here:
// The preamble of a file with //export must only contain declarations,
// otherwise C type conflicts can occur.
*/
import "C"
import (
	"sync"
)

var (
	callbacks = map[string]func(string) string{}
	cbMu      sync.RWMutex
)

// GoCallbackHandler is called by C++ (Slint) when a UI callback fires.
// It looks up the Go function registered for this callback name and runs it.
// The returned string is allocated with C.CString, C++ side is responsible for free()ing it.
//
//export GoCallbackHandler
func GoCallbackHandler(cName *C.char, cInput *C.char) *C.char {
	name := C.GoString(cName)
	input := C.GoString(cInput)

	cbMu.RLock()
	fn, ok := callbacks[name]
	cbMu.RUnlock()
	if !ok {
		return C.CString("")
	}
	return C.CString(fn(input))
}
