package mongo

import (
	"log"
	"reflect"
	"testing"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	mongoHelper "github.com/doctor-services/services/helper/mongo"
)

const (
	DbHost         = "localhost"
	DbPort         = 27017
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
func TestInitMongoConnection(t *testing.T) {
	dbhandler := &mongoHandler{
		host:     DbHost,
		port:     DbPort,
		database: DbName,
		autdb:    AuthDb,
		username: DbUser,
		password: DbPass,
	}
	defer dbhandler.CloseConnection()
	err := dbhandler.GetConnection()
	if err != nil {
		t.Fatalf("Error during create db session %v", err)
	}
	if !dbhandler.IsConnecting() {
		t.Error("Connection must be open after got connecting")
	}
}
func TestInitMongoConnectionFail(t *testing.T) {
	dbhandler := &mongoHandler{
		host:     DbHost,
		port:     DbPort,
		database: DbName,
		username: "Wronginf",
		password: DbPass,
	}
	defer dbhandler.CloseConnection()
	err := dbhandler.GetConnection()
	if err == nil {
		t.Fatalf("Connection must fail")
	}
}

func TestCloseConnection(t *testing.T) {
	dbhandler := &mongoHandler{
		host:     DbHost,
		port:     DbPort,
		database: DbName,
		autdb:    AuthDb,
		username: DbUser,
		password: DbPass,
	}
	defer dbhandler.CloseConnection()
	err := dbhandler.GetConnection()
	if dbhandler.IsConnecting() != true {
		t.Errorf("Connection must be opened after created connection but got %v", dbhandler.IsConnecting())
	}
	dbhandler.CloseConnection()
	if dbhandler.IsConnecting() == true {
		t.Error("After called close, connection must be closed")
	}
	if err != nil {
		t.Fatalf("Error during create db session %v", err)
	}
}
func TestInsertItem(t *testing.T) {
	newMessageID := bson.NewObjectId()
	message := map[string]interface{}{
		"content":      "This is test message",
		"_id":          newMessageID,
		"actorID":      1,
		"targetUserID": 1,
	}
	dbhandler := &mongoHandler{
		host:     DbHost,
		port:     DbPort,
		database: DbName,
		username: DbUser,
		password: DbPass,
	}
	defer dbhandler.CloseConnection()
	err := dbhandler.GetConnection()
	if err != nil {
		t.Fatalf("Fail to init db session: %s", err.Error())
	}
	_, err = dbhandler.AddNewItem(CollectionName, message)
	if err != nil {
		t.Fatalf("Insert item must not return error")
	}
}

func TestInsertItemDisconnect(t *testing.T) {
	newMessageID := bson.NewObjectId()
	message := map[string]interface{}{
		"content":      "This is test message",
		"_id":          newMessageID,
		"actorID":      1,
		"targetUserID": 1,
	}
	dbhandler := &mongoHandler{
		host:     DbHost,
		port:     DbPort,
		database: DbName,
		username: DbUser,
		password: DbPass,
	}
	defer dbhandler.CloseConnection()
	err := dbhandler.GetConnection()
	if err != nil {
		t.Fatalf("Fail to init db session: %s", err.Error())
	}
	dbhandler.connection.LogoutAll()
	_, err = dbhandler.AddNewItem(CollectionName, message)
	if err == nil {
		t.Fatalf("Insert item must return error")
	}
}

func TestInsertItemDontId(t *testing.T) {
	message := map[string]interface{}{
		"code":         "aaaa",
		"userId":       "289",
		"content":      "This is test message",
		"actorID":      1,
		"targetUserID": 1,
	}
	dbhandler := &mongoHandler{
		host:     DbHost,
		port:     DbPort,
		database: DbName,
		username: DbUser,
		password: DbPass,
	}
	defer dbhandler.CloseConnection()
	err := dbhandler.GetConnection()
	if err != nil {
		t.Fatalf("Fail to init db session: %s", err.Error())
	}
	_, err = dbhandler.AddNewItem(CollectionName, message)
	if err != nil {
		t.Fatalf("Insert item must not return error")
	}
}

func TestInsertAndFindById(t *testing.T) {
	newMessageID := bson.NewObjectId()
	createdAt := time.Now()
	message := map[string]interface{}{
		"content":      "This is test message",
		"_id":          newMessageID,
		"actorID":      1,
		"targetUserID": 1,
		"createdAt":    createdAt,
	}
	dbhandler := &mongoHandler{
		host:     DbHost,
		port:     DbPort,
		database: DbName,
		username: DbUser,
		password: DbPass,
	}
	defer dbhandler.CloseConnection()
	err := dbhandler.GetConnection()
	if err != nil {
		t.Fatalf("Fail to init db session: %s", err.Error())
	}
	_, err = dbhandler.AddNewItem(CollectionName, message)
	if err != nil {
		t.Fatalf("Insert item must not return error")
	}
	actualMessage, err := dbhandler.FindItemByID(CollectionName, newMessageID)
	if err != nil {
		t.Fatalf("Error during find message by ID: %s", err.Error())
	}
	if message["content"] != actualMessage["content"] {
		t.Fatalf("Found and inserted not match: expected %v but got %v", message, actualMessage)
	}
}

func TestFindAll(t *testing.T) {
	dbhandler, err := initDbHandler()
	defer dbhandler.CloseConnection()
	if err != nil {
		t.Fatalf("Fail when init db")
	}
	eachGetAll := []string{"productGroup", "notifyTemplate", "productGroupTopic", "subscriabletopic", "userdevicetoken"}
	for _, v := range eachGetAll {
		results, err := dbhandler.GetAllItemsByKey(CollectionName, 10, 1, "DESC", "createdAt", map[string]interface{}{"actorID": 1}, v)
		if err != nil {
			t.Fatalf("Error when get all items %s", err.Error())
		}
		if results.PageSize > 10 {
			t.Fatalf("Number of returned items cannot be greater than limit")
		}
	}
}

func TestFindAllFilterNoify(t *testing.T) {
	dbhandler, err := initDbHandler()
	defer dbhandler.CloseConnection()
	if err != nil {
		t.Fatalf("Fail when init db")
	}
	filters := map[string]interface{}{
		"actorid":       "adminid",
		"targetuserid":  "agencyid",
		"targetgroupid": "majorversion",
		"type":          "minorversion",
		"seen":          "patchversion",
	}
	results, err := dbhandler.GetAllItems(CollectionName, 10, 1, "DESC", "createdAt", filters)
	if err != nil {
		t.Fatalf("Error when get all items %s", err.Error())
	}
	if results.PageSize > 10 {
		t.Fatalf("Number of returned items cannot be greater than limit")
	}
}

func TestRemoveItemByID(t *testing.T) {
	dbhandler, err := initDbHandler()
	defer dbhandler.CloseConnection()
	newMessageID := bson.NewObjectId()
	createdAt := time.Now()
	message := map[string]interface{}{
		"content":      "This is test message",
		"_id":          newMessageID,
		"actorID":      1,
		"targetUserID": 1,
		"createdAt":    createdAt,
	}
	_, err = dbhandler.AddNewItem(CollectionName, message)
	if err != nil {
		t.Fatalf("Insert item must not return error")
	}
	find, err := dbhandler.FindItemByID(CollectionName, newMessageID)
	if err != nil {
		t.Fatalf("Error during find message by ID: %s", err.Error())
	}
	stringID, _ := newMessageID.MarshalText()
	dbhandler.RemoveItemByID(DbName, find["_id"].(string))
	_, err = dbhandler.FindItemByID(DbName, string(stringID))
	if err == nil {
		t.Fatalf("After deleting, find must return error")
	}
}

func TestInvalidRemoveItemByID(t *testing.T) {
	dbhandler, err := initDbHandler()
	defer dbhandler.CloseConnection()
	newMessageID := bson.NewObjectId()
	createdAt := time.Now()
	message := map[string]interface{}{
		"content":      "This is test message",
		"_id":          newMessageID,
		"actorID":      1,
		"targetUserID": 1,
		"createdAt":    createdAt,
	}
	_, err = dbhandler.AddNewItem(CollectionName, message)
	if err != nil {
		t.Fatalf("Insert item must not return error")
	}
	errRemove := dbhandler.RemoveItemByID(DbName, "fdsafas")
	if errRemove == nil {
		t.Fatalf("Remove must not return error")
	}
}

func TestInvalidFindID(t *testing.T) {
	dbhandler, err := initDbHandler()
	defer dbhandler.CloseConnection()
	newMessageID := bson.NewObjectId()
	createdAt := time.Now()
	message := map[string]interface{}{
		"content":      "This is test message",
		"_id":          newMessageID,
		"actorID":      1,
		"targetUserID": 1,
		"createdAt":    createdAt,
	}
	_, err = dbhandler.AddNewItem(CollectionName, message)
	if err != nil {
		t.Fatalf("Insert item must not return error")
	}
	find, err := dbhandler.FindItemByID(CollectionName, "")
	if err == nil {
		t.Fatalf("Find id must be return error: %s", err.Error())
		t.Fatalf("result: %+v", find)
	}
	dbhandler.RemoveItemByID(CollectionName, string(newMessageID))
}

// func TestFindByUserID(t *testing.T) {
// 	dbhandler, err := initDbHandler()
// 	defer dbhandler.CloseConnection()
// 	newMessageID := bson.NewObjectId()
// 	createdAt := time.Now()
// 	message := map[string]interface{}{
// 		"content":      "This is test message",
// 		"_id":          newMessageID,
// 		"actorID":      1,
// 		"userID":       1,
// 		"targetUserID": 1,
// 		"createdAt":    createdAt,
// 	}
// 	_, err = dbhandler.AddNewItem(CollectionName, message)
// 	if err != nil {
// 		t.Fatalf("Insert item must not return error")
// 	}
// 	find, err := dbhandler.FindItemByUserID(CollectionName, 1)
// 	if err != nil {
// 		t.Fatalf("Find id must be return error: %s", err.Error())
// 		t.Fatalf("result: %+v", find)
// 	}
// 	dbhandler.RemoveItemByID(CollectionName, find["_id"].(string))
// }

func TestUpdateBy(t *testing.T) {
	dbhandler, err := initDbHandler()
	defer dbhandler.CloseConnection()
	newMessageID := bson.NewObjectId()
	createdAt := time.Now()
	message := map[string]interface{}{
		"content":      "This is test message",
		"_id":          newMessageID,
		"actorID":      2,
		"targetUserID": 12,
		"createdAt":    createdAt,
		"seen":         false,
	}
	_, err = dbhandler.AddNewItem(CollectionName, message)
	if err != nil {
		t.Fatalf("Error during save message by ID: %s", err.Error())
	}
	insertedItem, err := dbhandler.FindItemByID(CollectionName, newMessageID)
	// t.Errorf("Inserted Item id: %v", insertedItem["_id"])
	if err != nil {
		t.Fatalf("Error during find message by ID: %s", err.Error())
	}
	insertedItem["seen"] = !(insertedItem["seen"]).(bool)
	clonedItem := mongoHelper.CloneStringMap(insertedItem)
	selector := map[string]interface{}{
		"targetUserID": "",
	}
	err = dbhandler.UpdateBy(CollectionName, selector, insertedItem)
	if err != nil {
		t.Fatalf("Update by id must not return error but got %s", err.Error())
	}
	if !reflect.DeepEqual(insertedItem, clonedItem) {
		t.Fatalf("Update must not modify original item")
	}
	dbhandler.RemoveItemByID(CollectionName, insertedItem["_id"].(string))
}

func TestUpdateByID(t *testing.T) {
	dbhandler, err := initDbHandler()
	defer dbhandler.CloseConnection()
	newMessageID := bson.NewObjectId()
	createdAt := time.Now()
	message := map[string]interface{}{
		"content":      "This is test message",
		"_id":          newMessageID,
		"actorID":      2,
		"targetUserID": 12,
		"createdAt":    createdAt,
		"seen":         false,
	}
	_, err = dbhandler.AddNewItem(CollectionName, message)
	if err != nil {
		t.Fatalf("Error during save message by ID: %s", err.Error())
	}
	insertedItem, err := dbhandler.FindItemByID(CollectionName, newMessageID)
	if err != nil {
		t.Fatalf("Error during find message by ID: %s", err.Error())
	}
	insertedItem["seen"] = !(insertedItem["seen"]).(bool)
	clonedItem := mongoHelper.CloneStringMap(insertedItem)
	err = dbhandler.UpdateByID(CollectionName, (insertedItem["_id"].(string)), insertedItem)
	if err != nil {
		t.Fatalf("Update by id must not return error but got %s", err.Error())
	}
	if !reflect.DeepEqual(insertedItem, clonedItem) {
		t.Fatalf("Update must not modify original item")
	}
	dbhandler.RemoveItemByID(CollectionName, insertedItem["_id"].(string))
}

func TestUpdateByIDDisconnect(t *testing.T) {
	dbhandler, err := initDbHandler()
	defer dbhandler.CloseConnection()
	newMessageID := bson.NewObjectId()
	createdAt := time.Now()
	message := map[string]interface{}{
		"content":      "This is test message",
		"_id":          newMessageID,
		"actorID":      2,
		"targetUserID": 12,
		"createdAt":    createdAt,
		"seen":         false,
	}
	_, err = dbhandler.AddNewItem(CollectionName, message)
	if err != nil {
		t.Fatalf("Error during save message by ID: %s", err.Error())
	}
	insertedItem, err := dbhandler.FindItemByID(CollectionName, newMessageID)
	if err != nil {
		t.Fatalf("Error during find message by ID: %s", err.Error())
	}
	dbhandler.connection.LogoutAll()
	err = dbhandler.UpdateByID(CollectionName, (insertedItem["_id"].(string)), insertedItem)
	if err == nil {
		t.Fatalf("Update by id must return error but got %s", err.Error())
	}
}

func TestInvalidUpdateByID(t *testing.T) {
	dbhandler, err := initDbHandler()
	defer dbhandler.CloseConnection()
	newMessageID := bson.NewObjectId()
	createdAt := time.Now()
	message := map[string]interface{}{
		"content":      "This is test message",
		"_id":          newMessageID,
		"actorID":      2,
		"targetUserID": 12,
		"createdAt":    createdAt,
		"seen":         false,
	}
	_, err = dbhandler.AddNewItem(CollectionName, message)
	if err != nil {
		t.Fatalf("Error during save message by ID: %s", err.Error())
	}
	insertedItem, err := dbhandler.FindItemByID(CollectionName, newMessageID)
	if err != nil {
		t.Fatalf("Error during find message by ID: %s", err.Error())
	}
	insertedItem["seen"] = !(insertedItem["seen"]).(bool)
	err = dbhandler.UpdateByID(CollectionName, "", insertedItem)
	if err == nil {
		t.Fatalf("Update by id must return error but got %s", err.Error())
	}
}

func TestInvalidObjectIDError_Error(t *testing.T) {
	type fields struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := InvalidObjectIDError{
				message: tt.fields.message,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("InvalidObjectIDError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initDbHandler(t *testing.T) {
	tests := []struct {
		name    string
		want    *mongoHandler
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := initDbHandler()
			if (err != nil) != tt.wantErr {
				t.Errorf("initDbHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initDbHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mongoHandler_AddNewItem(t *testing.T) {
	type fields struct {
		host       string
		port       int
		database   string
		autdb      string
		username   string
		password   string
		connection *mgo.Session
	}
	type args struct {
		dataName string
		item     map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mongoHandler{
				host:       tt.fields.host,
				port:       tt.fields.port,
				database:   tt.fields.database,
				autdb:      tt.fields.autdb,
				username:   tt.fields.username,
				password:   tt.fields.password,
				connection: tt.fields.connection,
			}
			got, err := m.AddNewItem(tt.args.dataName, tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("mongoHandler.AddNewItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mongoHandler.AddNewItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
