package mysql

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/CustodiaJS/custodiajs-core/databaseservices/services"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlDatabaseService struct {
	mutex *sync.Mutex
	db    *sql.DB
}

func (o *MysqlDatabaseService) CreateNewLink() services.DbServiceLinkinterface {
	return nil
}

func NewMySqlDatabaseService(host, username, password, database, alias string) (services.DatabaseServiceInterface, error) {
	// Es wird versucht den MySQL Dienst vorzubereiten
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, database))
	if err != nil {
		return nil, fmt.Errorf("NewMySqlDatabaseService: " + err.Error())
	}

	// Es wird geprfüt ob die Verbindung vorhanden ist
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("NewMySqlDatabaseService: " + err.Error())
	}

	// Das Objekt wird erstellt
	obj := &MysqlDatabaseService{mutex: &sync.Mutex{}, db: db}

	// Log
	fmt.Printf("MySQL Database service created, conntected to '%s'\n", host)

	// Das Objekt wird zurückgegeben
	return obj, nil
}
