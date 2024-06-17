package databaseservices

import (
	"sync"

	"github.com/CustodiaJS/custodiajs-core/databaseservices/services"
)

type DbService struct {
	mutex                *sync.Mutex
	databaseServiceTable map[string]services.DatabaseServiceInterface
}
