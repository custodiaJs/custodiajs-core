// wrapper.h
#include "lib_bridge.h"

typedef struct {
    const char* err;
    const char* name;
    uint version;
} STARTUP_RESULT;

typedef VmModule* (*LIB_LOAD)();
typedef void (*LIB_STOP)();

CFunctionReturnData cgo_call_function(C_VM_FUNCTION* function);
STARTUP_RESULT cgo_load_external_lib(const char* lib_path);
C_VM_FUNCTION_LIST cgo_get_global_functions();
C_VM_OBJECT_LIST cgo_get_global_object();
C_VM_MODULES_LIST cgo_get_modules();
void cgo_unload_lib();