package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"strings"

	"github.com/dimfeld/httptreemux"
	"gopkg.in/mgo.v2"
)

type SearchResult struct {
	URL       string `json:"url"`
	Snippet   string `json:"snippet"`
	Thumbnail string `json:"thumbnail"`
	Context   string `json:"context"`
}

type SearchResultResponse []SearchResult

type SearchQuery struct {
	SearchID string    `json:"-" bson:"search_id"`
	Term     string    `json:"term" bson:"term"`
	When     time.Time `json:"when" bson:"when"`
}

type SearchQueryResponse []SearchQuery

func main() {

	port := os.Getenv("PORT")
	if !strings.Contains(port, ":") {
		port = ":" + port
	}
	api := NewCustomSearchAPI()

	masterDB, err := mgo.Dial(os.Getenv("MONGO_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	h := Handlers{
		api:      api,
		masterDB: masterDB,
	}

	r := httptreemux.New()
	r.GET("/api/imagesearch/:query", h.ImageSearch)
	r.GET("/api/latest/imagesearch", h.Latest)

	log.Printf("Listening on port %s\n", port)
	http.ListenAndServe(port, r)
}
