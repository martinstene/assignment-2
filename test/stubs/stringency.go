package stubs

import (
	"assignment-2/structs"
	"assignment-2/util"
	"encoding/json"
	"net/http"
	"strings"
)

// PolicyHandler returns the json output mimicking the stringency api.
func PolicyHandler(w http.ResponseWriter, r *http.Request) {
	jsonData := structs.PolicyData{
		PolicyActions: []structs.PolicyActions{
			{
				PolicyTypeCode:          "C1",
				PolicyValueDisplayField: "School closing",
			},
			{
				PolicyTypeCode:          "C2",
				PolicyValueDisplayField: "Workplace closing",
			},
		}, Data: structs.StringencyData{
			Deaths:           632,
			Stringency:       69.44,
			CountryCode:      "NOR",
			StringencyActual: 69.44,
			Confirmed:        77169,
			DateValue:        "2021-03-10",
		},
	}

	if strings.Contains(r.URL.Path, "1999-05-08") {
		jsonData = structs.PolicyData{
			PolicyActions: []structs.PolicyActions{
				{
					PolicyTypeCode:          "C1",
					PolicyValueDisplayField: "School closing",
				},
				{
					PolicyTypeCode:          "C2",
					PolicyValueDisplayField: "Workplace closing",
				},
			},
			Data: structs.StringencyData{
				Msg: "Data unavailable",
			},
		}
	}

	util.ContentType(w)
	encoder := json.NewEncoder(w)

	err := encoder.Encode(jsonData)
	if err != nil {
		return
	}
}
