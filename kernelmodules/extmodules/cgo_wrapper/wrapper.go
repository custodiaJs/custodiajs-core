package cgowrapper

/*
#include <stdlib.h>
#include "wrapper.h"
*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"
	"vnh1/filesystem"
	"vnh1/types"
)

type CGOWrappedLibModuleFunction struct {
	name        string
	functionPtr *C.CVmFunction
}

type CGOWrappedLibFunctionReturn struct {
	dtype string
	data  interface{}
}

type CGOWrappedLibFunctionParameter struct {
	dtype string
	data  interface{}
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

func (o *CGOWrappedLibModuleFunction) Call(parms ...*CGOWrappedLibFunctionParameter) (*CGOWrappedLibFunctionReturn, error) {
	// Setup eines defer-Blocks zur Abfangung von Panics
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("LoadWrappedCGOLibModule failed: %v", r)
		}
	}()

	// Es wird eine neue Liste an Parametern erzeugt
	cgoParmList := C.cgo_new_function_parm_list()

	// Die Funktion wird aufgerufen
	retData := C.cgo_call_function(o.functionPtr, cgoParmList)

	// Es wird geprüft ob ein CGO-Calling Panic aufgetreten ist
	if err != nil {
		return nil, &types.ExtModCGOPanic{ErrorValue: fmt.Errorf("CGOWrappedLibModuleFunction->Call: " + err.Error())}
	}

	// Es wird ermittelt ob ein Fehler aufgetreten ist
	if retData.ErrorMsg != nil {
		return nil, &types.ExtModCGOPanic{ErrorValue: fmt.Errorf("CGOWrappedLibModuleFunction->Call: " + C.GoString(retData.ErrorMsg))}
	}

	// Prüfe den Datentyp und handle entsprechend
	var value *CGOWrappedLibFunctionReturn
	switch retData.returnData._type {
	case C.NONE:
		value = &CGOWrappedLibFunctionReturn{dtype: "null"}
	case C.STRING:
		value = &CGOWrappedLibFunctionReturn{dtype: "string", data: C.GoString(retData.returnData.string_data)}
	case C.ERROR:
		value = &CGOWrappedLibFunctionReturn{dtype: "error", data: C.GoString(retData.returnData.error_data)}
	case C.BYTES:
		value = &CGOWrappedLibFunctionReturn{dtype: "bytes", data: retData.returnData.byte_data}
	case C.INT:
		value = &CGOWrappedLibFunctionReturn{dtype: "int", data: retData.returnData.int_data}
	case C.FLOAT:
		value = &CGOWrappedLibFunctionReturn{dtype: "float", data: retData.returnData.float_data}
	case C.BOOLEAN:
		value = &CGOWrappedLibFunctionReturn{dtype: "bool", data: retData.returnData.bool_data}
	case C.TIMESTAMP:
		value = &CGOWrappedLibFunctionReturn{dtype: "timestamp", data: C.GoString(retData.returnData.timestamp_data)}
	case C.OBJECT:
		value = &CGOWrappedLibFunctionReturn{dtype: "object", data: retData.returnData.object_data}
	case C.ARRAY:
		value = &CGOWrappedLibFunctionReturn{dtype: "array", data: retData.returnData.array_data}
	case C.FUNCTION:
	case C.UINT:
	default:
		return nil, &types.ExtModFunctionCallError{ErrorValue: fmt.Errorf("unkown datatype")}
	}

	// Die Daten werden zurückgegeben
	return value, nil
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

	// Der Datentyp des Module Libs wird ermittelt, danach wird versucht
	// passend zum Datentyp die Lib einzulesen
	var lib C.STARTUP_RESULT
	var libPath *C.char
	if filesystem.FileExists(pathv) {
		// Pfad zur Shared Library
		libPath = C.CString(pathv)

		// Es wird überprüft, ob es sich um eine UNIX .SO-Datei,
		// eine Windows .DLL-Datei, oder um eine Apple DYLIB-Datei handelt.
		// Außerdem wird überprüft, ob es sich um eine .PY-Skriptdatei oder um ein Python-Modul handelt.
		// Sollte es sich nicht um eines dieser Dateiformate handeln, wird ein Fehler ausgelöst.
		if strings.HasSuffix(pathv, ".so") {
			// Es wird geprüft ob es sich um eine Unix .so Datei handelt
			if !filesystem.IsUnixSOFile(pathv) {
				defer C.free(unsafe.Pointer(libPath))
				return nil, fmt.Errorf("LoadWrappedCGOLibModule: invalid lib file")
			}

			// Rufe die Wrapper-Funktion auf
			lib = C.cgo_load_external_dynamic_unix_lib(libPath)
			if C.GoString(lib.err) != "" {
				defer C.free(unsafe.Pointer(libPath))
				return nil, fmt.Errorf("LoadWrappedCGOLibModule: lib loading c error: %s", C.GoString(lib.err))
			}
		} else if strings.HasSuffix(pathv, ".dll") {
			// Es wird geprüft ob es sich um eine Windows .DLL handelt
			if !filesystem.IsWindowsDLL(pathv) {
				defer C.free(unsafe.Pointer(libPath))
				return nil, fmt.Errorf("LoadWrappedCGOLibModule: not supported data binary format")
			}

			// Rufe die Wrapper-Funktion auf
			lib = C.cgo_load_external_win32_dynamic_lib(libPath)
			if C.GoString(lib.err) != "" {
				defer C.free(unsafe.Pointer(libPath))
				return nil, fmt.Errorf("LoadWrappedCGOLibModule: lib loading c error: %s", C.GoString(lib.err))
			}
		} else if strings.HasSuffix(pathv, ".dylib") {
			// Es wird gerprüft ob der Header korrekt ist
			if !filesystem.IsDylib(pathv) {
				defer C.free(unsafe.Pointer(libPath))
				return nil, fmt.Errorf("LoadWrappedCGOLibModule: not supported data binary format")
			}

			// Rufe die Wrapper-Funktion auf
			lib = C.cgo_load_external_macos_dynamic_lib(libPath)
			if C.GoString(lib.err) != "" {
				defer C.free(unsafe.Pointer(libPath))
				return nil, fmt.Errorf("LoadWrappedCGOLibModule: lib loading c error: %s", C.GoString(lib.err))
			}
		} else {
			defer C.free(unsafe.Pointer(libPath))
			return nil, fmt.Errorf("LoadWrappedCGOLibModule: not supported data binary format")
		}
	} else {
		return nil, fmt.Errorf("LoadWrappedCGOLibModule: unkown path %s", pathv)
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
