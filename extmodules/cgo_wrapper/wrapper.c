#include <stdlib.h>
#include "wrapper.h"
#include "lib_bridge.h"

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