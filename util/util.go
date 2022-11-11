package util

import (
	"assignment-2/structs"
	"encoding/json"
	"net/http"
)

// ContentType is used for creating a json content type
func ContentType(w http.ResponseWriter) {
	w.Header().Add("content-type", "application/json")
}

// ConvertMapToJson gets a map from the database and marshals it and unmarshal' it into a struct
// and returns this new struct with the data.
func ConvertMapToJson(m map[string]interface{}) structs.Notification {
	notification := structs.Notification{}
	dataFromDb, err := json.Marshal(m)
	if err != nil {
		return structs.Notification{}
	}

	err = json.Unmarshal(dataFromDb, &notification)
	if err != nil {
		return structs.Notification{}
	}
	return notification
}
