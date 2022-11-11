package notifications

import (
	"net/http"
)

// NotificationHandler Takes in methods written in the firebase.go file.
// The purpose of this handler is to view, add and delete webhooks.
func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ViewWebhook(w, r)
	case http.MethodPost:
		AddWebhook(w, r)
	case http.MethodDelete:
		DeleteWebhook(w, r)
	default:
		http.Error(w, "Method not supported.", http.StatusNotImplemented)
	}
}
