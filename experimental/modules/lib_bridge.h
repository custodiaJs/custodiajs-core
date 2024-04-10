// header

#ifndef LIBRARY_H
#define LIBRARY_H

// Stellt den Funktionspointer für die Geteilte Funktion dar
typedef int (*FUNCTION_PTR)();

// Stellt eine Funtion dar
typedef struct {
    const char* name;
    FUNCTION_PTR fptr;
} SHARED_FUNCTION;

typedef struct {
    SHARED_FUNCTION *array;
    size_t size;
    size_t capacity;
} SHARED_FUNCTION_ARRAY;

// Stellt eine Shared Lib dar
typedef struct {
    SHARED_FUNCTION_ARRAY* shared_functions;
    const char* name;
    int version;
} SHARED_LIB;

// Definiert die basisfunktion um die Lib zu laden
SHARED_LIB new_shared_lib_config(const char* name, int version);

// Definiert die Funktion zum Hinzufügen neuer VM Funktionen
int add_global_function(SHARED_LIB* slib, const char* name, FUNCTION_PTR fptr);

// Array Funktionen
void init_shared_function_array(SHARED_FUNCTION_ARRAY *arr);
void add_shared_function_array(SHARED_FUNCTION_ARRAY *pa, const char *name, FUNCTION_PTR fptr);

#endif /* LIBRARY_H */
