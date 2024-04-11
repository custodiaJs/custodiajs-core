package cgowrapper

/*
#cgo CFLAGS: -I./c_lib
#include <stdlib.h>
#include "wrapper.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type CGOWrappedLibModuleFunction struct {
	name        string
	functionPtr *C.CVmFunction
}

type CGOWrappedLibModule struct {
	path             *C.char
	lib              unsafe.Pointer
	unload_lib       func()
	name             string
	version          uint
	global_functions []*CGOWrappedLibModuleFunction
}

func (o *CGOWrappedLibModule) Unload() {
	// Setup eines defer-Blocks zur Abfangung von Panics
	defer func() {
		// Überprüfe, ob ein Panic aufgetreten ist
		if r := recover(); r != nil {
			// Wenn ja, konvertiere das Recovered-Objekt in einen Fehler
			fmt.Println(fmt.Errorf("LoadWrappedCGOLibModule failed: %v", r))
		}
	}()

	// Es wird versucht die LIB zu entladen
	o.unload_lib()

	// Es wird aufgeärumt
	C.free(unsafe.Pointer(o.path))
}

func (o *CGOWrappedLibModule) GetName() string {
	return o.name
}

func (o *CGOWrappedLibModule) GetVersion() uint {
	return o.version
}

func (o *CGOWrappedLibModule) GetGlobalFunctions() []*CGOWrappedLibModuleFunction {
	return o.global_functions
}

func (o *CGOWrappedLibModuleFunction) Call() (string, error) {
	// Setup eines defer-Blocks zur Abfangung von Panics
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(fmt.Errorf("LoadWrappedCGOLibModule failed: %v", r))
		}
	}()

	retData := C.cgo_call_function(o.functionPtr)

	// Prüfe den Datentyp und handle entsprechend
	switch retData._type {
	case C.NONE:
		fmt.Println("Keine Daten zurückgegeben")
	case C.STRING:
		fmt.Println("String Daten:", C.GoString(retData.string_data))
	case C.ERROR:
		fmt.Println("Fehlermeldung:", C.GoString(retData.error_data))
	case C.BYTES:
		// Handle byte data, assuming it returns a null-terminated string for simplicity
		fmt.Println("Byte Daten:", C.GoString(retData.byte_data))
	case C.INT:
		fmt.Println("Integer Daten:", int(retData.int_data))
	case C.FLOAT:
		fmt.Println("Float Daten:", float64(retData.float_data))
	case C.BOOLEAN:
		fmt.Println("Boolean Daten:", bool(retData.bool_data))
	case C.TIMESTAMP:
		fmt.Println("Timestamp Daten:", C.GoString(retData.timestamp_data))
	case C.OBJECT:
		fmt.Println("Object Daten:", retData.object_data)
	case C.ARRAY:
		fmt.Println("Array Daten:", retData.array_data)
	default:
		fmt.Println("Unbekannter Typ")
	}

	// Die Daten werden zurückgegeben
	return "", nil
}

func (o *CGOWrappedLibModuleFunction) GetName() string {
	return o.name
}

func LoadWrappedCGOLibModule(pathv string) (*CGOWrappedLibModule, error) {
	// Initialisiere einen Rückgabewert für den Fehler auf nil
	var err error

	// Setup eines defer-Blocks zur Abfangung von Panics
	defer func() {
		// Überprüfe, ob ein Panic aufgetreten ist
		if r := recover(); r != nil {
			// Wenn ja, konvertiere das Recovered-Objekt in einen Fehler
			err = fmt.Errorf("LoadWrappedCGOLibModule failed: %v", r)
		}
	}()

	// Pfad zur Shared Library
	libPath := C.CString(pathv)

	// Rufe die Wrapper-Funktion auf
	lib := C.cgo_load_external_lib(libPath)
	if C.GoString(lib.err) != "" {
		panic(C.GoString(lib.err))
	}

	// Es werden alle Verfügbaren Funktionen abgerufen
	c_functions := make([]*CGOWrappedLibModuleFunction, 0)
	cgo_global_functions := C.cgo_get_global_functions()
	for i := 0; i < int(cgo_global_functions.size); i++ {
		// Berechne die Adresse des i-ten Elements
		functionPtr := (*C.CVmFunction)(unsafe.Pointer(uintptr(unsafe.Pointer(cgo_global_functions.array)) + uintptr(i)*unsafe.Sizeof(*cgo_global_functions.array)))

		// Die Funktion wird zwischengespeichert
		c_functions = append(c_functions, &CGOWrappedLibModuleFunction{name: C.GoString(functionPtr.name), functionPtr: functionPtr})
	}

	// Wird ausgeführt um die LIB zu Entladen
	unloadFunc := func() {
		C.cgo_unload_lib()
	}

	// Das Rückgabe Objekt wird erstellt
	returnObj := &CGOWrappedLibModule{
		path:             libPath,
		global_functions: c_functions,
		lib:              unsafe.Pointer(libPath),
		name:             C.GoString(lib.name),
		version:          uint(lib.version),
		unload_lib:       unloadFunc,
	}

	// Das Objekt wird zurückgegeben
	return returnObj, err
}
