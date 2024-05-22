package postgresql

import (
	"database/sql"
	"fmt"
	"sync"
	"vnh1/databaseservices/services"

	_ "github.com/lib/pq"
)

type PostgreSqlDatabaseService struct {
	mutex *sync.Mutex
	db    *sql.DB
}

func (o *PostgreSqlDatabaseService) CreateNewLink() services.DbServiceLinkinterface {
	return nil
}

func NewPostgreSqlDatabaseService(host, username, password, database, alias string) (services.DatabaseServiceInterface, error) {
	// Es wird versucht den MySQL Dienst vorzubereiten
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, host, database))
	if err != nil {
		return nil, fmt.Errorf("NewPostgreSqlDatabaseService: " + err.Error())
	}

	// Es wird geprfüt ob die Verbindung vorhanden ist
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("NewPostgreSqlDatabaseService: " + err.Error())
	}

	// Das Objekt wird erstellt
	obj := &PostgreSqlDatabaseService{mutex: &sync.Mutex{}, db: db}

	// Log
	fmt.Printf("PostgresSQL Database service created, conntected to '%s'\n", host)

	// Das Objekt wird zurückgegeben
	return obj, nil
}
