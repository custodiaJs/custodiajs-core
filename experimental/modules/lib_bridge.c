// code
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "lib_bridge.h"
#include <stdbool.h>

// Erstellt eine neue Lib Configugartion
VmModule* new_vm_module(const char* name, int version) {
    // Alloziere Speicher für die VmModule Struktur
    VmModule* result = malloc(sizeof(VmModule));
    if (result == NULL) {
        // Speicherreservierung fehlgeschlagen
        return NULL;
    }

    result->name = strdup(name); // strdup alloziert ebenfalls Speicher
    if (result->name == NULL) {
        // strdup fehlgeschlagen, räume auf und gebe NULL zurück
        free(result);
        return NULL;
    }

    result->version = version;

    // Initialisiere das nvm_function_list Array
    // Achte darauf, dass du Speicher für nvm_function_list allozierst oder es initialisiert, bevor du es verwendest
    result->nvm_function_list = malloc(sizeof(C_VM_FUNCTION_LIST));
    if (result->nvm_function_list == NULL) {
        // Speicherreservierung für nvm_function_list fehlgeschlagen, räume auf und gebe NULL zurück
        free(result->name);
        free(result);
        return NULL;
    }
    init_shared_function_array(result->nvm_function_list);

    // Initialisiere das nvm_objects
    result->nvm_objects = (C_VM_OBJECT_LIST*)malloc(sizeof(C_VM_OBJECT_LIST));
    if (result->nvm_objects == NULL) {
        // Speicherreservierung für nvm_objects fehlgeschlagen, räume auf und gebe NULL zurück
        free(result->nvm_function_list); // Hier sollte zusätzlich eine Freigabefunktion aufgerufen werden, falls init_shared_function_array Ressourcen zuweist
        free(result->name);
        free(result);
        return NULL;
    }
    init_vm_object_list(result->nvm_objects);

    // Initialisiere das nvm_modules
    result->nvm_modules = (C_VM_MODULES_LIST*)malloc(sizeof(C_VM_MODULES_LIST));
    if (result->nvm_modules == NULL) {
        // Speicherreservierung für nvm_modules fehlgeschlagen, räume auf und gebe NULL zurück
        // Stelle sicher, dass du alle zuvor reservierten Ressourcen freigibst
        free(result->nvm_function_list); // Zusätzliche Freigabefunktion, falls notwendig
        free(result->name);
        free(result->nvm_objects); // Vergiss nicht, auch nvm_objects zu behandeln
        free(result);
        return NULL;
    }
    init_vm_modules_list(result->nvm_modules);

    // Rückgabe
    return result;
}

// Fügt eine neue Globale Funktion hinzu
int add_global_function(VmModule* slib, const char* name, FUNCTION_PTR fptr) {
    add_shared_function_array(slib->nvm_function_list, name, fptr);
    printf("add function %s\n", name);
    return 0;
}

// Erstellt ein neues Array mit geteilten Funktionen
void init_shared_function_array(C_VM_FUNCTION_LIST *arr) {
    arr->array = NULL;
    arr->size = 0;
    arr->capacity = 0;
}

// Fügt eine Geteilte Funktion zum Sharing Array hinzu
int add_shared_function_array(C_VM_FUNCTION_LIST *pa, const char *name, FUNCTION_PTR fptr) {
    if (pa->size == pa->capacity) {
        size_t newCapacity = (pa->capacity == 0) ? 1 : pa->capacity * 2;
        C_VM_FUNCTION *newArray = (C_VM_FUNCTION *)realloc(pa->array, newCapacity * sizeof(C_VM_FUNCTION));
        if (!newArray) {
            fprintf(stderr, "Speicherzuweisung fehlgeschlagen\n");
            return 1; // Fehlerbehandlung könnte hier verbessert werden
        }
        pa->array = newArray;
        pa->capacity = newCapacity;
    }
    pa->array[pa->size].name = strdup(name); // strdup kopiert den String und weist ihm den neuen Speicher zu
    pa->array[pa->size].fptr = fptr; // Setze den Funktionszeiger
    pa->size++;
    return 0;
}

