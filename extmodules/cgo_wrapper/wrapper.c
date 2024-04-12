// wrapper.c

#include <dlfcn.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "wrapper.h"
#include "lib_bridge.h"

// Lädt die Lib
STARTUP_RESULT cgo_load_external_lib(const char* lib_path) {
    // Das Ergebniss wird zurückgegeben
    STARTUP_RESULT result;

    // Es wird versucht die lib zu laden
    void* lib = dlopen(lib_path, RTLD_LAZY);
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
    VmModule* vm_module = lib_load();

    // Die Lib wird Extrahiert
    CWrappedModuleLib* module = (CWrappedModuleLib*)malloc(sizeof(CWrappedModuleLib));
    module->vm_module = vm_module;
    module->lib = lib;

    // Der Wert wird zurückgegeben
    result.err = "";
    result.name = vm_module->name;
    result.version = vm_module->version;
    result.moduleLib = module;

    // Rückgabe
    return result;
}

// Gibt alle Funktionen zurück
CVmFunctionList cgo_get_global_functions(CWrappedModuleLib* vm_module) {
    // Stelle sicher, dass vm_module gültig und initialisiert ist.
    if (vm_module == NULL || vm_module->vm_module->nvm_function_list == NULL) {
        // Hier solltest du entscheiden, wie du mit dieser Situation umgehen willst.
        // Zum Beispiel könntest du eine leere CVmFunctionList zurückgeben:
        CVmFunctionList empty = {0};
        return empty;
    }

    // Gibt eine Kopie der CVmFunctionList Struktur zurück
    return *(vm_module->vm_module->nvm_function_list);
}

// Gibt alle Verfügbaren Module zurück
CVmModulesList cgo_get_modules(CWrappedModuleLib* vm_module) {
    // Stelle sicher, dass vm_module gültig und initialisiert ist.
    if (vm_module == NULL || vm_module->vm_module->nvm_modules == NULL) {
        // Hier solltest du entscheiden, wie du mit dieser Situation umgehen willst.
        // Zum Beispiel könntest du eine leere CVmFunctionList zurückgeben:
        CVmModulesList empty = {0};
        return empty;
    }

    // Gibt eine Kopie der CVmFunctionList Struktur zurück
    return *(vm_module->vm_module->nvm_modules);
}

// Gibt alle Globalen Objekte zurück
CVmObjectList cgo_get_global_object(CWrappedModuleLib* vm_module) {
    // Stelle sicher, dass vm_module gültig und initialisiert ist.
    if (vm_module == NULL || vm_module->vm_module->nvm_objects == NULL) {
        // Hier solltest du entscheiden, wie du mit dieser Situation umgehen willst.
        // Zum Beispiel könntest du eine leere CVmObjectList zurückgeben:
        CVmObjectList empty = {0};
        return empty;
    }

    // Gibt eine Kopie der CVmObjectList Struktur zurück
    return *(vm_module->vm_module->nvm_objects);
}

// Wir verwendet um eine Module Funktion aufzurufen
CFunctionReturnData cgo_call_function(CVmFunction* function) {
    CFunctionReturnData res = function->fptr();
    return res;
}

// Entlädt die Lib und gibt das vm_module frei
void cgo_unload_lib(CWrappedModuleLib* vm_module) {
    if (vm_module != NULL) {
        // Erst das Modul freigeben, wenn es existiert
        if (vm_module->vm_module != NULL) {
            // Angenommen, free_module ist die korrekte Freigabefunktion
            free_module(vm_module->vm_module);
            vm_module->vm_module = NULL;
        }

        // Dann die Bibliothek schließen, wenn sie geladen wurde
        if (vm_module->lib != NULL) {
            dlclose(vm_module->lib);
            vm_module->lib = NULL;
        }

        // Schließlich das Hauptobjekt freigeben
        free(vm_module);
    }
}

