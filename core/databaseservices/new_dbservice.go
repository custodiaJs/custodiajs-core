package databaseservices

import (
	"sync"
	"vnh1/core/databaseservices/services"
)

// Ein neuer DbService wird erstellt
func NewDbService() *DbService {
	return &DbService{
		mutex:                &sync.Mutex{},
		databaseServiceTable: make(map[string]services.DatabaseServiceInterface),
	}
}
