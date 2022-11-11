package stringency

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
)

// TestMain starts a main sequence for the test file
func TestMain(m *testing.M) {
	notifications.RunFirebase("../../CredentialsFirestore.json")
	defer notifications.CloseFirebase()
	policyStub := httptest.NewServer(http.HandlerFunc(stubs.PolicyHandler))
	defer policyStub.Close()
	// Changing the normal URL to the new test URL
	constants.ChangePolicyURL(policyStub.URL + "/")

	m.Run()
}

func TestIfStringencyIsCorrect(t *testing.T) {

	result := GetStringencyStruct("NOR", "2021-03-10")

	expected := structs.PolicyResponse{
		CountryCode: "NOR",
		Scope:       "2021-03-10",
		Stringency:  69.44,
		Policies:    2,
	}

	assert.Equal(t, expected, result)
}

func TestIfStringencyIsIncorrect(t *testing.T) {

	result := GetStringencyStruct("NOR", "2021-03-10")

	expected := structs.PolicyResponse{
		CountryCode: "NOR",
		Scope:       "1999-05-08",
		Stringency:  51.85,
		Policies:    20,
	}

	assert.NotEqual(t, expected, result)
}

func TestIfStringencyHasMsgUnavailable(t *testing.T) {

	result := GetStringencyStruct("NOR", "1999-05-08")

	expected := structs.PolicyResponse{
		CountryCode: "NOR",
		Scope:       "1999-05-08",
		Stringency:  -1,
		Policies:    0,
	}

	assert.Equal(t, expected, result)
}

func TestCasesStatusCode(t *testing.T) {
	casesCode := status.GetStatusCodePolicy()

	assert.Equal(t, http.StatusOK, casesCode)
}
