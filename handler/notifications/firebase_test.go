package notifications

import (
	"assignment-2/constants"
	"assignment-2/structs"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestMain starts a main sequence for the test file
func TestMain(m *testing.M) {
	RunFirebase("../../CredentialsFirestore.json")
	defer CloseFirebase()

	m.Run()
}

func TestAddWebhook(t *testing.T) {
	before := GetNumberOfWebhooks()

	jsonData := structs.Notification{
		Url:       "localhost:8080/client",
		Country:   "Switzerland",
		Calls:     4,
		WebHookId: "gerwuyihjofkøwj2324",
	}

	jsonMarshal, err := json.Marshal(jsonData)
	if err != nil {
		return
	}

	request, err := http.NewRequest(http.MethodPost, constants.Notifications, bytes.NewBuffer(jsonMarshal))
	if err != nil {
		return
	}

	// This section was made with the help from a co-student Mathias W.S
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(NotificationHandler)

	handler.ServeHTTP(recorder, request)

	result := recorder.Result()

	_, err = ioutil.ReadAll(result.Body)
	if err != nil {
		return
	}

	AddWebhook(recorder, request)

	after := GetNumberOfWebhooks()

	assert.NotEqual(t, before, after)
}

func TestGetDocument(t *testing.T) {
	var result structs.Notification

	expected := structs.Notification{
		Url:          "localhost:8080/client",
		Country:      "Switzerland",
		Calls:        4,
		CurrentCalls: 0,
		WebHookId:    "nC1fnMeKh9aifAYnhopT",
	}

	GetDocument(constants.NotficationsDb, "nC1fnMeKh9aifAYnhopT", &result)

	assert.Equal(t, expected, result)
}

func TestWriteToDB(t *testing.T) {

	allBefore, err := client.Collection("testWrite").Documents(ctx).GetAll()
	if err != nil {
		return
	}

	lenBefore := len(allBefore)

	assert.Nil(t, err)

	cv19case := structs.MostRecentCases{
		Country:    "Test",
		Date:       "2022-01-01",
		Confirmed:  100,
		Recovered:  200,
		Deaths:     50,
		GrowthRate: 0.45678,
		TimeStamp:  time.Now(),
	}

	WriteToDB("testWrite", cv19case)

	all, err := client.Collection("testWrite").Documents(ctx).GetAll()
	if err != nil {
		return
	}

	assert.Nil(t, err)

	lenAfter := len(all)

	assert.NotEqual(t, lenBefore, lenAfter)
}

func TestIncrement(t *testing.T) {

	noti := structs.Notification{
		Url:          "localhost:8080/client",
		Country:      "Norway",
		Calls:        4,
		CurrentCalls: 0,
	}

	WriteToDB("testIncrement", noti)

	Increment("Norway", "testIncrement")

	iter := client.Collection("testIncrement").Where("Country", "==", "Norway").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		var webhook structs.Notification
		err = doc.DataTo(&webhook)

		assert.NotEqual(t, noti, webhook)
	}
}

func TestViewWebhook(t *testing.T) {

	m := structs.Notification{
		Url:       "localhost:8080/client",
		Country:   "Switzerland",
		Calls:     4,
		WebHookId: "fWP0WHcR26L3GWARHqlX",
	}

	result := ViewOneWebhook("fWP0WHcR26L3GWARHqlX")

	assert.Equal(t, m, result)
}

func TestGetNumberOfWebhooks(t *testing.T) {
	result := GetNumberOfWebhooks()

	getAll, err := client.Collection(constants.NotficationsDb).Documents(ctx).GetAll()

	assert.Nil(t, err)

	expected := len(getAll)

	assert.Equal(t, expected, result)
}

func TestInvalidNotification(t *testing.T) {
	expected := structs.Notification{}

	GetDocument(constants.NotficationsDb, "XDDDD", &expected)

	assert.Equal(t, structs.Notification{}, expected)
}
