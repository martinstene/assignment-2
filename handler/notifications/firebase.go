package notifications

import (
	"assignment-2/constants"
	"assignment-2/json_coder"
	"assignment-2/structs"
	"assignment-2/util"
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"strings"
	"time"
)

// Firebase context and client used by Firestore functions throughout the program.
var ctx context.Context
var client *firestore.Client

// RunFirebase initializes firestore DB
func RunFirebase(authentication string) {
	// Firebase initialisation
	ctx = context.Background()
	sa := option.WithCredentialsFile(authentication)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}
	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

func CloseFirebase() {
	err := client.Close()
	if err != nil {
		return
	}
}

func AddWebhook(w http.ResponseWriter, r *http.Request) {
	util.ContentType(w)

	var notification structs.Notification

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&notification); err != nil {
		// Error handling
		http.Error(w, "Reading of payload failed", http.StatusInternalServerError)
		return
	}

	fmt.Println("Received message ", notification)
	// Checks to see if the url is enpty
	if len(notification.Url) == 0 {
		// Error handling
		http.Error(w, "Your message appears to be empty. Ensure to terminate URI with /.", http.StatusBadRequest)
		return
	} else {
		// Add element in embedded structure.
		// Note: this structure is defined by the client.
		id, _, err := client.Collection(constants.NotficationsDb).Add(ctx, &notification)
		if err != nil {
			// Error handling
			http.Error(w, "The request did not meet the requirements.", http.StatusExpectationFailed)
			return
		}
		// Updates the WebhookID using the ID provided in the DB
		// Sets a timestamp using a built-in function in firestore.
		_, err = client.Collection(constants.NotficationsDb).Doc(id.ID).Update(ctx, []firestore.Update{
			{
				Path:  "WebHookId",
				Value: id.ID,
			},
			{
				Path:  "time",
				Value: firestore.ServerTimestamp,
			},
		})
		if err != nil {
			// Error handling
			http.Error(w, "Something went wrong", http.StatusBadRequest)
			return
		}

		if err != nil {
			// Error handling
			http.Error(w, "Error when adding notifications with URL: "+notification.Url+", Error: "+err.Error(), http.StatusBadRequest)
			return
		} else {
			// Entry added to collection.
			webhookID := structs.CreatedWebHookResponse{
				WebHookId: id.ID,
			}
			// Sets the status code to be 201
			w.WriteHeader(http.StatusCreated)

			encoder := json.NewEncoder(w)
			err = encoder.Encode(&webhookID)

			if err != nil {
				// Error handling
				http.Error(w, "An error occurred when encoding", http.StatusBadRequest)
				return
			}

		}
	}
}

// GetNumberOfWebhooks is used for the status endpoint and gets the total amount of webhooks registered in the DB
func GetNumberOfWebhooks() int {
	getDocuments, err := client.Collection(constants.NotficationsDb).Documents(ctx).GetAll()
	// Error handling
	if err != nil {
		log.Print("No documents found in the specified collection" + err.Error())
		return 0
	}
	numberOfWebhooks := len(getDocuments)
	return numberOfWebhooks
}

/*
	DeleteWebhook gets a webhook ID in the URL, then check if this ID exists in the DB,
	proceeds to delete the webhook if the IDs match.
*/
func DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	// ConvertMapToJson for embedded message ID from URL
	elem := strings.Split(r.URL.Path, "/")
	webHookID := elem[4]

	if len(webHookID) != 0 {
		// Gets the document from firestore.
		resp, err := client.Collection(constants.NotficationsDb).Doc(webHookID).Get(ctx)
		if err != nil {
			// Error handling
			if !resp.Exists() {
				http.Error(w, "No webhook using this ID, it's either deleted or never existed.", http.StatusGone)
				return
			}
			// If the ID is in the DB, delete the specific document.
			_, err := client.Collection(constants.NotficationsDb).Doc(webHookID).Delete(ctx)
			if err != nil {
				return
			}

		}
		// Feedback to the user.
		http.Error(w, "The webhook was deleted successfully.", http.StatusOK)
	}
}

func ViewOneWebhook(webHookID string) interface{} {
	// Retrieve specific message based on id (Firestore-generated hash)
	response := client.Collection(constants.NotficationsDb).Doc(webHookID)

	// Retrieve reference to document
	doc, err2 := response.Get(ctx)
	if err2 != nil {
		log.Printf("%v", "Error extracting body of returned document of message "+webHookID)
	}
	// A message map with string keys. Each key is one field, like "text" or "timestamp"
	m := doc.Data()
	if m == nil {
		log.Printf("%v", "The list seems to be empty")
	}
	return util.ConvertMapToJson(m)
}

/*
	ViewWebhook is created using the displayMessage method in example 11-Firestore-demo.
	It takes in an ID from the URL and checks if it's a valid ID and gets this document, which is then
	encoded from a map to json and printed out as a response. If the user doesn't input an ID it will return the entire
	list of webhooks stored in the DB.
*/
func ViewWebhook(w http.ResponseWriter, r *http.Request) {
	// Gets the webhookID from the last part of the URL
	elem := strings.Split(r.URL.Path, "/")
	webHookID := elem[4]

	if len(webHookID) != 0 {
		// Extract individual message
		m := ViewOneWebhook(webHookID)
		json_coder.PrettyPrint(w, m)
	}
	if len(webHookID) == 0 {
		iter := client.Collection(constants.NotficationsDb).Documents(ctx)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}

			// A message map with string keys. Each key is one field, like "text" or "timestamp"
			m := doc.Data()

			json_coder.PrettyPrint(w, m)
		}
	}
}

