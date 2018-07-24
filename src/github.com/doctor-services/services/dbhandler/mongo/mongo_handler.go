package mongo

import (
	mongoHelper "github.com/doctor-services/services/helper/mongo"

	"github.com/doctor-services/services/dbhandler"

	"log"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	paingHelper "github.com/doctor-services/helpers/paging"
)

type mongoHandler struct {
	host       string
	port       int
	database   string
	autdb      string
	username   string
	password   string
	connection *mgo.Session
}

func (m *mongoHandler) createMongoSession() (*mgo.Session, error) {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{m.host + ":" + strconv.Itoa(m.port)},
		Timeout:  60 * time.Second,
		Database: m.autdb,
		Username: m.username,
		Password: m.password,
	}
	// Create a session which maintains a pool of socket connections
	// to our MongoDBhandbhandler.
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Printf("[App.db]: Error during create mongo session: %s\n", err)
		return nil, err
	}
	return mongoSession, nil
}

// GetConnection get the singleton connection object
func (m *mongoHandler) GetConnection() error {
	if m.connection == nil {
		var err error
		m.connection, err = m.createMongoSession()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *mongoHandler) IsConnecting() bool {
	return m.connection != nil
}

func (m *mongoHandler) CloseConnection() {
	if m.connection != nil {
		m.connection.Close()
		m.connection = nil
	}
}

// InvalidObjectIDError is returned when wrong object id passed
type InvalidObjectIDError struct {
	message string
}

func (e InvalidObjectIDError) Error() string {
	return e.message
}

// GetAllItems get all items with paging infor
func (m *mongoHandler) GetAllItems(dataname string, limit int, page int, orderBy string,
	sortBy string, filters map[string]interface{}) (dbhandler.PagedResults, error) {
	//field dont empty
	filters["content"] = bson.M{"$exists": true}
	filters["seenAt"] = bson.M{"$exists": true}
	filters["createdAt"] = bson.M{"$exists": true}
	//create fillter
	if mongoHelper.IsMapContenNonEmptyKey(filters, "actorid") {
		filters["actorID"], _ = strconv.Atoi(filters["actorid"].(string))
		delete(filters, "actorid")
	} else {
		filters["actorID"] = bson.M{"$exists": true}
	}

	if mongoHelper.IsMapContenNonEmptyKey(filters, "targetuserid") {
		filters["targetUserID"], _ = strconv.Atoi(filters["targetuserid"].(string))
		delete(filters, "targetuserid")
	}

	if mongoHelper.IsMapContenNonEmptyKey(filters, "targetgroupid") {
		filters["targetGroupID"], _ = strconv.Atoi(filters["targetgroupid"].(string))
		delete(filters, "targetgroupid")
	}

	if mongoHelper.IsMapContenNonEmptyKey(filters, "type") {
		filters["notifyType"], _ = strconv.Atoi(filters["type"].(string))
		delete(filters, "type")
	} else {
		filters["notifyType"] = bson.M{"$exists": true}
	}

	if mongoHelper.IsMapContenNonEmptyKey(filters, "seen") {
		filters["seen"], _ = strconv.ParseBool(filters["seen"].(string))
	}

	// Make sure connection open
	err := m.GetConnection()
	if err != nil {
		log.Printf("[App.db]: Error during create mongo session: %s\n", err)
		return dbhandler.PagedResults{}, err
	}
	workingDBSession := m.connection.Copy()
	defer workingDBSession.Close()
	c := workingDBSession.DB(m.database).C(dataname)
	// Get total items by filters
	total, err := c.Find(filters).Count()
	if err != nil {
		log.Printf("[App.db]: Error during couting items: %s\n", err)
		return dbhandler.PagedResults{}, err
	}
	pagingInfor := paingHelper.NewPaginator(total, limit, page)
	// Create sortby string
	sortString := "+" + sortBy
	if strings.ToUpper(orderBy) == "DESC" {
		sortString = "-" + sortBy
	}
	// First we need to skip previous page items
	skip := (page * limit) - limit
	q := c.Find(filters).Sort(sortString).Skip(skip)
	//q := minquery.New(workingDBSession.DB(m.database), dataname, filters).Sort(sortString).Limit(skip)
	var items []interface{}
	err = q.Limit(limit).All(&items)
	genericItems := make([]map[string]interface{}, len(items))
	for index, item := range items {
		d := item.(bson.M)
		genericItems[index] = mongoHelper.CreateMapFromBsonM(d)
	}
	return dbhandler.PagedResults{
		Total:           total,
		CurrentPage:     page,
		TotalPage:       pagingInfor.TotalPage,
		PageSize:        len(genericItems),
		NextPage:        pagingInfor.NextPage,
		PreviousPage:    pagingInfor.PreviousPage,
		HasNextPage:     pagingInfor.HasNextPage,
		HasPreviousPage: pagingInfor.HasPreviousPage,
		Items:           genericItems,
	}, nil
}

// GetAllItemsProductGroup get all items with paging infor
func (m *mongoHandler) GetAllItemsByKey(dataname string, limit int, page int, orderBy string,
	sortBy string, filters map[string]interface{}, key string) (dbhandler.PagedResults, error) {
	// Make sure connection open
	err := m.GetConnection()
	if err != nil {
		log.Printf("[App.db]: Error during create mongo session: %s\n", err)
		return dbhandler.PagedResults{}, err
	}
	workingDBSession := m.connection.Copy()
	defer workingDBSession.Close()
	//create index
	//var cursorFields  []string
	c := workingDBSession.DB(m.database).C(dataname)
	if key == "productGroup" {
		//field dont empty
		filters["content"] = bson.M{"$exists": true}
		filters["createdAt"] = bson.M{"$exists": true}
		//create fillter
		if mongoHelper.IsMapContenNonEmptyKey(filters, "actorid") {
			filters["actorID"], _ = strconv.Atoi(filters["actorid"].(string))
			delete(filters, "actorid")
		} else {
			filters["actorID"] = bson.M{"$exists": true}
		}

		if mongoHelper.IsMapContenNonEmptyKey(filters, "type") {
			filters["notifyType"], _ = strconv.Atoi(filters["type"].(string))
			delete(filters, "type")
		} else {
			filters["notifyType"] = bson.M{"$exists": true}
		}

		//cursorFields = []string{"_id", "topic", "productGroupID"}
	}
	if key == "notifyTemplate" {
		filters["template"] = bson.M{"$exists": true}
		if mongoHelper.IsMapContenNonEmptyKey(filters, "name") {
			filters["nameTemplate"], _ = strconv.Atoi(filters["name"].(string))
			delete(filters, "name")
		} else {
			filters["nameTemplate"] = bson.M{"$exists": true}
		}

		if mongoHelper.IsMapContenNonEmptyKey(filters, "creatorid") {
			filters["creatorId"], _ = strconv.Atoi(filters["creatorid"].(string))
			delete(filters, "creatorid")
		} else {
			filters["creatorId"] = bson.M{"$exists": true}
		}
		//cursorFields = []string{"nameTemplate", "template", "_id"}
	}
	if key == "productGroupTopic" {
		if mongoHelper.IsMapContenNonEmptyKey(filters, "topic") {
			filters["topic"], _ = strconv.Atoi(filters["topic"].(string))
		} else {
			filters["topic"] = bson.M{"$exists": true}
		}

		if mongoHelper.IsMapContenNonEmptyKey(filters, "productgroupid") {
			filters["productGroupID"], _ = strconv.Atoi(filters["productgroupid"].(string))
			delete(filters, "productgroupid")
		} else {
			filters["productGroupID"] = bson.M{"$exists": true}
		}
		//cursorFields = []string{"topic", "productGroupID", "_id"}
	}
	if key == "subscriabletopic" {
		//create filter
		if mongoHelper.IsMapContenNonEmptyKey(filters, "userid") {
			filters["topicUserID"], _ = strconv.Atoi(filters["userid"].(string))
			delete(filters, "userid")
		} else {
			filters["topicUserID"] = bson.M{"$exists": true}
		}

		if mongoHelper.IsMapContenNonEmptyKey(filters, "topic") {
			filters["topic"], _ = filters["topic"].(string)
		} else {
			filters["topic"] = bson.M{"$exists": true}
			filters["topicMobile"] = bson.M{"$exists": true}
		}
		//cursorFields = []string{"topicUserID", "topic", "topicMobile", "_id"}
	}
	if key == "userdevicetoken" {
		//cursorFields = []string{"userID", "webToken", "iosToken", "androidToken", "_id"}
		filters["userID"] = bson.M{"$exists": true}
	}
	// Get total items by filters
	total, err := c.Find(filters).Count()
	if err != nil {
		log.Printf("[App.db]: Error during couting items: %s\n", err)
		return dbhandler.PagedResults{}, err
	}
	pagingInfor := paingHelper.NewPaginator(total, limit, page)
	// Create sortby string
	sortString := "+" + sortBy
	if strings.ToUpper(orderBy) == "DESC" {
		sortString = "-" + sortBy
	}
	// First we need to skip previous page items
	skip := (page * limit) - limit
	//q := minquery.New(workingDBSession.DB(m.database), dataname, filters).Sort(sortString).Limit(skip)
	q := c.Find(filters).Sort(sortString).Skip(skip)
	var items []interface{}
	err = q.Limit(limit).All(&items)
	// This will move the cursort to the last item need to skip
	//skipCursor, err := q.All(&items, cursorFields...)
	// Startting from last skipped item, we get data
	//_, err = q.Cursor(skipCursor).Limit(limit).All(&items, cursorFields...)
	// Cover bson items to generic slices items
	genericItems := make([]map[string]interface{}, len(items))
	for index, item := range items {
		doc := item.(bson.M)
		genericItems[index] = mongoHelper.CreateMapFromBsonM(doc)
	}
	return dbhandler.PagedResults{
		Total:           total,
		CurrentPage:     page,
		TotalPage:       pagingInfor.TotalPage,
		PageSize:        len(genericItems),
		NextPage:        pagingInfor.NextPage,
		PreviousPage:    pagingInfor.PreviousPage,
		HasNextPage:     pagingInfor.HasNextPage,
		HasPreviousPage: pagingInfor.HasPreviousPage,
		Items:           genericItems,
	}, nil
}

// GetAllItemsNoLimit get all items no limit
func (m *mongoHandler) GetAllItemsNoLimit(dataname string, key string) ([]map[string]interface{}, error) {
	err := m.GetConnection()
	if err != nil {
		log.Printf("[App.db]: Error during create mongo session: %s\n", err)
		return nil, err
	}
	workingDBSession := m.connection.Copy()
	defer workingDBSession.Close()
	c := workingDBSession.DB(m.database).C(dataname)
	filters := map[string]interface{}{}
	if key == "userdevicetoken" {
		filters["userID"] = bson.M{"$exists": true}
	}

	var items []interface{}
	err = c.Find(filters).All(&items)
	if err != nil {
		log.Printf("[App.db]: Error during get all items: %s\n", err)
		return nil, err
	}
	genericItems := make([]map[string]interface{}, len(items))
	for index, item := range items {
		doc := item.(bson.M)
		genericItems[index] = mongoHelper.CreateMapFromBsonM(doc)
	}
	return genericItems, err
}

func (m *mongoHandler) AddNewItem(dataName string, item map[string]interface{}) (map[string]interface{}, error) {
	// Make sure not modify original map
	willInsertDoc := mongoHelper.CloneStringMap(item)
	// Make sure connection open
	err := m.GetConnection()
	if err != nil {
		log.Printf("[App.db]: Error during save %+v\n. %s\n", item, err)
		return willInsertDoc, err
	}
	// Create unique id for item
	if providedID, ok := willInsertDoc["_id"]; !ok || providedID == nil || providedID == "" {
		willInsertDoc["_id"] = bson.NewObjectId()
	}
	if _, ok := willInsertDoc["_id"].(bson.ObjectId); !ok {
		willInsertDoc["_id"], err = mongoHelper.CreateObjectID(willInsertDoc["_id"])
		if err != nil {
			return willInsertDoc, err
		}
	}
	workingDBSession := m.connection.Copy()
	defer workingDBSession.Close()
	c := workingDBSession.DB(m.database).C(dataName)
	indexes, _ := c.Indexes()
	for _, index := range indexes {
		err = c.DropIndex(index.Key...)
		if err != nil {
		}
	}
	err = c.Insert(willInsertDoc)
	if err != nil {
		return item, err
	}
	// return hexid
	returnedID, _ := willInsertDoc["_id"].(bson.ObjectId).MarshalText()
	willInsertDoc["_id"] = string(returnedID)
	return willInsertDoc, err
}

func (m *mongoHandler) RemoveItemByID(dataName string, id interface{}) error {
	// Make sure connection open
	err := m.GetConnection()
	// Make sure to use correct object id
	objectID, err := mongoHelper.CreateObjectID(id)
	if err != nil {
		log.Printf("[App.db]: Error remove item %s. %s\n", id, err)
		return err
	}
	workingDBSession := m.connection.Copy()
	defer workingDBSession.Close()
	c := workingDBSession.DB(m.database).C(dataName)
	response := c.RemoveId(objectID)
	return response
}

func (m *mongoHandler) FindItemByID(dataName string, id interface{}) (map[string]interface{}, error) {
	var data map[string]interface{}
	// Make sure connection open
	err := m.GetConnection()
	if err != nil {
		log.Printf("[App.db]: Error remove item %s. %s\n", id, err)
		return data, err
	}
	// Make sure to use correct object id
	objectID, err := mongoHelper.CreateObjectID(id)
	if err != nil {
		log.Printf("[App.db]: Error during create object id %s. %s\n", id, err)
		return data, err
	}
	workingDBSession := m.connection.Copy()
	defer workingDBSession.Close()
	c := workingDBSession.DB(m.database).C(dataName)
	var found interface{}
	err = c.FindId(objectID).One(&found)
	if err != nil {
		log.Printf("[App.db]: Error find item %s. %s\n", id, err)
		return data, err
	}
	data = mongoHelper.CreateMapFromBsonM(found.(bson.M))
	return data, nil
}

func (m *mongoHandler) UpdateByID(dataName string, id interface{}, update map[string]interface{}) error {
	// Make sure connection open
	err := m.GetConnection()
	if err != nil {
		log.Printf("[App.db]: Error during get connection for updating item %s. %s\n", id, err)
		return err
	}
	workingDBSession := m.connection.Copy()
	defer workingDBSession.Close()
	c := workingDBSession.DB(m.database).C(dataName)
	// Make sure to use correct object id
	objectID, err := mongoHelper.CreateObjectID(id)
	if err != nil {
		log.Printf("[App.db]: Error during create object id %s. %s\n", id, err)
		return err
	}
	// Not allow to update id
	willUpdateDoc := mongoHelper.CloneStringMap(update)
	delete(willUpdateDoc, "_id")
	response := c.UpdateId(objectID, willUpdateDoc)
	return response
}

func (m *mongoHandler) UpdateBy(dataName string, selector interface{}, update map[string]interface{}) error {
	// Make sure connection open
	err := m.GetConnection()
	if err != nil {
		log.Printf("[App.db]: Error during get connection for updating item %s. %s\n", selector, err)
		return err
	}
	willUpdateDoc := mongoHelper.CloneStringMap(update)
	delete(willUpdateDoc, "_id")
	workingDBSession := m.connection.Copy()
	defer workingDBSession.Close()
	c := workingDBSession.DB(m.database).C(dataName)
	_, err = c.UpdateAll(selector, bson.M{"$set": willUpdateDoc})
	return err
}

// NewMongoHandler create a instance of mongo db
func NewMongoHandler(host string, port int, database string, authdb string,
	username string, password string) dbhandler.DatabaseHandler {
	return &mongoHandler{
		host:     host,
		port:     port,
		database: database,
		autdb:    authdb,
		username: username,
		password: password,
	}
}
