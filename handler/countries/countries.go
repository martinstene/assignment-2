package countries

import (
	"assignment-2/client"
	"assignment-2/constants"
	"assignment-2/json_coder"
	"log"
	"net/http"
	"strings"
)

func ConvertCountryCases(nameOfCountry string) string {
	/* Using a strings.Builder to create an url that meets the country
	* APIs demands and then uses that to decode a country info
	* by getting back a http request from the country api.
	 */
	var urlToSend strings.Builder
	urlToSend.WriteString(constants.Cca3Code)
	urlToSend.WriteString(nameOfCountry)
	urlToSend.WriteString(constants.Cca3Fields)

	link, err := client.GetResponseFromURL(urlToSend.String())
	if err != nil {
		log.Print("No content found", http.StatusNoContent)
		return ""
	}
	// uses the created url and gets a json response. this info is decoded after.
	country := json_coder.DecodeCountryInfo(link)
	// Returns the correct formatted word to the API
	return strings.Title(strings.ToLower(country.Name.Common))
}

func ConvertCountryPolicy(nameOfCountry string) string {
	/* Using a strings.Builder to create an url that meets the country
	* APIs demands and then uses that to decode a country info
	* by getting back a http request from the country api.
	 */
	var urlToSend strings.Builder
	urlToSend.WriteString(constants.CountryInfo)
	urlToSend.WriteString(nameOfCountry)

	link, err := client.GetResponseFromURL(urlToSend.String())
	if err != nil {
		log.Print("No content found", http.StatusNoContent)
		return ""
	}
	// uses the created url and gets a json response. this info is decoded after.
	country := json_coder.DecodePolicyCountryInfo(link)
	// Returns the correct formatted word to the API
	return strings.ToUpper(country.CCa3)
}
