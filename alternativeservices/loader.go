package alternativeservices

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdio.h>
#include "libshared.h"
*/
import "C"

import (
	"fmt"
	"vnh1/types"
	"vnh1/utils"
)

type getHelloWord func(number int) int

type AlternativeService struct {
}

func loadAlternativeService(path string) (*AlternativeService, error) {
	handle := C.dlopen(C.CString(path), C.RTLD_LAZY)
	bar := C.dlsym(handle, C.CString("GetHelloWord"))
	fmt.Printf("bar is at %p\n", bar)

	// Aufruf der C-Funktion, um GetHelloWord aufzurufen
	C.callHelloWord(bar)

	return nil, nil
}

func LoadAllAlternativeServices() []*AlternativeService {
	fmt.Println("Loading all Alternative Services...")
	// Es wird geprüft ob der Ordner vorhanden ist
	if !utils.FolderExists(string(types.UNIX_ALTERNATIVE_SERVICES)) {
		return nil
	}

	// Es werden alle .so Dateien aufgelistet
	dir, err := utils.WalkDir(string(types.UNIX_ALTERNATIVE_SERVICES), true)
	if err != nil {
		panic(err)
	}

	fmt.Println(dir)

	// Es werden alle .so Dateien herausgefiltertet
	libfiles := []*AlternativeService{}
	for _, item := range dir {
		if item.Extension != ".so" {
			continue
		}
		result, err := loadAlternativeService(item.Path)
		if err != nil {
			panic("Internal error: " + err.Error())
		}
		libfiles = append(libfiles, result)
	}

	// Die Libfiles werden zurückgegeben
	return libfiles
}
