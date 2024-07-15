#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include "wrapper.h"
#include "lib_bridge.h"

#include <signal.h>
#include <setjmp.h>


static jmp_buf env;

void signal_handler(int sig) {
    printf("Gefangen Signal %d\n", sig);
    longjmp(env, 1);
}

// Gibt alle Funktionen zurück
CVmFunctionList cgo_get_global_functions(CWrappedModuleLib* vm_module) {
    if (vm_module == NULL || vm_module->vm_module->nvm_function_list == NULL) {
        CVmFunctionList empty = {0};
        return empty;
    }

    return *(vm_module->vm_module->nvm_function_list);
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
CGO_FUNCTION_CALL_RETURN_STATE cgo_call_function(CVmFunction* function, CVmCallbackFunctionParameterList* parameters) {
    // Das Rückgabeobjekt wird erstellt
    CGO_FUNCTION_CALL_RETURN_STATE res;

    // Es wird geprüft ob ein Speicherfehler aufgetreten ist
    if (setjmp(env) == 0) {
        // Setze den Signal-Handler
        signal(SIGSEGV, signal_handler);

        // Die C-Funktion wird abgspeichert
        res.returnData = function->fptr(parameters);

        // Das Fehlerfeld wird geleert
        res.ErrorMsg = NULL;
    } else {
        // Der Fehler wird erzeugt
        res.ErrorMsg = "memory_error";
    }

    // Sollte das Parameter Objekt nicht leern sein, wird es nach dem aufrufen der Funktion zerstört
    if (parameters != NULL) free_cvm_callback_function_parameter_list(parameters);

    // Das Ergebnis wird zurückgegeben
    return res;
}

// Diese Funktion wird verwendet um eine neue CVmCallbackFunctionParameterList zu erstellen
CVmCallbackFunctionParameterList* cgo_new_function_parm_list() {
    CVmCallbackFunctionParameterList* cba;
    init_function_parms_list(cba);
    return cba;
}