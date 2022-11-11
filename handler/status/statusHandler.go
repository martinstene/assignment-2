package status

import (
	"assignment-2/constants"
	"assignment-2/handler/notifications"
	"assignment-2/json_coder"
	"assignment-2/structs"
	"fmt"
	"log"
	"net/http"
	"time"
)

var startTime = time.Now()

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		json_coder.PrettyPrint(w, Combine())
	default:
		http.Error(w, "Method not supported.", http.StatusNotImplemented)
	}
}

// Uptime Uses the time.Since() and checks how long it was since the API started
func Uptime() int {
	return int(time.Since(startTime).Seconds())
}

// GetStatusCodeCases Uses a GET call to the root link to the university API. Then return it's statuscode.
// To get a health check for the graphql I used this source to find an easy string to check:
// https://www.apollographql.com/docs/apollo-server/monitoring/health-checks/
func GetStatusCodeCases() int {
	resp, err := http.Get(constants.Cv19Cases + "?query=%7B__typename%7D")
	if err != nil {
		log.Print(err)
	}
	return resp.StatusCode
}

// GetStatusCodePolicy Uses a GET call to the root link to the country API. Then return it's statuscode.
func GetStatusCodePolicy() int {
	resp, err := http.Get(constants.Cv19PolicyRoot)
	if err != nil {
		log.Print(err)
	}
	return resp.StatusCode
}

// Uses a GET call to the root link to the country API. Then return it's statuscode.
func getStatusCountry() int {
	resp, err := http.Get(constants.CountryRoot)
	if err != nil {
		log.Print(err)
	}
	return resp.StatusCode
}

// Combine the struct into one using the methods above.
func Combine() structs.Status {
	return structs.Status{
		CasesApi:   GetStatusCodeCases(),
		PolicyApi:  GetStatusCodePolicy(),
		CountryApi: getStatusCountry(),
		WebHooks:   notifications.GetNumberOfWebhooks(),
		Version:    "v1",
		Uptime:     fmt.Sprint(Uptime()) + " seconds",
	}
}