// Löscht eine List mit Geteilten Funktionen
void free_vm_function_list(C_VM_FUNCTION_LIST *pa) {
    for (size_t i = 0; i < pa->size; i++) {
        free((void*)pa->array[i].name); // Da strdup verwendet wurde, muss der Speicher freigegeben werden
    }
    free(pa->array);
    pa->array = NULL;
    pa->size = 0;
    pa->capacity = 0;
}

// Erstellt eine neue Objekt liste
void init_vm_object_list(C_VM_OBJECT_LIST *list) {
    if (!list) return; // Sicherstellen, dass der übergebene Zeiger gültig ist

    list->C_VM_OBJECT = NULL; // Zu Beginn gibt es keine Objekte
    list->size = 0;
    list->capacity = 0;
}

// Fügt ein VM Objekt der VM Objekte Liste hinzu
int add_vm_object(C_VM_OBJECT_LIST* list, C_VM_OBJECT* object) {
    if (!list || !object) return -1;

    // Erweitere die Liste, wenn nötig
    if (list->size == list->capacity) {
        size_t new_capacity = list->capacity == 0 ? 1 : list->capacity * 2;
        void** new_objects = (void**)realloc(list->C_VM_OBJECT, new_capacity * sizeof(void*));
        if (!new_objects) return -1;

        list->C_VM_OBJECT = new_objects;
        list->capacity = new_capacity;
    }

    ((C_VM_OBJECT**)list->C_VM_OBJECT)[list->size++] = object;
    return 0; // Erfolg
}

// Gibt die VM Objekte Liste frei
void free_vm_object_list(C_VM_OBJECT_LIST* list) {
    if (!list) return;

    // Gehe durch die Liste und gib jedes C_VM_OBJECT frei
    for (size_t i = 0; i < list->size; i++) {
        C_VM_OBJECT* object = ((C_VM_OBJECT**)list->C_VM_OBJECT)[i];
        free(object->name); // Name freigeben
        // Hier müsste auch nvm_function_list freigegeben werden, falls notwendig
        free(object); // Objekt freigeben
    }

    // Gib den Array von Zeigern frei
    free(list->C_VM_OBJECT);

    // Gib die Liste selbst frei
    free(list);
}

// Initalisiert eine neue Module Liste
void init_vm_modules_list(C_VM_MODULES_LIST *list) {
    if (list == NULL) return;
    list->array = NULL;
    list->size = 0;
    list->capacity = 0;
}

// Fügt ein Module hinzu
int add_vm_module(C_VM_MODULES_LIST *list, C_VM_MODULE *module) {
    if (list == NULL || module == NULL) return -1;

    // Überprüfe, ob die Liste erweitert werden muss
    if (list->size == list->capacity) {
        size_t new_capacity = list->capacity > 0 ? list->capacity * 2 : 1;
        C_VM_MODULE *new_array = (C_VM_MODULE*)realloc(list->array, new_capacity * sizeof(C_VM_MODULE));
        if (new_array == NULL) return -1; // Fehler beim Erweitern der Liste

        list->array = new_array;
        list->capacity = new_capacity;
    }

    // Füge das neue Modul hinzu und erhöhe die Größe
    list->array[list->size++] = *module;
    return 0; // Erfolg
}

// Gibt ein Module frei
void free_vm_module(C_VM_MODULE* module) {
    if (module == NULL) return;

    // Freigabe des Namens, falls vorhanden
    if (module->name != NULL) {
        free(module->name);
        module->name = NULL; // Sicherstellen, dass der Zeiger nach der Freigabe nicht mehr verwendet wird
    }

    // Freigabe der nvm_function_list, falls vorhanden
    if (module->nvm_function_list != NULL) {
        // Angenommen, es gibt eine Funktion namens free_vm_function_list, die die Liste freigibt
        free_vm_function_list(module->nvm_function_list);
        // Nun den Speicher der Struktur selbst freigeben
        free(module->nvm_function_list);
        module->nvm_function_list = NULL; // Verhindert Dangling Pointer
    }

    // Optional: Falls C_VM_MODULE dynamisch alloziert wurde, gib das Modul selbst frei
    free(module);
}

