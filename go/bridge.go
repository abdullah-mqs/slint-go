package main

/*
#cgo CFLAGS: -I${SRCDIR}/../bridge
#cgo LDFLAGS: -L${SRCDIR}/../build/dev -lbridge
#include <stdlib.h>
#include "bridge.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Bridge struct{}

func NewBridge() *Bridge { return &Bridge{} }

func (b *Bridge) Init() { C.slint_bridge_init() }

func (b *Bridge) LoadUI(path string) error {
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))
	if C.slint_bridge_load(cs) != 0 {
		return fmt.Errorf("UI load fail: %s", path)
	}
	return nil
}

func (b *Bridge) Run() { C.slint_bridge_run() }

func (b *Bridge) OnCallback(name string, fn func(string) string) {
	cbMu.Lock()
	callbacks[name] = fn
	cbMu.Unlock()
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	C.slint_bridge_on_callback(cn)
}

func (b *Bridge) SetString(name string, v string) {
	cn := C.CString(name)
	cv := C.CString(v)
	defer C.free(unsafe.Pointer(cn))
	defer C.free(unsafe.Pointer(cv))
	C.slint_bridge_set_string(cn, cv)
}

func (b *Bridge) SetInt(name string, v int) {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	C.slint_bridge_set_int(cn, C.int(v))
}

func (b *Bridge) SetBool(name string, v bool) {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	iv := 0
	if v {
		iv = 1
	}
	C.slint_bridge_set_bool(cn, C.int(iv))
}

func (b *Bridge) SetFloat(name string, v float64) {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	C.slint_bridge_set_float(cn, C.double(v))
}

func (b *Bridge) SetColor(name string, hex string) {
	cn := C.CString(name)
	cv := C.CString(hex)
	defer C.free(unsafe.Pointer(cn))
	defer C.free(unsafe.Pointer(cv))
	C.slint_bridge_set_color(cn, cv)
}

func (b *Bridge) SetImage(name string, path string) {
	cn := C.CString(name)
	cp := C.CString(path)
	defer C.free(unsafe.Pointer(cn))
	defer C.free(unsafe.Pointer(cp))
	C.slint_bridge_set_image(cn, cp)
}

func (b *Bridge) SetImageBytes(name string, data []byte) {
	if len(data) == 0 {
		return
	}
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	C.slint_bridge_set_image_bytes(cn, (*C.uchar)(unsafe.Pointer(&data[0])), C.int(len(data)))
}

func (b *Bridge) GetString(name string) string {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	cs := C.slint_bridge_get_string(cn)
	if cs == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cs))
	return C.GoString(cs)
}

func (b *Bridge) GetInt(name string) int {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	return int(C.slint_bridge_get_int(cn))
}

func (b *Bridge) GetBool(name string) bool {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	return C.slint_bridge_get_bool(cn) != 0
}
