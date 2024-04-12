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
	cmodule          *C.CWrappedModuleLib
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
	C.cgo_unload_lib(o.cmodule)

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

func (o *CGOWrappedLibModuleFunction) Call() (string, interface{}, error) {
	// Setup eines defer-Blocks zur Abfangung von Panics
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("LoadWrappedCGOLibModule failed: %v", r)
		}
	}()

	// Die Funktion wird aufgerufen
	retData := C.cgo_call_function(o.functionPtr)

	// Prüfe den Datentyp und handle entsprechend
	var value interface{}
	var valueType string
	switch retData._type {
	case C.NONE:
		valueType = "null"
		value = nil
	case C.STRING:
		value = C.GoString(retData.string_data)
		valueType = "string"
	case C.ERROR:
		value = C.GoString(retData.error_data)
		valueType = "error"
	case C.BYTES:
		value = retData.byte_data
		valueType = "bytes"
	case C.INT:
		value = int(retData.int_data)
		valueType = "int"
	case C.FLOAT:
		value = float64(retData.float_data)
		valueType = "float"
	case C.BOOLEAN:
		value = bool(retData.bool_data)
		valueType = "bool"
	case C.TIMESTAMP:
		value = C.GoString(retData.timestamp_data)
		valueType = "timestamp"
	case C.OBJECT:
		value = retData.object_data
		valueType = "object"
	case C.ARRAY:
		value = retData.array_data
		valueType = "array"
	default:
		return "", nil, fmt.Errorf("unkown datatype")
	}

	// Die Daten werden zurückgegeben
	return valueType, value, err
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
	cgo_global_functions := C.cgo_get_global_functions(lib.moduleLib)
	for i := 0; i < int(cgo_global_functions.size); i++ {
		// Berechne die Adresse des i-ten Elements
		functionPtr := (*C.CVmFunction)(unsafe.Pointer(uintptr(unsafe.Pointer(cgo_global_functions.array)) + uintptr(i)*unsafe.Sizeof(*cgo_global_functions.array)))

		// Die Funktion wird zwischengespeichert
		c_functions = append(c_functions, &CGOWrappedLibModuleFunction{name: C.GoString(functionPtr.name), functionPtr: functionPtr})
	}

	// Das Rückgabe Objekt wird erstellt
	returnObj := &CGOWrappedLibModule{
		path:             libPath,
		global_functions: c_functions,
		lib:              unsafe.Pointer(libPath),
		name:             C.GoString(lib.name),
		version:          uint(lib.version),
		cmodule:          lib.moduleLib,
	}

	// Das Objekt wird zurückgegeben
	return returnObj, err
}
