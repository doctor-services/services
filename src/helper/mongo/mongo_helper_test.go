package mongo

import (
	"reflect"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestIsBsonContainNonEmptyKey(t *testing.T) {
	bsonObj := bson.M{"test": "test"}
	if IsBsonMContenNonEmptyKey(bsonObj, "test") != true {
		t.Error("Must return true if value of a key is not empty")
	}
	if IsBsonMContenNonEmptyKey(bsonObj, "non_exist") == true {
		t.Error("Must return false if value of a key is not exist")
	}
}

func TestIsMapContenNonEmptyKey(t *testing.T) {
	m := map[string]interface{}{
		"rsc": 3711,
		"r":   nil,
		"gri": 1908,
		"adg": 912,
	}
	if IsMapContenNonEmptyKey(m, "rsc") != true {
		t.Error("Must return true if value of a key is not empty")
	}
	if IsMapContenNonEmptyKey(m, "non_exist") == true {
		t.Error("Must return false if value of a key is not exist")
	}
	if IsMapContenNonEmptyKey(m, "r") == true {
		t.Error("Must return false if value of a key is nil")
	}
}

func TestCreateMapFromBsonM(t *testing.T) {
	b := bson.M{
		"rsc": 3711,
		"r":   nil,
		"gri": 1908,
		"adg": 912,
		"_id": bson.NewObjectId(),
	}
	m := CreateMapFromBsonM(b)
	if len(m) != 5 {
		t.Errorf("Result map must have 5 value but got %d", len(m))
	}
	rsc, ok := m["rsc"].(int)
	if !ok {
		t.Errorf("Expected int but got %v", reflect.TypeOf(rsc))
	}
	if rsc != 3711 {
		t.Error("After converted to map, value must no be changed")
	}
}
func TestCloneStringMap(t *testing.T) {
	m := map[string]interface{}{
		"rsc": 3711,
		"r":   nil,
		"gri": 1908,
		"adg": 912,
	}
	r := CloneStringMap(m)
	if len(m) != len(r) {
		t.Error("Cloned map must have the same lenght as source")
	}

	for k := range m {
		if m[k] != r[k] {
			t.Error("Every elements in source and target must be equal")
		}
	}
}

func TestCreateObjectID(t *testing.T) {
	id := bson.NewObjectId()
	r, e := CreateObjectID(id)
	if e != nil {
		t.Error("Must not return error in valid id")
	}
	if bson.IsObjectIdHex(r.String()) {
		t.Error("Must return a valid id")
	}
	r, e = CreateObjectID("id")
	if e == nil {
		t.Error("Must return error in not valid id")
	}
}
