// wrapper.c

#include <dlfcn.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "wrapper.h"
#include "lib_bridge.h"

void* lib;                  // Speichert die LIB ab (.so / .dll)
VmModule* vm_module;        // Speichert das Verwendetbare lib Module Struct ab

// Lädt die Lib
STARTUP_RESULT cgo_load_external_lib(const char* lib_path) {
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
    vm_module = lib_load();

    // Der Wert wird zurückgegeben
    result.err = "";
    result.name = vm_module->name;
    result.version = vm_module->version;

    // Rückgabe
    return result;
}

// Gibt alle Funktionen zurück
CVmFunctionList cgo_get_global_functions() {
    // Stelle sicher, dass vm_module gültig und initialisiert ist.
    if (vm_module == NULL || vm_module->nvm_function_list == NULL) {
        // Hier solltest du entscheiden, wie du mit dieser Situation umgehen willst.
        // Zum Beispiel könntest du eine leere CVmFunctionList zurückgeben:
        CVmFunctionList empty = {0};
        return empty;
    }

    // Gibt eine Kopie der CVmFunctionList Struktur zurück
    return *(vm_module->nvm_function_list);
}

// Gibt alle Verfügbaren Module zurück
CVmModulesList cgo_get_modules() {
    // Stelle sicher, dass vm_module gültig und initialisiert ist.
    if (vm_module == NULL || vm_module->nvm_modules == NULL) {
        // Hier solltest du entscheiden, wie du mit dieser Situation umgehen willst.
        // Zum Beispiel könntest du eine leere CVmFunctionList zurückgeben:
        CVmModulesList empty = {0};
        return empty;
    }

    // Gibt eine Kopie der CVmFunctionList Struktur zurück
    return *(vm_module->nvm_modules);
}

// Gibt alle Globalen Objekte zurück
CVmObjectList cgo_get_global_object() {
    // Stelle sicher, dass vm_module gültig und initialisiert ist.
    if (vm_module == NULL || vm_module->nvm_objects == NULL) {
        // Hier solltest du entscheiden, wie du mit dieser Situation umgehen willst.
        // Zum Beispiel könntest du eine leere CVmObjectList zurückgeben:
        CVmObjectList empty = {0};
        return empty;
    }

    // Gibt eine Kopie der CVmObjectList Struktur zurück
    return *(vm_module->nvm_objects);
}

// Wir verwendet um eine Module Funktion aufzurufen
CFunctionReturnData cgo_call_function(CVmFunction* function) {
    CFunctionReturnData res = function->fptr();
    return res;
}

// Entlädt die Lib
void cgo_unload_lib() {
    if (lib) {
        dlclose(lib);
        lib = NULL;
    }

    if (vm_module) {
        free_module(vm_module);
        vm_module = NULL;
    }
}
