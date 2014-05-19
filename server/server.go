package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/smoogle/webdav"
)

func main() {
	dir, err := ioutil.TempDir("", "webdav")	
	if err != nil {
		log.Fatalf("could not create temporary directory: %v", err)
	}

	// http.StripPrefix is not working, webdav.Server has no knowledge
	// of stripped component, but needs for COPY/MOVE methods.
	// Destination path is supplied as header and needs to be stripped.
	http.Handle("/webdav/", &webdav.Server{
		Fs:         webdav.Dir(dir),
		TrimPrefix: "/webdav/",
		Listings:   true,
	})

	http.HandleFunc("/", index)

	log.Println("Listening on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q\n", r.URL.Path)
}
