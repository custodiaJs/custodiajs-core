// wrapper.c

#include <dlfcn.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "wrapper.h"
#include "lib_bridge.h"

// Speichert die Geladene LIB ab
void* lib;

// Lädt die Lib
STARTUP_RESULT load_external_lib(const char* lib_path) {
    // Das Ergebniss wird zurückgegeben
    STARTUP_RESULT result;

    // Es wird versucht die lib zu laden
    lib = dlopen(lib_path, RTLD_LAZY);
    if (!lib) {
        result.err = "cant_open_lib";
        return result;
    }

    // Die Startup Funktion wird geladen
    LIB_LOAD lib_load = (LIB_LOAD) dlsym(lib, "lib_load");
    if (!lib_load) {
        dlclose(lib);
        result.err = "cant_call_function";
        return result;
    }

    // Die Shutdown Funktion wird geladen
    LIB_STOP lib_stop = (LIB_STOP) dlsym(lib, "lib_stop");
    if (!lib_stop) {
        dlclose(lib);
        result.err = "cant_call_function";
        return result;
    }

    // Die Lib wird geladen
    SHARED_LIB slib = lib_load();

    // Der Wert wird zurückgegeben
    result.err = "";
    result.lib = slib;
    return result;
}

// Entlädt die Lib
void unload_lib() {
    if (lib == NULL) return;
}