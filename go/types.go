package main

/*
#cgo CFLAGS: -I${SRCDIR}/../bridge
#cgo LDFLAGS: -L${SRCDIR}/../build/dev -lbridge
#include <stdlib.h>
#include "bridge.h"
*/
import "C"
import "unsafe"

// SetStringArray: Go []string -> Slint [string]
func (b *Bridge) SetStringArray(name string, items []string) {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))

	cs := make([]*C.char, len(items))
	for i, s := range items {
		cs[i] = C.CString(s)
	}
	defer func() {
		for _, p := range cs {
			C.free(unsafe.Pointer(p))
		}
	}()

	if len(items) == 0 {
		C.slint_bridge_set_string_array(cn, nil, 0)
		return
	}
	C.slint_bridge_set_string_array(cn, (**C.char)(unsafe.Pointer(&cs[0])), C.int(len(items)))
}

// SetIntArray: Go []int -> Slint [int]
func (b *Bridge) SetIntArray(name string, items []int) {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))

	ci := make([]C.int, len(items))
	for i, v := range items {
		ci[i] = C.int(v)
	}
	if len(ci) == 0 {
		C.slint_bridge_set_int_array(cn, nil, 0)
		return
	}
	C.slint_bridge_set_int_array(cn, &ci[0], C.int(len(ci)))
}

// SetFloatArray: Go []float64 -> Slint [float]
func (b *Bridge) SetFloatArray(name string, items []float64) {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))

	if len(items) == 0 {
		C.slint_bridge_set_float_array(cn, nil, 0)
		return
	}
	C.slint_bridge_set_float_array(cn, (*C.double)(unsafe.Pointer(&items[0])), C.int(len(items)))
}

// SlintStruct
type SlintStruct struct {
	h unsafe.Pointer
}

func NewStruct() *SlintStruct {
	return &SlintStruct{h: C.slint_bridge_struct_new()}
}

// Free: struct use hone ke baad call karo (ya defer s.Free())
func (s *SlintStruct) Free() {
	if s.h != nil {
		C.slint_bridge_struct_free(s.h)
		s.h = nil
	}
}

func (s *SlintStruct) SetString(field, value string) {
	cf := C.CString(field)
	cv := C.CString(value)
	defer C.free(unsafe.Pointer(cf))
	defer C.free(unsafe.Pointer(cv))
	C.slint_bridge_struct_set_string(s.h, cf, cv)
}

func (s *SlintStruct) SetInt(field string, value int) {
	cf := C.CString(field)
	defer C.free(unsafe.Pointer(cf))
	C.slint_bridge_struct_set_int(s.h, cf, C.int(value))
}

func (s *SlintStruct) SetFloat(field string, value float64) {
	cf := C.CString(field)
	defer C.free(unsafe.Pointer(cf))
	C.slint_bridge_struct_set_float(s.h, cf, C.double(value))
}

func (s *SlintStruct) SetBool(field string, value bool) {
	cf := C.CString(field)
	defer C.free(unsafe.Pointer(cf))
	iv := 0
	if value {
		iv = 1
	}
	C.slint_bridge_struct_set_bool(s.h, cf, C.int(iv))
}

func (s *SlintStruct) SetColor(field, hex string) {
	cf := C.CString(field)
	ch := C.CString(hex)
	defer C.free(unsafe.Pointer(cf))
	defer C.free(unsafe.Pointer(ch))
	C.slint_bridge_struct_set_color(s.h, cf, ch)
}

// SetStruct
func (b *Bridge) SetStruct(name string, s *SlintStruct) {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	C.slint_bridge_set_struct(cn, s.h)
}

func (b *Bridge) SetEnum(name string, variant string) {
	b.SetString(name, variant)
}

// SetStructArray: []*SlintStruct -> Slint [Struct]
func (b *Bridge) SetStructArray(name string, items []*SlintStruct) {
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))

	if len(items) == 0 {
		C.slint_bridge_set_struct_array(cn, nil, 0)
		return
	}
	ptrs := make([]unsafe.Pointer, len(items))
	for i, s := range items {
		ptrs[i] = s.h
	}
	C.slint_bridge_set_struct_array(cn, (*unsafe.Pointer)(unsafe.Pointer(&ptrs[0])), C.int(len(ptrs)))
}
