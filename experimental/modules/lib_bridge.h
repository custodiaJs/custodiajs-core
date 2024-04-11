// header

#ifndef LIBRARY_H
#define LIBRARY_H

#include <stdbool.h>

/*
    Stellt die Basisfunktionen dar
*/

// Stellt die VM Datentypen bereit
typedef enum {
    NONE,
    STRING,
    ERROR,
    BYTES,
    INT,
    FLOAT,
    BOOLEAN,
    TIMESTAMP,
    OBJECT,
    ARRAY
} CVmDatatype;

// Stellt die Returnwerte eines Funktionsaufrufes dar
typedef struct {
    CVmDatatype type;
    char* string_data;
    char* error_data;
    char* byte_data;
    char* timestamp_data;
    int int_data;
    float float_data;
    bool bool_data;
    void* object_data;
    void* array_data;
} CFunctionReturnData;

// Stellt den Funktionspointer für die Geteilte Funktion dar
typedef CFunctionReturnData (*FUNCTION_PTR)();

// Stellt eine Funtion dar
typedef struct {
    char* name;
    FUNCTION_PTR fptr;
} C_VM_FUNCTION;

// Stellt eine Funktionsliste dar
typedef struct {
    C_VM_FUNCTION *array;
    size_t size;
    size_t capacity;
} C_VM_FUNCTION_LIST;

// Stellt eine Liste von VM Objekten dar
typedef struct {
    void *C_VM_OBJECT;
    size_t size;
    size_t capacity;
} C_VM_OBJECT_LIST;

// Stellt ein VM Objekt dar
typedef struct {
    char* name;
    C_VM_FUNCTION_LIST* nvm_function_list;
    C_VM_OBJECT_LIST* nvm_objects;
} C_VM_OBJECT;

// Stellt ein Module dar
typedef struct {
    char* name;
    C_VM_FUNCTION_LIST* nvm_function_list;
} C_VM_MODULE;

// Stellt eine Liste von Modulen dar
typedef struct {
    C_VM_MODULE *array;
    size_t size;
    size_t capacity;
} C_VM_MODULES_LIST;

// Stellt eine Shared Lib dar
typedef struct {
    C_VM_FUNCTION_LIST* nvm_function_list;
    C_VM_MODULES_LIST* nvm_modules;
    C_VM_OBJECT_LIST* nvm_objects;
    char* name;
    int version;
} VmModule;

// Definiert die basisfunktion um die Lib zu laden
VmModule* new_vm_module(const char* name, int version);

// Array Funktionen
void init_vm_object_list(C_VM_OBJECT_LIST *list);
void init_vm_modules_list(C_VM_MODULES_LIST *list);
void init_shared_function_array(C_VM_FUNCTION_LIST *arr);
int add_shared_function_array(C_VM_FUNCTION_LIST *pa, const char *name, FUNCTION_PTR fptr);
int add_vm_module(C_VM_MODULES_LIST *list, C_VM_MODULE *module);
void free_vm_modules_list(C_VM_MODULES_LIST *list);
void free_vm_object_list(C_VM_OBJECT_LIST* list);
void free_vm_function_list(C_VM_FUNCTION_LIST *pa);
void free_module(VmModule* module);

/*
    Stellt die LIB Funktionen bereit
*/
int add_global_function(VmModule* slib, const char* name, FUNCTION_PTR fptr);

// Erstellt einen neuen Leeren Rückgabe wert
CFunctionReturnData CFunctionReturnData_NewEmpty();

// Erstellt aus einem String, ein Rückgabewert
CFunctionReturnData CFunctionReturnData_NewString(const char*);

// Erstellt einen Fehler, Rückgabewert
CFunctionReturnData CFunctionReturnData_NewError(const char*);

// Erstellt einen neuen Byte Datensatz
CFunctionReturnData CFunctionReturnData_NewByteData(const char*);

// Erstellt einen neuen Integer Datensatz
CFunctionReturnData CFunctionReturnData_NewInt(int);

// Erstellt einen neuen Float Datensatz
CFunctionReturnData CFunctionReturnData_NewFloat(float);

// Erstellt einen neuen Bool Datensatz
CFunctionReturnData CFunctionReturnData_NewBool(bool);

// Erstellt einen neuen Timestamp
CFunctionReturnData CFunctionReturnData_NewTimestamp(const char*);

// Erstellt ein neues Leeres Objekt
CFunctionReturnData CFunctionReturnData_NewObject();

// Erstellt ein neues Leers Array
CFunctionReturnData CFunctionReturnData_NewArray();

#endif /* LIBRARY_H */
