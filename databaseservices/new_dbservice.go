package databaseservices

import (
	"sync"

	"github.com/CustodiaJS/custodiajs-core/databaseservices/services"
)

// Ein neuer DbService wird erstellt
func NewDbService() *DbService {
	return &DbService{
		mutex:                &sync.Mutex{},
		databaseServiceTable: make(map[string]services.DatabaseServiceInterface),
	}
}
