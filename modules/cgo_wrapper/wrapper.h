// wrapper.h
#include "lib_bridge.h"

typedef struct {
    const char* err;
    const char* name;
    uint version;
} STARTUP_RESULT;

typedef VmModule* (*LIB_LOAD)();
typedef void (*LIB_STOP)();

// Ruft die in bereitgestellte Funktion auf
CFunctionReturnData cgo_call_function(CVmFunction* function);

// Wird verwenet um die Lib zu Laden
STARTUP_RESULT cgo_load_external_lib(const char* lib_path);

// Gibt alle Verfügbaren Globalen Funktionen zurück
CVmFunctionList cgo_get_global_functions();

// Gibt alle Verfügbaren Globalen Objekte zurück
CVmObjectList cgo_get_global_object();

// Gibt alle Verfügbaren Modules zurück
CVmModulesList cgo_get_modules();

// Entlädt die Lib
void cgo_unload_lib();