package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
)

var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("<h1>Hello World! - 1st webserver with Go</h1>"))

	// 1st arg: where we want to write the output
	// 2nd arg: data we want to pass to the template
	tpl.Execute(w, nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	params := u.Query()
	searchKey := params.Get("query")
	page := params.Get("page")
	if page == "" {
		page = "1"
	}

	fmt.Println("Search Query is: ", searchKey)
	fmt.Println("Results page is: ", page)
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

	// Instantiate a file server object and passing the directory
	// where all our static files are placed
	fs := http.FileServer(http.Dir("assets"))

	// Telling our router to use this File Server obj
	// for all paths beginning with "/assets/" prefix
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Handle search for the search path "/search"
	mux.HandleFunc("/search", searchHandler)

	// Handle function for the root path "/"
	mux.HandleFunc("/", indexHandler)

	// starts the server using "port" var parameters
	http.ListenAndServe(":"+port, mux)
}
