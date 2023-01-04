package main

import (
	"context"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"net/url"
	"time"
)

// SetupClient creates an oauth2-enabled http client by launching a browser
// on the user's machine and listening on the loopback address
// for a response. You'll need to register the loopback address
// and port as a redirect_url in your Authentication Server.
// The client sent to the channel will be configured with OAuth2 tokens.
// Once the token is received, the http server will be shut down.
func SetupClient(conf *oauth2.Config, ch chan *http.Client) {

	// set up http server
	ctx := context.Background()
	router := http.NewServeMux()
	srv := &http.Server{
		Addr:    ListenAddress,
		Handler: router,
	}
	srv.RegisterOnShutdown(func() {
		log.Println("Shutting server down...")
	})

	// make a channel we can use to shut down the listening server
	// once we're done with it
	tok := make(chan *oauth2.Token)
	go func() {
		// we got the token we need to make a client,
		// we can shut the server down here
		ch <- conf.Client(ctx, <-tok)
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("%s\n", err)
		}
	}()

	// our handler function will get the token, so we have to pass it the channel
	router.HandleFunc("/oauth/callback", createCallbackHandler(ctx, conf, tok))

	// launch the server in a goroutine so this function doesn't block
	go func() {
		log.Printf("Listening on %s...\n", ListenAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("%s\n", err)
		}
	}()

	// open a web browser pointing to the auth server
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	log.Printf("Sleeping...")
	time.Sleep(10 * time.Second)
	log.Printf("Opening browser at %s\n", url)
	_ = open.Run(url)
}

// Creates handler that receives the authorization code from the auth server
// and exchanges it for an access token. The token is sent on the channel
// and a success message is returned to the user.
func createCallbackHandler(ctx context.Context, conf *oauth2.Config, ch chan *oauth2.Token) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Exchange authorization code for access token
		queryParts, _ := url.ParseQuery(r.URL.RawQuery)

		if queryParts["code"] == nil && len(queryParts["code"]) > 0 {
			msg := fmt.Sprintf("Token exchange failed! Missing 'code' parameter.")
			_, _ = fmt.Fprintf(w, msg)
			w.WriteHeader(http.StatusBadRequest)
			log.Println(msg)
			return
		}

		code := queryParts["code"][0]
		err, tok := exchange(code, conf, ctx)

		if err != nil {
			msg := fmt.Sprintf("Token exchange failed! %s", err.Error())
			_, _ = fmt.Fprintf(w, msg)
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
		ch <- tok

		// Write success message back to browser
		msg := `
		<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN">
		<html>
		<head>
		</head>
		Success. You may close this window.
		</body>
		</html>
		`
		_, _ = fmt.Fprintf(w, msg)
	}
}

// Exchanges an authorization code for an access token
func exchange(code string, conf *oauth2.Config, ctx context.Context) (error, *oauth2.Token) {
	log.Printf("Code: %s\n", code)
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		return err, nil
	}
	log.Printf("Token Type: %s\n", tok.TokenType)
	log.Printf("Access Token: %s\n", tok.AccessToken)
	return nil, tok
}
