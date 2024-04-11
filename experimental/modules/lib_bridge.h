// header

#ifndef LIBRARY_H
#define LIBRARY_H

// Stellt den Funktionspointer für die Geteilte Funktion dar
typedef int (*FUNCTION_PTR)();

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
} VM_MODULE;

// Definiert die basisfunktion um die Lib zu laden
VM_MODULE* new_vm_module_config(const char* name, int version);

// Array Funktionen
void init_vm_object_list(C_VM_OBJECT_LIST *list);
void init_vm_modules_list(C_VM_MODULES_LIST *list);
void init_shared_function_array(C_VM_FUNCTION_LIST *arr);
int add_shared_function_array(C_VM_FUNCTION_LIST *pa, const char *name, FUNCTION_PTR fptr);
int add_global_function(VM_MODULE* slib, const char* name, FUNCTION_PTR fptr);
int add_vm_module(C_VM_MODULES_LIST *list, C_VM_MODULE *module);
void free_vm_modules_list(C_VM_MODULES_LIST *list);
void free_vm_object_list(C_VM_OBJECT_LIST* list);
void free_vm_function_list(C_VM_FUNCTION_LIST *pa);
void free_module(VM_MODULE* module);

#endif /* LIBRARY_H */
