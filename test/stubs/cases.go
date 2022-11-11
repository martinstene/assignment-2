package stubs

import (
	"assignment-2/structs"
	"assignment-2/util"
	"encoding/json"
	"net/http"
)

// CasesHandler returns the json output mimicking the cases api.
func CasesHandler(w http.ResponseWriter, r *http.Request) {
	graphQlData := structs.Data{
		Data: structs.Country{
			Country: structs.InformationCases{
				Name: "Norway",
				MostRecent: structs.MostRecentCases{
					Country:    "Norway",
					Date:       "2022-04-30",
					Confirmed:  1426013,
					Deaths:     2932,
					GrowthRate: 0.0001430766673516579,
					Recovered:  0,
				},
			},
		},
	}

	util.ContentType(w)
	encoder := json.NewEncoder(w)

	err := encoder.Encode(graphQlData)
	if err != nil {
		return
	}
}
