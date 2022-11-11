package status

import (
	"assignment-2/handler/notifications"
	"assignment-2/structs"
	"assignment-2/test/stubs"
	"fmt"
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

	policyStub := httptest.NewServer(http.HandlerFunc(stubs.PolicyHandler))
	defer policyStub.Close()

	notificationStub := httptest.NewServer(http.HandlerFunc(notifications.NotificationHandler))
	defer notificationStub.Close()

	m.Run()
}

func TestUptime(t *testing.T) {
	time.Sleep(2 * time.Second)
	assert.Equal(t, 2, Uptime())
}

func TestNumberOfWebhooks(t *testing.T) {
	numberOfHooks := notifications.GetNumberOfWebhooks()

	assert.Equal(t, 2, numberOfHooks)
}

func TestStatusHandler(t *testing.T) {
	result := Combine()

	expected := structs.Status{
		CasesApi:  http.StatusOK,
		PolicyApi: http.StatusOK,
		WebHooks:  notifications.GetNumberOfWebhooks(),
		Version:   "v1",
		Uptime:    fmt.Sprint(Uptime()) + " seconds",
	}

	assert.Equal(t, expected, result)
}