// Gibt eine Module Liste frei
void free_vm_modules_list(C_VM_MODULES_LIST *list) {
    if (list == NULL) return;

    // Gehe durch die Liste und gib jedes Modul frei
    for (size_t i = 0; i < list->size; i++) {
        free_vm_module(&list->array[i]); // Angenommen, free_vm_module ist eine Funktion, die ein C_VM_MODULE freigibt
    }

    // Gib den Array von Modulen frei
    free(list->array);

    // Setze die Liste zurück
    list->array = NULL;
    list->size = 0;
    list->capacity = 0;
}

// Zerstört eine ganze Shared LIB
void free_module(VmModule* lib) {
    if (lib == NULL) return;

    // Freigabe des `name` Feldes
    if (lib->name != NULL) {
        free(lib->name);
        lib->name = NULL;
    }

    // Freigabe der `nvm_function_list`, falls vorhanden
    if (lib->nvm_function_list != NULL) {
        // Angenommen, `free_vm_function_list` ist deine Freigabefunktion für `C_VM_FUNCTION_LIST`
        free_vm_function_list(lib->nvm_function_list);
        // Freigabe der `nvm_function_list` Struktur selbst
        free(lib->nvm_function_list);
        lib->nvm_function_list = NULL;
    }

    // Freigabe der `nvm_modules`, falls vorhanden
    if (lib->nvm_modules != NULL) {
        // Angenommen, `free_vm_modules_list` ist deine Freigabefunktion für `C_VM_MODULES_LIST`
        free_vm_modules_list(lib->nvm_modules);
        // Freigabe der `nvm_modules` Struktur selbst
        free(lib->nvm_modules);
        lib->nvm_modules = NULL;
    }

    // Freigabe der `nvm_objects`, falls vorhanden
    if (lib->nvm_objects != NULL) {
        // Angenommen, `free_vm_object_list` ist deine Freigabefunktion für `C_VM_OBJECT_LIST`
        free_vm_object_list(lib->nvm_objects);
        // Freigabe der `nvm_objects` Struktur selbst
        lib->nvm_objects = NULL;
    }

    // Zum Schluss, freigabe der `VmModule` Struktur selbst
    free(lib);
}

// Erstellt einen neuen Leeren Datensatz
CFunctionReturnData CFunctionReturnData_NewEmpty() {
    CFunctionReturnData res;
    res.type = NONE;
    res.string_data = NULL;
    res.error_data = NULL;
    res.byte_data = NULL;
    res.timestamp_data = NULL;
    res.int_data = 0;
    res.float_data = 0.0;
    res.bool_data = false;
    res.object_data = NULL;
    res.array_data = NULL;
    return res;
}

// Erstellt aus einem String, ein Rückgabewert
CFunctionReturnData CFunctionReturnData_NewString(const char* str_value) {
    CFunctionReturnData res;
    res.type = STRING;
    res.string_data = strdup(str_value);  // Kopiert str_value und weist es zu

    // Stelle sicher, dass strdup erfolgreich war
    if (res.string_data == NULL) {
        // Fehlerbehandlung, z.B. durch Setzen eines Fehlers oder Rückgabe eines speziellen Werts
        res.type = ERROR;
        res.error_data = "unkown internal c error";
    }

    res.byte_data = NULL;
    res.timestamp_data = NULL;
    res.int_data = 0;
    res.float_data = 0.0;
    res.bool_data = false;
    res.object_data = NULL;
    res.array_data = NULL;

    return res;
}

// Erstellt einen Fehler, Rückgabewert
CFunctionReturnData CFunctionReturnData_NewError(const char* error_value) {
    CFunctionReturnData res;
    res.type = ERROR;
    res.error_data = strdup(error_value);  // Kopiert error_value und weist es zu

    // Stelle sicher, dass strdup erfolgreich war
    if (res.error_data == NULL) {
        // Fehlerbehandlung, z.B. durch Setzen eines Fehlers oder Rückgabe eines speziellen Werts
        res.type = ERROR;
        res.error_data = "unkown internal c error";
    }

    res.string_data = NULL;
    res.byte_data = NULL;
    res.timestamp_data = NULL;
    res.int_data = 0;
    res.float_data = 0.0;
    res.bool_data = false;
    res.object_data = NULL;
    res.array_data = NULL;

    return res;
}

