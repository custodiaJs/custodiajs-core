// wrapper.h
#include "lib_bridge.h"

typedef struct {
    const char* err;
    SHARED_LIB lib;
} STARTUP_RESULT;

typedef SHARED_LIB (*LIB_LOAD)();
typedef void (*LIB_STOP)();

STARTUP_RESULT load_external_lib(const char* lib_path);
void unload_lib();
