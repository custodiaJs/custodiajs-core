package databaseservices

import (
	"sync"
	"vnh1/databaseservices/services"
)

type DbService struct {
	mutex                *sync.Mutex
	databaseServiceTable map[string]services.DatabaseServiceInterface
}
