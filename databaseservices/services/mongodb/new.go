package mongodb

import (
	"context"
	"fmt"
	"sync"
	"time"
	"vnh1/databaseservices/services"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbDatabaseService struct {
	mutex *sync.Mutex
	db    *mongo.Client
	host  string
}

func (o *MongoDbDatabaseService) CreateNewLink() services.DbServiceLinkinterface {
	return nil
}

func (o *MongoDbDatabaseService) connectToMongo(errorChan chan error) {
	// Ein Kontext mit Timeout von 10 Sekunden erstellen
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Verbindung zur MongoDB-Datenbank herstellen
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", o.host)))
	if err != nil {
		errorChan <- fmt.Errorf("MongoDbDatabaseService->connectToMongo: " + err.Error())
		return
	}

	// Es wird geprüft ob die Verbindung aufgebaut ist
	if err := client.Ping(ctx, nil); err != nil {
		errorChan <- fmt.Errorf("MongoDbDatabaseService->connectToMongo: " + err.Error())
		return
	}

	// Das Clientobjekt wird zwischengespeichert
	o.db = client

	// Warten, bis der Kontext beendet wird
	<-ctx.Done()
}

func NewMongoDbDatabaseService(host, username, password, database, alias string) (services.DatabaseServiceInterface, error) {
	// das MongoDB Objekt wird erstellt
	errorChan := make(chan error)
	mdbService := &MongoDbDatabaseService{mutex: &sync.Mutex{}, db: nil, host: host}

	// Eine Go-Routine starten, um die Verbindung aufzubauen
	go func() {
		mdbService.connectToMongo(errorChan)
		close(errorChan)
	}()

	// Es wird auf den Chan gewartet
	if err := <-errorChan; err != nil {
		return nil, fmt.Errorf("NewMongoDbDatabaseService: " + err.Error())
	}

	// Log
	fmt.Printf("MongoDB Database service created, conntected to '%s'\n", host)

	// Das Objekt wird zurückgegeben
	return mdbService, nil
}
