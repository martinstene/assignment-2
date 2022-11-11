package constants

const (
	// Root Cases Policy Status Notifications are endpoints constants.
	Root          = "/"
	Cases         = "/corona/v1/cases/"
	Policy        = "/corona/v1/policy/"
	Status        = "/corona/v1/status/"
	Notifications = "/corona/v1/notifications/"

	Cv19PolicyRoot = "https://covidtrackerapi.bsg.ox.ac.uk/api"

	// Cca3Code are service APIs for CCA3 conversion
	Cca3Code    = "https://restcountries.com/v3.1/alpha/"
	Cca3Fields  = "?fields=name"
	CountryInfo = "https://restcountries.com/v3.1/name/"
	CountryRoot = "https://restcountries.com"

	// Collection name in Firestore
	NotficationsDb = "notifications"
	CasesDb        = "cases"
	PolicyDb       = "stringency"
)

var (
	// Cv19Cases Cv19PolicyRoot Cv19Policy are services APIs
	Cv19Cases  = "https://covid19-graphql.now.sh"
	Cv19Policy = "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/"
)

// ChangeCasesURL ChangePolicyURL are used for stubbing
func ChangeCasesURL(url string) {
	Cv19Cases = url
}

func ChangePolicyURL(url string) {
	Cv19Policy = url
}
