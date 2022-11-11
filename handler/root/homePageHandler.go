package root

import (
	"fmt"
	"log"
	"net/http"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		showInstructions(w, r)
	default:
		http.Error(w, "Method not supported.", http.StatusNotImplemented)
	}
}

// Shows the user a short guide in the "/" page.
func showInstructions(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "There are a few different endpoints described in the readme. "+
		"These short instructions will give you the ability to look it up from within the API.\n\n"+
		"The first endpoint is the cases endpoint: /corona/v1/cases/{:country_name}"+
		"\nYou can also search like this: /corona/v1/cases/{:cca3_code}\n\n"+
		"The next endpoint is the stringency where scope is optional: /corona/v1/policy/{:cca3_code}{?scope=YYYY-MM-DD}"+
		"\nYou can also search like this: /corona/v1/policy/{:country_name}{?scope=YYYY-MM-DD}\n\n"+
		"The third endpoint is the status endpoint: /corona/v1/status/\n\n"+
		"And the fourth one is the notifications endpoint which has multiple features which you can explore using software such as Postman."+
		"\nThe basic endpoint is this: /corona/v1/notifications/\n\n"+
		"Have fun!")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Something weird happened when writing", http.StatusInternalServerError)
		return
	}
}
