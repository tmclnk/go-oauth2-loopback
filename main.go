package main

import (
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var ClientID = os.Getenv("CLIENT_ID")
var AuthURL = os.Getenv("AUTH_URL")
var TokenURL = os.Getenv("TOKEN_URL")
var RedirectURL = os.Getenv("REDIRECT_URL")
var ListenAddress = os.Getenv("LISTEN_ADDRESS")

func main() {
	var conf = &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: "",
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  AuthURL,
			TokenURL: TokenURL,
		},
		RedirectURL: RedirectURL,
	}

	client := make(chan *http.Client)
	quit := make(chan struct{}) // signal-only channel
	go func() {
		exampleResourceCall(<-client)
		close(client)
		quit <- struct{}{}
		close(quit)
	}()
	SetupClient(conf, client)

	// Don't shut down until notified
	<-quit
}

// Make a call to a secured resource
func exampleResourceCall(client *http.Client) {
	log.Println("Created client...")
	const resourceUrl = "https://spring.users.runpaste.com/users/123"
	resp, err := client.Get(resourceUrl)
	log.Println("Got resource response...")
	if err != nil {
		log.Fatalf("Resource retrieval: %s\n", err)
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	log.Printf(string(bytes))
}
