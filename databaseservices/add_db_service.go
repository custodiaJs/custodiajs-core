package databaseservices

import (
	"fmt"
	"strings"

	"github.com/CustodiaJS/custodiajs-core/databaseservices/services"
	"github.com/CustodiaJS/custodiajs-core/databaseservices/services/mongodb"
	"github.com/CustodiaJS/custodiajs-core/databaseservices/services/mysql"
	"github.com/CustodiaJS/custodiajs-core/databaseservices/services/postgresql"
	"github.com/CustodiaJS/custodiajs-core/databaseservices/services/sqlite"
	"github.com/CustodiaJS/custodiajs-core/vmdb"
)

// Es wird ein neuer Datenbankdienst hinzugefügt
func (o *DbService) AddDatabaseService(dmdb *vmdb.VMEntryBaseData) error {
	// Der Mutex wird verwendet
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// Es wird geprüft ob der Datenbank Dienst bereits hinzugefügt wurde
	if _, found := o.databaseServiceTable[strings.ToLower(string(dmdb.GetDatabaseFingerprint()))]; found {
		return fmt.Errorf("DbService->AddDatabaseService: database service always in table")
	}

	// Die Hostadresse wird estellt
	hostAdr := fmt.Sprintf("%s:%d", dmdb.Host, dmdb.Port)

	// Es wird ermittelt ob es sich um ein zulässigen Datenbanktypen handelt
	var dbservice services.DatabaseServiceInterface
	var err error
	switch dmdb.Type {
	case "mysql":
		// Es wird versucht den MySQL Dienst zu erstellen
		dbservice, err = mysql.NewMySqlDatabaseService(hostAdr, dmdb.Username, dmdb.Password, dmdb.Database, dmdb.Alias)
	case "mongodb":
		// Es wird versucht den MongoDB Dienst zu erstellen
		dbservice, err = mongodb.NewMongoDbDatabaseService(hostAdr, dmdb.Username, dmdb.Password, dmdb.Database, dmdb.Alias)
	case "postgresql":
		// der PostgresSql Dienst wird ersellt
		dbservice, err = postgresql.NewPostgreSqlDatabaseService(dmdb.Host, dmdb.Username, dmdb.Password, dmdb.Database, dmdb.Alias)
	case "sqlite":
		// der PostgresSql Dienst wird ersellt
		dbservice, err = sqlite.NewSqliteDatabaseService(dmdb.Host, dmdb.Alias)
	default:
		return fmt.Errorf(fmt.Sprintf("DbService->AddDatabaseService: unkown database type '%s'", dmdb.Type))
	}

	// Es wird geprüft ob ein Fehler aufgetreten ist
	if err != nil {
		return fmt.Errorf("DbService->AddDatabaseService: " + err.Error())
	}

	// Der Database Dienst wird zwischengeseichert
	o.databaseServiceTable[strings.ToLower(string(dmdb.GetDatabaseFingerprint()))] = dbservice

	// Log
	fmt.Printf("Database service added: %s, %s, %s\n", dmdb.Type, dmdb.Alias, strings.ToUpper(string(dmdb.GetDatabaseFingerprint())))

	// Der Vorgang wurde ohne Fehler durchgeführt
	return nil
}

// Es wird ein Link für einen Spezifischen Dienst erstellt
func (o *DbService) GetDBServiceLink(fingerprint vmdb.DatabaseFingerprint) (services.DbServiceLinkinterface, error) {
	// Der Mutex wird verwendet
	o.mutex.Lock()
	defer o.mutex.Unlock()

	// Es wird versucht den Datenbank Dienst zu ermitteln
	dbService, found := o.databaseServiceTable[strings.ToLower(string(fingerprint))]
	if !found {
		return nil, fmt.Errorf("DbService->GetDBServiceLink: not found " + strings.ToUpper(string(fingerprint)))
	}

	// Es wird ein neuer VM Link für diesen Dienst erstellt
	dbServiceLink := dbService.CreateNewLink()

	// Der Link wird zwischengespeichert
	o.databaseServiceTable[strings.ToLower(string(fingerprint))] = dbService

	// Der Link wird zurückgegeben
	return dbServiceLink, nil
}
