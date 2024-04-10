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
	defer C.free(unsafe.Pointer(libPath))

	// Rufe die Wrapper-Funktion auf
	lib := C.load_external_lib(libPath)
	errV := C.GoString(lib.err)
	if errV != "" {
		panic(errV)
	}

	// Die Daten werden Extrahiert
	name := C.GoString(lib.lib.name)
	version := int(lib.lib.version)

	fmt.Println(name, version)

	// Die Lib wird entladen
	C.unload_lib()
}
