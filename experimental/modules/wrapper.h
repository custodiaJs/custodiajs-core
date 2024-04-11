// wrapper.h
#include "lib_bridge.h"

typedef struct {
    const char* err;
    const char* name;
    uint version;
} STARTUP_RESULT;

typedef VM_MODULE* (*LIB_LOAD)();
typedef void (*LIB_STOP)();

STARTUP_RESULT load_external_lib(const char* lib_path);
C_VM_FUNCTION_LIST get_global_functions();
C_VM_OBJECT_LIST get_global_object();
C_VM_MODULES_LIST get_modules();
void unload_lib();