// Erstellt einen neuen Byte Datensatz
CFunctionReturnData CFunctionReturnData_NewByteData(const char* bytes_value) {
    CFunctionReturnData res;
    res.type = BYTES;
    res.byte_data = strdup(bytes_value);

    // Stelle sicher, dass strdup erfolgreich war
    if (res.byte_data == NULL) {
        // Fehlerbehandlung, z.B. durch Setzen eines Fehlers oder Rückgabe eines speziellen Werts
        res.type = ERROR;
        res.error_data = "unkown internal c error";
    }

    res.string_data = NULL;
    res.timestamp_data = NULL;
    res.int_data = 0;
    res.float_data = 0.0;
    res.bool_data = false;
    res.object_data = NULL;
    res.array_data = NULL;

    return res;
}

// Erstellt einen neuen Integer Datensatz
CFunctionReturnData CFunctionReturnData_NewInt(int int_value) {
    CFunctionReturnData res;
    res.type = INT;
    res.int_data = int_value;
    res.string_data = NULL;
    res.error_data = NULL;
    res.byte_data = NULL;
    res.timestamp_data = NULL;
    res.float_data = 0.0;
    res.bool_data = false;
    res.object_data = NULL;
    res.array_data = NULL;
    return res;
}

// Erstellt einen neuen Float Datensatz
CFunctionReturnData CFunctionReturnData_NewFloat(float float_value) {
    CFunctionReturnData res;
    res.type = FLOAT;
    res.float_data = float_value;
    res.string_data = NULL;
    res.error_data = NULL;
    res.byte_data = NULL;
    res.timestamp_data = NULL;
    res.int_data = 0;
    res.bool_data = false;
    res.object_data = NULL;
    res.array_data = NULL;
    return res;
}

// Erstellt einen neuen Bool Datensatz
CFunctionReturnData CFunctionReturnData_NewBool(bool bool_value) {
    CFunctionReturnData res;
    res.type = BOOLEAN;
    res.bool_data = bool_value;
    res.string_data = NULL;
    res.error_data = NULL;
    res.byte_data = NULL;
    res.timestamp_data = NULL;
    res.int_data = 0;
    res.float_data = 0.0;
    res.object_data = NULL;
    res.array_data = NULL;
    return res;
}

// Erstellt einen neuen Timestamp
CFunctionReturnData CFunctionReturnData_NewTimestamp(const char* timesptamp_value) {
    CFunctionReturnData res;
    res.type = TIMESTAMP;
    res.timestamp_data = strdup(timesptamp_value);

    // Stelle sicher, dass strdup erfolgreich war
    if (res.timestamp_data == NULL) {
        // Fehlerbehandlung, z.B. durch Setzen eines Fehlers oder Rückgabe eines speziellen Werts
        res.type = ERROR;
        res.error_data = "unkown internal c error";
    }

    res.string_data = NULL;
    res.byte_data = NULL;
    res.int_data = 0;
    res.float_data = 0.0;
    res.bool_data = false;
    res.object_data = NULL;
    res.array_data = NULL;

    return res;
}

// Erstellt ein neues Leeres Objekt
CFunctionReturnData CFunctionReturnData_NewObject() {
    CFunctionReturnData res;
    res.type = OBJECT;
    res.string_data = NULL;
    res.error_data = NULL;
    res.byte_data = NULL;
    res.timestamp_data = NULL;
    res.int_data = 0;
    res.float_data = 0.0;
    res.bool_data = false;
    res.array_data = NULL;
    return res;
}

// Erstellt ein neues Leers Array
CFunctionReturnData CFunctionReturnData_NewArray() {
    CFunctionReturnData res;
    res.type = ARRAY;
    res.string_data = NULL;
    res.error_data = NULL;
    res.byte_data = NULL;
    res.timestamp_data = NULL;
    res.int_data = 0;
    res.float_data = 0.0;
    res.bool_data = false;
    res.object_data = NULL;
    return res;
}