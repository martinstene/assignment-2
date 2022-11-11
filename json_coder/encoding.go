package json_coder

import (
	"assignment-2/util"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
PrettyPrint gotten from 02-JSON-demo
Using an interface so that no extra method is needed.
Takes in a decoded json list and reformats it so that it looks
Pretty, hence the name. Uses a responseWriter to write it out on the API.
*/
func PrettyPrint(w http.ResponseWriter, completedList interface{}) {
	util.ContentType(w)
	output, err := json.MarshalIndent(completedList, "", "  ")
	if err != nil {
		http.Error(w, "Error during pretty printing", http.StatusInternalServerError)
		return
	}
	_, err = fmt.Fprintf(w, string(output)+"\n")
	if err != nil {
		http.Error(w, "Something wrong happened while printing", http.StatusExpectationFailed)
		return
	}
}
