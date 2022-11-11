package linkshortener

import (
	"io"
	"log"
	"net/http"
)

func getLink(w http.ResponseWriter, r *http.Request) {
	log.Println("getting link from the url...")
	io.WriteString(w, "fetched url from out DB: example.com")
}

func createLink(w http.ResponseWriter, r *http.Request) {
	log.Println("creating the link given in the post request...")
	io.WriteString(w, "created toekn: test_token")
}
