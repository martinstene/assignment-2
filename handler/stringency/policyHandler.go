package stringency

import (
	"assignment-2/client"
	"assignment-2/constants"
	"assignment-2/handler/countries"
	"assignment-2/handler/notifications"
	"assignment-2/json_coder"
	"assignment-2/structs"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

func PolicyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getResultsPolicy(w, r)
	default:
		http.Error(w, "Method not supported.", http.StatusNotImplemented)
	}
}

/*
	getResultsPolicy gets the URL, checks if the scope is inputted, uses both the country and the scope
	which is inputted in the URL and uses Sprint to format this URL. Increments in the DB and checks if the entry
	already exists and deals with that accordingly. If there is no information, return an error, this is by choice.
*/
func getResultsPolicy(w http.ResponseWriter, r *http.Request) {
	// path.Base gets the last element of the url
	countryPath := path.Base(r.URL.Path)
	// If the user inputs something else then a caa3 code, convert this to a normal country name
	if len(countryPath) != 3 {
		countryPath = countries.ConvertCountryPolicy(countryPath)
	}
	// Sets the scope to be today unless specified.
	scope := time.Now().Format("2006-01-02")
	// The user is able to specify the scope and separates it with "=".
	if len(r.URL.RawQuery) != 0 {
		parseRawQuery := strings.Split(r.URL.RawQuery, "=")
		// gets the first element after the equals sign.
		scope = parseRawQuery[1]
	}
	if notifications.CheckIfExistsPolicy(w, constants.PolicyDb, strings.ToUpper(countryPath), scope) {
		finalizedStruct := GetStringencyStruct(countryPath, scope)
		// Even if the user searches for a cached result it is still a call and is counted towards the webhook
		// Makes sure that all calls gets incremented the same by converting the country-path to the country name
		if len(countryPath) == 3 {
			go notifications.Increment(countries.ConvertCountryCases(countryPath), constants.NotficationsDb)
		} else if len(countryPath) != 3 {
			go notifications.Increment(countryPath, constants.NotficationsDb)
		}
		// Writes the response to the DB
		notifications.WriteToDB(constants.PolicyDb, finalizedStruct)
		// Prints
		json_coder.PrettyPrint(w, finalizedStruct)
	}
}

// GetStringencyStruct takes in the url and gets the response and decodes it.
// if the message is data unavailable it returns a standard struct with the inputted
func GetStringencyStruct(country string, scope string) structs.PolicyResponse {
	// formats the link
	link := fmt.Sprint(constants.Cv19Policy + strings.ToUpper(country) + "/" + scope)

	url, err := client.GetResponseFromURL(link)
	if err != nil {
		log.Print("Something went wrong when fetching the data from the stringency API")
	}
	policyData := json_coder.DecodePolicyInfo(url)

	if policyData.Data.Msg == "Data unavailable" {
		return structs.PolicyResponse{
			CountryCode: strings.ToUpper(country),
			Scope:       scope,
			Stringency:  -1,
			Policies:    0,
		}
	}

	finalizedStruct := CombineStructs(policyData.Data, policyData.PolicyActions)
	return finalizedStruct
}

// countPolicies checks if the stringency type code is "NONE", which means that there are no policies.
func countPolicies(policies []structs.PolicyActions) int {
	if policies[0].PolicyTypeCode == "NONE" {
		return 0
	}
	return len(policies)
}

func stringencyEmpty(stringency structs.StringencyData) float64 {
	if stringency.StringencyActual == 0 {
		return stringency.Stringency
	}
	return stringency.StringencyActual
}

// CombineStructs Combines the data info with the actions info and returns the combined struct.
func CombineStructs(data structs.StringencyData, actions []structs.PolicyActions) structs.PolicyResponse {
	return structs.PolicyResponse{
		CountryCode: data.CountryCode,
		Scope:       data.DateValue,
		Stringency:  stringencyEmpty(data),
		Policies:    countPolicies(actions),
	}
}
