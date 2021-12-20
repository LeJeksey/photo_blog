package db

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func GetStructAsBsonM(u interface{}) (bson.M, error) {
	uBytes, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	var updateData bson.M
	err = json.Unmarshal(uBytes, &updateData)
	return updateData, err
}
