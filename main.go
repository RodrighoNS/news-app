package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var tpl = template.Must(template.ParseFiles("index.html"))

var apiKey *string

type Source struct {
	ID   interface{} `json:"id"`
	Name string      `json:"name"`
}

type Article struct {
	Source      Source    `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}

type Results struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

type Search struct {
	SearchKey  string
	Nextpage   int
	TotalPages int
	Results    Results
}

func (a *Article) FormatPublishedDate() string {
	year, month, day := a.PublishedAt.Date()
	return fmt.Sprintf("%v %d, %d", month, day, year)
}

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
	searchKey := params.Get("q")
	page := params.Get("page")
	if page == "" {
		page = "1"
	}

	fmt.Println("Search Query is: ", searchKey)
	fmt.Println("Results page is: ", page)

	// new inst of "Search" struct
	search := &Search{}
	search.SearchKey = searchKey

	// convert "page" into an integer "next"
	next, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "Unexpected server error", http.StatusInternalServerError)
		return
	}

	// "next" being assigned to "Nextpage"
	search.Nextpage = next
	pageSize := 20 // can be 0 to 100

	endpoint := fmt.Sprintf(
		"https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%d&apiKey=%s&sortBy=publishedAt&language=en", url.QueryEscape(search.SearchKey), pageSize, search.Nextpage, *apiKey)

	resp, err := http.Get(endpoint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&search.Results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	search.TotalPages = int(math.Ceil(float64(search.Results.TotalResults / (pageSize))))

	err = tpl.Execute(w, search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func main() {
	apiKey = flag.String("apiKey", "", "Newsapi.org access key")
	flag.Parse()

	if *apiKey == "" {
		log.Fatal("apiKey must be set")
	}

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
