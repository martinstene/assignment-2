package client

import (
	"log"
	"net/http"
)

// GetResponseFromURL method that takes an url and gets a json response from the webpage.
// Gotten from 05-REST-client in main.go
func GetResponseFromURL(url string) (*http.Response, error) {
	resp, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Print("Error in creating request:", err.Error())
	}

	// Setting the content type header
	resp.Header.Add("content-type", "application/json")

	// Instantiate the client
	client := &http.Client{}

	// Issue request
	res, err := client.Do(resp)
	if err != nil {
		log.Print("Error in response:", err.Error())
	}

	return res, nil
}
