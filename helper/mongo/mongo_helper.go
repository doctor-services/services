package mongo

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

// IsBsonMContenNonEmptyKey checks whether a bson object has a key with non-empty value
func IsBsonMContenNonEmptyKey(data bson.M, key string) bool {
	val, ok := data[key]
	return ok && val != nil
}

// IsMapContenNonEmptyKey check a value of a key in a map empty or not
func IsMapContenNonEmptyKey(data map[string]interface{}, key string) bool {
	val, ok := data[key]
	return ok && val != nil
}

// CreateMapFromBsonM convert bson.M to generic type
func CreateMapFromBsonM(doc bson.M) map[string]interface{} {
	var message map[string]interface{}
	message = map[string]interface{}(doc)
	// Set id to id string
	if IsBsonMContenNonEmptyKey(doc, "_id") {
		objectIDText, _ := doc["_id"].(bson.ObjectId).MarshalText()
		message["_id"] = string(objectIDText)
	}
	return message
}

// CloneStringMap copyes a string map
func CloneStringMap(source map[string]interface{}) map[string]interface{} {
	resultMap := make(map[string]interface{}, len(source))
	for key, value := range source {
		resultMap[key] = value
	}
	return resultMap
}

// CreateObjectID create a mongo object ID from a valid object ID interface
func CreateObjectID(id interface{}) (bson.ObjectId, error) {
	if _, ok := id.(bson.ObjectId); !ok {
		// Incase input is string
		stringID, ok := id.(string)
		if ok {
			if !bson.IsObjectIdHex(stringID) {
				return bson.ObjectId(""), errors.New("Wrong id format")
			}
			return bson.ObjectIdHex(stringID), nil
		}
		bytesID, ok := id.([]byte)
		if !ok {
			return bson.ObjectId(""), errors.New("Unsuported input: only support string and []byte")
		}
		// create a (may be invalid) object type of ObjectId
		var result = bson.ObjectId(bytesID)
		err := result.UnmarshalText(bytesID)
		if err != nil {
			return bson.ObjectId(""), errors.New("Wrong id format")
		}

		return result, nil
	}
	return id.(bson.ObjectId), nil
}
