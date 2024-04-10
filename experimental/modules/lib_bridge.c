// code

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "lib_bridge.h"

SHARED_LIB new_shared_lib_config(const char* name, int version) {
    // Alloziere Speicher für die SHARED_LIB Struktur
    SHARED_LIB result;

    result.name = strdup(name); // strdup alloziert ebenfalls Speicher
    result.version = version;

    // Initialisiere das shared_functions Array
    init_shared_function_array(&result.shared_functions);
    return result;
}

SHARED_FUNCTION new_shared_function(const char* name, FUNCTION_PTR fptr) {
    SHARED_FUNCTION result;
    result.name = strdup(name);
    result.fptr = fptr;
    return result;
}

int add_global_function(SHARED_LIB* slib, const char* name, FUNCTION_PTR fptr) {
    //add_shared_function_array(slib->shared_functions, name, fptr);
    printf("add function %s\n", name);
    return 0;
}

void init_shared_function_array(SHARED_FUNCTION_ARRAY *arr) {
    arr->array = NULL;
    arr->size = 0;
    arr->capacity = 0;
}

void add_shared_function_array(SHARED_FUNCTION_ARRAY *pa, const char *name, FUNCTION_PTR fptr) {
    if (pa->size == pa->capacity) {
        size_t newCapacity = (pa->capacity == 0) ? 1 : pa->capacity * 2;
        SHARED_FUNCTION *newArray = (SHARED_FUNCTION *)realloc(pa->array, newCapacity * sizeof(SHARED_FUNCTION));
        if (!newArray) {
            fprintf(stderr, "Speicherzuweisung fehlgeschlagen\n");
            return; // Fehlerbehandlung könnte hier verbessert werden
        }
        pa->array = newArray;
        pa->capacity = newCapacity;
    }
    pa->array[pa->size].name = strdup(name); // strdup kopiert den String und weist ihm den neuen Speicher zu
    pa->array[pa->size].fptr = fptr; // Setze den Funktionszeiger
    pa->size++;
}

void free_shared_function_array(SHARED_FUNCTION_ARRAY *pa) {
    for (size_t i = 0; i < pa->size; i++) {
        free((void*)pa->array[i].name); // Da strdup verwendet wurde, muss der Speicher freigegeben werden
    }
    free(pa->array);
    pa->array = NULL;
    pa->size = 0;
    pa->capacity = 0;
}