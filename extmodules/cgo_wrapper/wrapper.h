// wrapper.h
#include "lib_bridge.h"

typedef struct {
    void* lib;
    VmModule* vm_module; 
} CWrappedModuleLib;

typedef struct {
    const char* err;
    const char* name;
    uint version;
    CWrappedModuleLib* moduleLib;
} STARTUP_RESULT;

typedef VmModule* (*LIB_LOAD)();
typedef void (*LIB_STOP)();

typedef struct {
    CFunctionReturnData returnData;
    char* ErrorMsg;
} CGO_FUNCTION_CALL_RETURN_STATE;

// Ruft die in bereitgestellte Funktion auf
CGO_FUNCTION_CALL_RETURN_STATE cgo_call_function(CVmFunction*, CVmCallbackFunctionParameterList*);

// Wird verwenet um die Lib zu Laden
STARTUP_RESULT cgo_load_external_dynamic_unix_lib(const char*);
STARTUP_RESULT cgo_load_external_win32_dynamic_lib(const char*);
STARTUP_RESULT cgo_load_external_macos_dynamic_lib(const char*);

// Gibt alle Verfügbaren Globalen Funktionen zurück
CVmFunctionList cgo_get_global_functions(CWrappedModuleLib*);

// Gibt alle Verfügbaren Globalen Objekte zurück
CVmObjectList cgo_get_global_object(CWrappedModuleLib*);

// Diese Funktion wird verwendet um eine neue CVmCallbackFunctionParameterList zu erstellen
CVmCallbackFunctionParameterList* cgo_new_function_parm_list();

// Entlädt die Lib
void cgo_unload_lib(CWrappedModuleLib*);
