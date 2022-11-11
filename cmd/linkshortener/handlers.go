package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
	"github.com/homayoonalimohammadi/go-link-shortener/linkshortener/internal/app/core"
	"github.com/lib/pq"
)

type CreateRequest struct {
	Url string
}

type CreateResponse struct {
	Token string
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to the Link Shortener!")
}

// Get token stats given as a path parameter.
// These stats include original url and other potential meta data.
func getTokenStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token, ok := vars["token"]
	if !ok {
		http.Error(w, "unable to extract token from the url", http.StatusBadRequest)
		return
	}

	// retrieve url from the given token
	url, err := linkShortener.Postgres.GetUrl(token)
	if err != nil {
		log.Println("unable to retrieve url from postgres:", err)
		io.WriteString(w, "unable to retrieve token data")
		return
	}
	io.WriteString(w, fmt.Sprintf("fetched url from \"%s\": %s", token, url))
}

// Create unique token from the link given in the request body and returns it.
// If the url already has a token, returns the token.
func createToken(w http.ResponseWriter, r *http.Request) {
	log.Println("creating the link given in the post request...")

	// check for correct header
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "invalid Content-Type: must be application/json"
			http.Error(w, msg, http.StatusBadRequest)
		}
	}

	// set max bytes to prevent too large request body
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	decoder := json.NewDecoder(r.Body)

	// disallowing unknown fields
	decoder.DisallowUnknownFields()

	var req CreateRequest
	err := decoder.Decode(&req)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// generate token
	token := core.GenerateToken(10)
	url := req.Url
	if isComplete := strings.HasPrefix(url, "https://"); !isComplete {
		url = "https://" + url
	}

	// save token in DB(s) asynchronously
	err = linkShortener.Postgres.Save(url, token)
	if err == nil {
		log.Printf("saved token for \"%s\": %s", url, token)
		io.WriteString(w, fmt.Sprintf("saved token for \"%s\": %s", url, token))
		return
	}

	// handle secondary errors
	pqError, ok := err.(*pq.Error)
	if !ok || pqError.Code.Name() != "unique_violation" {
		log.Println("error saving record:", err)
		io.WriteString(w, "something went wrong... try again later.")
		return
	}

	// handle unique_violation error and token retrieval
	log.Println("attemped to save duplicate url or token")
	token, err = linkShortener.Postgres.GetToken(url)
	if err != nil {
		http.Error(w, "something went wrong... try again later", http.StatusInternalServerError)
		return
	}
	io.WriteString(w, fmt.Sprintf("a token is already available for this url: %s", token))
}

// Redirects to the original url given a token as a path parameter.
// Resembles the main functionality of the app that by heading to the
// https://<address>/<token> you will get redirected to the original url.
func redirectToOriginal(w http.ResponseWriter, r *http.Request) {
	log.Println("redirecting to the original link...")
	vars := mux.Vars(r)
	token, ok := vars["token"]
	if !ok {
		io.WriteString(w, "unable to extract token from the url")
		return
	}
	url, err := linkShortener.Postgres.GetUrl(token)
	if err != nil {
		io.WriteString(w, "no url is registered for this token")
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}
