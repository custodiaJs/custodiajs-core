package main

/*
#cgo LDFLAGS: -ldl
#include <stdlib.h>
#include "wrapper.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

//export myGoCallback
func myGoCallback() {
	fmt.Println("Dies ist mein Go-Callback!")
}

func main() {
	// Pfad zur Shared Library
	libPath := C.CString("/home/fluffelbuff/Dokumente/Projekte/vnh1/testlib/clib/lib.so")
	defer C.free(unsafe.Pointer(libPath))

	// Rufe die Wrapper-Funktion auf
	lib := C.load_external_lib(libPath)
	fmt.Println(lib)

	C.callGoCallback(C.callback_func(unsafe.Pointer(C.myGoCallback)))

	// Die Lib wird entladen
	C.unload_lib()
}
