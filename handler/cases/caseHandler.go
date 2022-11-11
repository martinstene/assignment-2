package cases

import (
	"assignment-2/constants"
	"assignment-2/handler/countries"
	"assignment-2/handler/notifications"
	"assignment-2/json_coder"
	"assignment-2/structs"
	"net/http"
	"path"
	"strings"
)

func CaseHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getResults(w, r)
	default:
		http.Error(w, "Method not supported.", http.StatusNotImplemented)
	}
}

/*
   GetCases gets the information about the cases
   and uses the country given to query this into graphql,
   get a response and decode it to a new struct, write it to db
   and then return the decoded struct.
*/
func GetCases(country string) (structs.MostRecentCases, error) {
	jsonData := structs.QueryGraph(country)
	response := json_coder.GetGraphqlBodyCases(jsonData)
	// Writes the response to the DB
	notifications.WriteToDB(constants.CasesDb, response)
	return response, nil
}

/*
	getResults gets a country from the url, then increments this in the DB.
	Checks if this country is already in the DB and returns accordingly.
*/
func getResults(w http.ResponseWriter, r *http.Request) {
	//splits up the name, and gets the part of the path that contains the search, also replaces space with %20
	nameToSend := strings.ReplaceAll(path.Base(r.URL.Path), " ", "%20")
	// If a cca3 code in inputted converts this to a country name instead so that the graphQL api will recognize it
	if len(nameToSend) == 3 {
		nameToSend = countries.ConvertCountryCases(nameToSend)
	}
	//Even if the user searches for a cached result it is still a call and is counted towards the webhook
	go notifications.Increment(strings.Title(strings.ToLower(nameToSend)), constants.NotficationsDb)
	if notifications.CheckIfExistsCases(w, constants.CasesDb, strings.Title(strings.ToLower(nameToSend))) {
		// turns the user-input to lowercase then uses Title to make the first letter uppercase.
		var completedQuery, err = GetCases(strings.Title(strings.ToLower(nameToSend)))
		if err != nil {
			http.Error(w, "Something went wrong when fetching the data"+err.Error(), http.StatusExpectationFailed)
			return
		}
		// prints out the completedQuery
		json_coder.PrettyPrint(w, completedQuery)
	}
}
