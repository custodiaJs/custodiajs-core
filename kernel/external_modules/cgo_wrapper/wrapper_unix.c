// wrapper_unix.c

#if defined(__linux__) || defined(__APPLE__) || defined(__FreeBSD__) || defined(__OpenBSD__) || defined(__NetBSD__) || defined(__unix__)

#include <dlfcn.h>
#include <stdlib.h>
#include "wrapper.h"
#include "lib_bridge.h"

// Lädt eine Unix .SO Lib
STARTUP_RESULT cgo_load_external_dynamic_unix_lib(const char* lib_path) {
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

// Lädt eine Win32 .DLL Lib
STARTUP_RESULT cgo_load_external_win32_dynamic_lib(const char* lib_path) {
    STARTUP_RESULT notSupported;
    notSupported.err = "win32_not_supported_unix_host";
    return notSupported;
}

// Lädt eine MacOS .DYLIB Lib
STARTUP_RESULT cgo_load_external_macos_dynamic_lib(const char* lib_path) {
    STARTUP_RESULT notSupported;
    notSupported.err = "win32_not_supported_unix_host";
    return notSupported;
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