package main

import (
	"assignment-2/constants"
	"assignment-2/handler/cases"
	"assignment-2/handler/notifications"
	"assignment-2/handler/root"
	"assignment-2/handler/status"
	"assignment-2/handler/stringency"
	"log"
	"net/http"
	"os"
)

// handleRequests sets the port to be used.
func handleRequests() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Println("Port is not set, setting it to: " + "8080")
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// This is a main() function
func main() {
	// Initializes firebase
	notifications.RunFirebase("CredentialsFirestore.json")
	defer notifications.CloseFirebase()

	// Initializes endpoints
	http.HandleFunc(constants.Root, root.HomePageHandler)
	http.HandleFunc(constants.Cases, cases.CaseHandler)
	http.HandleFunc(constants.Notifications, notifications.NotificationHandler)
	http.HandleFunc(constants.Policy, stringency.PolicyHandler)
	http.HandleFunc(constants.Status, status.StatusHandler)

	// Sets the port and listens
	handleRequests()
}
