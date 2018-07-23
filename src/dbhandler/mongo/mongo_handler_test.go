package mongo

import (
	"log"
	"reflect"
	"testing"
)

const (
	DbHost         = "localhost"
	DbPort         = 27018
	DbUser         = "root"
	DbPass         = "root"
	DbName         = "test_database"
	AuthDb         = "admin"
	CollectionName = "test_collection"
)

func initDbHandler() (*mongoHandler, error) {
	dbhandler := &mongoHandler{
		host:     DbHost,
		port:     DbPort,
		database: DbName,
		autdb:    AuthDb,
		username: DbUser,
		password: DbPass,
	}
	err := dbhandler.GetConnection()
	if err != nil {
		log.Printf("Fail to init db session: %s", err.Error())
	}
	return dbhandler, err
}

func TestNewMongoHandlerConnection(t *testing.T) {
	expectedHandler := &mongoHandler{
		host:     DbHost,
		port:     DbPort,
		database: DbName,
		autdb:    AuthDb,
		username: DbUser,
		password: DbPass,
	}
	dbhandler := NewMongoHandler(DbHost, DbPort, DbName, AuthDb, DbUser, DbPass)

	if !reflect.DeepEqual(expectedHandler, dbhandler) {
		t.Fatalf("NewMongoHandler fail: expected %v but got %v", expectedHandler, dbhandler)
	}
}
