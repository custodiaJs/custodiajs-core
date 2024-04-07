package sqlite

import (
	"database/sql"
	"fmt"
	"sync"
	"vnh1/core/databaseservices/services"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteDatabaseService struct {
	mutex *sync.Mutex
	db    *sql.DB
}

func (o *SqliteDatabaseService) CreateNewLink() services.DbServiceLinkinterface {
	return nil
}

func NewSqliteDatabaseService(sqlfName, alias string) (services.DatabaseServiceInterface, error) {
	// Es wird versucht den MySQL Dienst vorzubereiten
	db, err := sql.Open("sqlite3", sqlfName)
	if err != nil {
		return nil, fmt.Errorf("NewSqliteDatabaseService: " + err.Error())
	}

	// Es wird geprfüt ob die Verbindung vorhanden ist
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("NewSqliteDatabaseService: " + err.Error())
	}

	// Das Objekt wird erstellt
	obj := &SqliteDatabaseService{mutex: &sync.Mutex{}, db: db}

	// Log
	fmt.Printf("SQLite Database service created, used file '%s'\n", sqlfName)

	// Das Objekt wird zurückgegeben
	return obj, nil
}
