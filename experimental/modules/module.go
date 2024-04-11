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
	lib := C.cgo_load_external_lib(libPath)
	if C.GoString(lib.err) != "" {
		panic(C.GoString(lib.err))
	}

	// Es werden alle Verfügbaren Funktionen abgerufen
	sharedFunctions := C.cgo_get_global_functions()
	for i := 0; i < int(sharedFunctions.size); i++ {
		// Berechne die Adresse des i-ten Elements
		functionPtr := (*C.C_VM_FUNCTION)(unsafe.Pointer(uintptr(unsafe.Pointer(sharedFunctions.array)) + uintptr(i)*unsafe.Sizeof(*sharedFunctions.array)))

		// Die Funktion wird aufgerufen
		res := C.cgo_call_function(functionPtr)
		fmt.Println("RES:", res)
	}

	// Wird aufgerufen wenn die Funktion zueende ist
	defer func() {
		C.cgo_unload_lib()
	}()

	name := C.GoString(lib.name)
	fmt.Println(name, lib.version)
}
