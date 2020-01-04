package main

import (
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello World!</h1>"))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// creates a new Multiplexor server
	// a request multiplexer matches the URL of incoming requests against
	// a list of registered paths and calls the associated handler for the path whenever a match is found.
	mux := http.NewServeMux()

	// Handle function for the root path "/"
	mux.HandleFunc("/", indexHandler)

	// starts the server using "port" var parameters
	http.ListenAndServe(":"+port, mux)
}
