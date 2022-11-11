package structs

import "time"

// QueryGraph and Cases endpoint structs
func QueryGraph(countryToBeSearched string) map[string]string {
	jsonData := map[string]string{
		"query": ` 
			{
				country(name:"` + countryToBeSearched + `") {
				name
					mostRecent {
						date(format: "yyyy-MM-dd"),
						confirmed,
						recovered,
						deaths,
						growthRate
					}
				}
			}
		`,
	}
	return jsonData
}

type Data struct {
	Data Country `json:"data"`
}

type Country struct {
	Country InformationCases `json:"country"`
}

type InformationCases struct {
	Name       string          `json:"name"`
	MostRecent MostRecentCases `json:"mostRecent"`
}

type MostRecentCases struct {
	Country    string    `json:"country"`
	Date       string    `json:"date"`
	Confirmed  int       `json:"confirmed"`
	Recovered  int       `json:"recovered"`
	Deaths     int       `json:"deaths"`
	GrowthRate float64   `json:"growthRate"`
	TimeStamp  time.Time `firestore:"time" json:"-"`
}

// The Status endpoint struct
type Status struct {
	CasesApi   int    `json:"cases_api"`
	PolicyApi  int    `json:"policy_api"`
	CountryApi int    `json:"country_api"`
	WebHooks   int    `json:"web_hooks"`
	Version    string `json:"version"`
	Uptime     string `json:"uptime"`
}

// The stringency endpoint structs
type PolicyData struct {
	PolicyActions []PolicyActions `json:"policyActions"`
	Data          StringencyData  `json:"stringencyData"`
}

type PolicyActions struct {
	PolicyTypeCode          string `json:"policy_type_code"`
	PolicyValueDisplayField string `json:"policy_value_display_field,omitempty"`
}

type StringencyData struct {
	DateValue        string  `json:"date_value"`
	CountryCode      string  `json:"country_code"`
	Confirmed        int     `json:"confirmed"`
	Deaths           int     `json:"deaths"`
	StringencyActual float64 `json:"stringency_actual"`
	Stringency       float64 `json:"stringency"`
	Msg              string  `json:"msg"`
}

type PolicyResponse struct {
	CountryCode string    `json:"country_code"`
	Scope       string    `json:"scope"`
	Stringency  float64   `json:"stringency"`
	Policies    int       `json:"policies"`
	TimeStamp   time.Time `firestore:"time" json:"-"`
}

// The Notification endpoint structs including webhook structs.
type Notification struct {
	Url          string `json:"url"`
	Country      string `json:"country"`
	Calls        int    `json:"calls"`
	CurrentCalls int    `json:"CurrentCalls"`
	WebHookId    string `json:"WebHookId"`
}

type CreatedWebHookResponse struct {
	WebHookId string `json:"web_hook_id"`
}

type WebHookInvocation struct {
	WebHookId string `json:"web_hook_id"`
	Country   string `json:"country"`
	Calls     int    `json:"calls"`
}

// CountryInfo and CountryCca3 and PolicyCountry Structs for the converting from cca3 to country name and from name to cca3
type CountryInfo struct {
	Name CountryCca3 `json:"name"`
}

type CountryCca3 struct {
	Common string `json:"common"`
}

type PolicyCountry struct {
	CCa3 string `json:"cca3"`
}
