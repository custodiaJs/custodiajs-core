// wrapper.h

// Sterllt den LIB Struct dar
typedef struct EXTERNAL_LIB {
    uint state;
    void* lib;
} EXTERNAL_LIB;

typedef void (*callback_func)();

const char* load_external_lib(const char* lib_path);
void callGoCallback(callback_func callback);
extern void myGoCallback();
const char* initialize();
void unload_lib();