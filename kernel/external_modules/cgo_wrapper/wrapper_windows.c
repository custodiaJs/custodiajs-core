// wrapper_windows.c
#ifdef _WIN32

#include <windows.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "wrapper.h"
#include "lib_bridge.h"

// Lädt die Lib
STARTUP_RESULT cgo_load_external_dynamic_unix_lib(const char* lib_path) {
    // Das Ergebniss wird zurückgegeben
    STARTUP_RESULT result;

    // Handle für die geladene DLL
    HINSTANCE hDll;

    // Funktionszeiger für die Funktion in der DLL
    DLLFunctionPointer functionPtr;

    // Lade die DLL
    hDll = LoadLibraryA(lib_path);
    if (hDll == NULL) {
        fprintf(stderr, "Fehler beim Laden der DLL.\n");
        return 1;
    }

    // Hole den Funktionszeiger für die Funktion in der DLL
    functionPtr = (DLLFunctionPointer)GetProcAddress(hDll, "NameOfFunctionInDLL");
    if (functionPtr == NULL) {
        fprintf(stderr, "Fehler beim Holen des Funktionszeigers.\n");
        return 1;
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

#endif