/*
	WriteToDB takes in a collection-name and an interface to add a document to a collection
	and adds a timestamp regardless of what collection or what interface.
*/
func WriteToDB(collectionName string, data interface{}) {
	id, _, err := client.Collection(collectionName).Add(ctx, &data)
	// Error handling
	if err != nil {
		log.Fatalln("Something went wrong when adding")
		return
	}
	// Adds timestamp to the document.
	_, _ = client.Collection(collectionName).Doc(id.ID).Update(ctx, []firestore.Update{
		{
			Path:  "time",
			Value: firestore.ServerTimestamp,
		},
	})
}

// GetDocument gets data from a collection and a document in the database.
func GetDocument(collection string, document string, structToExtractTo interface{}) {
	// Gets a document from the collection on the DB
	res := client.Collection(collection).Doc(document)
	// Gets the document reference and sets the context
	doc, err := res.Get(ctx)

	if !doc.Exists() {
		log.Print("The doc does not exist, returning.")
		return
	}

	if err != nil {
		log.Printf("Error getting document: %v", err)
		return
	}
	// Sends the data received in the document to a struct
	err = doc.DataTo(&structToExtractTo)

	if err != nil {
		log.Printf("Error extracting data into struct: %v", err)
		return
	}
}

// CheckIfExistsPolicy goes through all documents with a countryCode and a scope like the ones inputted
// and checks if this exists, if it does, prints it from the database.
func CheckIfExistsPolicy(w http.ResponseWriter, collection string, country string, date string) bool {
	iter := client.Collection(collection).
		Where("CountryCode", "==", country).
		Where("Scope", "==", date).
		Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
			return false
		}
		if doc.Exists() {
			documentSnapshot, err := doc.Ref.Get(ctx)
			if err != nil {
				return false
			}
			json_coder.PrettyPrint(w, documentSnapshot.Data())
			return false
		}
	}
	return true
}

// CheckIfExistsCases goes through all documents with a country name like the one inputted
// and checks if this exists, if it does, prints it from the database.
func CheckIfExistsCases(w http.ResponseWriter, collection string, country string) bool {
	// Create a struct variable to use it for checking if the case is in the DB.
	var dataFromDatabase structs.MostRecentCases

	iter := client.Collection(collection).
		Where("Country", "==", country).
		Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		// Gets the document from the DB to get the timestamp used for cache purging.
		GetDocument(constants.CasesDb, doc.Ref.ID, &dataFromDatabase)
		// If it's not more than 8 hours since the document was updated, it will get the document
		// from the DB and print this out to the user.
		if !(time.Since(dataFromDatabase.TimeStamp).Hours() > (time.Hour * 8).Hours()) {
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
				return false
			}
			if doc.Exists() {
				documentSnapshot, err := doc.Ref.Get(ctx)
				if err != nil {
					return false
				}
				json_coder.PrettyPrint(w, documentSnapshot.Data())
				return false
			}
		}
		// If the document IS older than 8 hours, this will be deleted from the DB and it will
		// return true. This will send it into the normal logic method which will do an API call
		// and get the results.
		if time.Since(dataFromDatabase.TimeStamp).Hours() > (time.Hour * 8).Hours() {
			documentSnapshot, err := doc.Ref.Get(ctx)
			if err != nil {
				return false
			}
			_, err = documentSnapshot.Ref.Delete(ctx)
			if err != nil {
				http.Error(w, "Error when deleting document.", http.StatusBadRequest)
				return false
			}
			return true
		}
	}
	return true
}

/*
	Increment checks if a country matches, then goes through each document containing this country,
	then updating this document by incrementing the CurrentCalls field. Then calls the CheckIfCallsMatch method.
*/
func Increment(country string, collection string) {
	iter := client.Collection(collection).Where("Country", "==", country).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
			return
		}
		_, err = doc.Ref.Update(ctx, []firestore.Update{
			{
				Path:  "CurrentCalls",
				Value: firestore.Increment(1),
			},
		})
	}
	CheckIfCallsMatch(country)
}

/*
	CheckIfCallsMatch gets the notifications collection, checks if the country matches, the iterates through all
	entries which are matching. And if the webhook calls matches the current calls, then trigger the webhook.
*/
func CheckIfCallsMatch(country string) {
	iter := client.Collection(constants.NotficationsDb).
		Where("Country", "==", country).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
			return
		}

		var webhook structs.Notification
		err = doc.DataTo(&webhook)
		if webhook.Calls == webhook.CurrentCalls {
			// Uses "go" to trigger a new thread
			go trigger(webhook)
		}
	}
}

/*
	trigger is used for the webhook invocation. It takes in a struct
	of type Notification and converts this to the specified struct for invocation.
	Marshals this data and sends a post request to the url specified in the webhook.
*/
func trigger(webhook structs.Notification) {
	webHookInvocation := structs.WebHookInvocation{
		WebHookId: webhook.WebHookId,
		Country:   webhook.Country,
		Calls:     webhook.Calls,
	}
	jsonData, _ := json.Marshal(webHookInvocation)
	_, err := http.Post(webhook.Url, "application/json", bytes.NewBuffer(jsonData))
	// Error handling
	if err != nil {
		log.Print("Error during request creation. Error:", err)
		return
	}
}
