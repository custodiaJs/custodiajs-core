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
    UINT,
    FLOAT,
    BOOLEAN,
    TIMESTAMP,
    OBJECT,
    FUNCTION,
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
    uint uint_data;
    void* callback_data;
} CFunctionReturnData;

// Stellt eine Liste von Funktionsargumenten dar
typedef struct {
    void *array;
    size_t size;
    size_t capacity;
} CVmCallbackFunctionParameterList;

// Stellt den Funktionspointer für die Geteilte Funktion dar
typedef CFunctionReturnData (*FUNCTION_PTR)(CVmCallbackFunctionParameterList*);

// Stellt eine Funtion dar
typedef struct {
    char* name;
    FUNCTION_PTR fptr;
} CVmFunction;

// Stellt eine Funktionsliste dar
typedef struct {
    CVmFunction *array;
    size_t size;
    size_t capacity;
} CVmFunctionList;

// Stellt eine Liste von VM Objekten dar
typedef struct {
    void *CVmObject;
    size_t size;
    size_t capacity;
} CVmObjectList;

// Stellt ein VM Objekt dar
typedef struct {
    char* name;
    CVmFunctionList* nvm_function_list;
    CVmObjectList* nvm_objects;
} CVmObject;

// Stellt ein Module dar
typedef struct {
    char* name;
    CVmFunctionList* nvm_function_list;
} CVmModule;

// Stellt eine Liste von Modulen dar
typedef struct {
    CVmModule *array;
    size_t size;
    size_t capacity;
} CVmModulesList;

// Stellt eine Shared Lib dar
typedef struct {
    CVmFunctionList* nvm_function_list;
    CVmModulesList* nvm_modules;
    CVmObjectList* nvm_objects;
    char* name;
    int version;
} VmModule;

// Definiert die basisfunktion um die Lib zu laden
VmModule* new_vm_module(const char*, int);

// Array Funktionen
void free_cvm_callback_function_parameter_list(CVmCallbackFunctionParameterList*);
int add_shared_function_array(CVmFunctionList*, const char*, FUNCTION_PTR);
void init_function_parms_list(CVmCallbackFunctionParameterList*);
void init_shared_function_array(CVmFunctionList*);
int add_vm_module(CVmModulesList*, CVmModule*);
void free_vm_function_list(CVmFunctionList*);
void free_vm_modules_list(CVmModulesList*);
void init_vm_modules_list(CVmModulesList*);
void free_vm_object_list(CVmObjectList*);
void init_vm_object_list(CVmObjectList*);
void free_module(VmModule*);


// Stellt die LIB Funktionen bereit
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

// Erstellt ein neues Objekt
CVmObject* CVmObject_New();

#endif /* LIBRARY_H */
