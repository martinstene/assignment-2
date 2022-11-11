package json_coder

import (
	"assignment-2/constants"
	"assignment-2/structs"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

/*
	GetGraphqlBodyCases marshals the map it will get from the graphql response, sends a request to get this map.
	Gets the response and reads this. After reading, it will unmarshal this map into a struct Data, which will get
	cast into another struct MostRecentCases, which will be returned.
*/
func GetGraphqlBodyCases(jsonData map[string]string) structs.MostRecentCases {
	jsonValue, _ := json.Marshal(jsonData)
	// Sends a post request to the graphql to get a map in return.
	request, err := http.NewRequest(http.MethodPost, constants.Cv19Cases, bytes.NewBuffer(jsonValue))
	client := &http.Client{}
	response, err := client.Do(request)
	defer func(Body io.ReadCloser) {
		// Error handling
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	var dataStruct structs.Data
	err = json.Unmarshal(data, &dataStruct)
	// Creates a struct with the decoded information.
	mostRecentCases := structs.MostRecentCases{
		Country:    dataStruct.Data.Country.Name,
		Date:       dataStruct.Data.Country.MostRecent.Date,
		Confirmed:  dataStruct.Data.Country.MostRecent.Confirmed,
		Recovered:  dataStruct.Data.Country.MostRecent.Recovered,
		Deaths:     dataStruct.Data.Country.MostRecent.Deaths,
		GrowthRate: dataStruct.Data.Country.MostRecent.GrowthRate,
	}
	return mostRecentCases
}

// DecodePolicyInfo is a method that takes a http request and decodes the json body
func DecodePolicyInfo(httpResponse *http.Response) structs.PolicyData {
	decoder := json.NewDecoder(httpResponse.Body)
	var policies structs.PolicyData

	if err := decoder.Decode(&policies); err != nil {
		log.Print(err, http.StatusNoContent)
	}

	return policies
}

// DecodeCountryInfo is a method that takes a http request and decodes the json body
func DecodeCountryInfo(httpResponse *http.Response) structs.CountryInfo {
	decoder := json.NewDecoder(httpResponse.Body)
	var countries structs.CountryInfo

	if err := decoder.Decode(&countries); err != nil {
		log.Print(err, http.StatusNoContent)
	}

	return countries
}

// DecodePolicyCountryInfo is a method that takes a http request and decodes the json body
func DecodePolicyCountryInfo(httpResponse *http.Response) structs.PolicyCountry {
	decoder := json.NewDecoder(httpResponse.Body)
	var countries []structs.PolicyCountry

	if err := decoder.Decode(&countries); err != nil {
		log.Print(err, http.StatusNoContent)
	}
	// Returns the first cca3 code to show.
	return countries[0]
}
