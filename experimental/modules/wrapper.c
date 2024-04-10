#include <dlfcn.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "wrapper.h"

// Speichert die Geladene LIB ab
EXTERNAL_LIB global_lib;

// Lädt die Lib
const char* load_external_lib(const char* lib_path) {
    // Es wird versucht die lib zu laden
    void* lib = dlopen(lib_path, RTLD_LAZY);
    if (!lib) return "cant_open_lib";

    // Die Lib wird zwischengespeichert
    global_lib.state = 1;
    global_lib.lib = lib;

    // Der Wert wird zurückgegeben
    return "ok";
}

// Entlädt die Lib
void unload_lib() {
    if (global_lib.lib == NULL) return;
    dlclose(global_lib.lib);
}

// Initalisiert die Library
const char* initialize() {
    return "ok";
}

// Funktion, die den Callback aufruft
void callGoCallback(callback_func callback) {
    // Aufruf des übergebenen Go-Callbacks
    callback();
}
