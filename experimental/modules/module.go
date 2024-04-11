package modules

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

func TestLoad() {
	// Pfad zur Shared Library
	libPath := C.CString("/home/fluffelbuff/Schreibtisch/lib.so")
	defer func() {
		C.free(unsafe.Pointer(libPath))
	}()

	// Rufe die Wrapper-Funktion auf
	lib := C.load_external_lib(libPath)
	if C.GoString(lib.err) != "" {
		panic(C.GoString(lib.err))
	}

	// Es werden alle Verfügbaren Funktionen abgerufen
	sharedFunctions := C.get_global_functions()
	fmt.Println(sharedFunctions)

	// Wird aufgerufen wenn die Funktion zueende ist
	defer func() {
		C.unload_lib()
	}()

	name := C.GoString(lib.name)
	fmt.Println(name, lib.version)
}
