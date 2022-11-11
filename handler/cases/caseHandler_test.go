package cases

import (
	"assignment-2/constants"
	"assignment-2/handler/notifications"
	"assignment-2/handler/status"
	"assignment-2/structs"
	"assignment-2/test/stubs"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestMain starts a main sequence for the test file
func TestMain(m *testing.M) {
	notifications.RunFirebase("../../CredentialsFirestore.json")
	defer notifications.CloseFirebase()
	casesStub := httptest.NewServer(http.HandlerFunc(stubs.CasesHandler))
	defer casesStub.Close()
	// Changing the normal URL to the new test URL
	constants.ChangeCasesURL(casesStub.URL)

	m.Run()
}

func TestIfCasesIsCorrect(t *testing.T) {

	result, err := GetCases("Norway")

	expected := structs.MostRecentCases{
		Country:    "Norway",
		Date:       "2022-04-30",
		Confirmed:  1426013,
		Recovered:  0,
		Deaths:     2932,
		GrowthRate: 0.0001430766673516579,
		TimeStamp:  time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	assert.Equal(t, expected, result, "The structs are equal.")
	assert.Nil(t, err)
}

func TestIfCasesIsIncorrect(t *testing.T) {

	result, err := GetCases("Norway")

	expected := structs.MostRecentCases{
		Country:    "Norway",
		Date:       "1999-08-05",
		Confirmed:  1426013,
		Recovered:  0,
		Deaths:     2932,
		GrowthRate: 0.0001430766673516579,
		TimeStamp:  time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	assert.NotEqualf(t, expected, result, "The structs are not equal.")
	assert.Nil(t, err)
}

func TestCasesStatusCode(t *testing.T) {
	casesCode := status.GetStatusCodeCases()

	assert.Equal(t, http.StatusOK, casesCode)
}
