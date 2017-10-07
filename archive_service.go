package main

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ArchiveSearch stores the search query in the database
func ArchiveSearch(db *mgo.Session, query string) error {
	search := SearchQuery{
		SearchID: bson.NewObjectId().Hex(),
		Term:     query,
		When:     time.Now(),
	}

	return db.DB("").C("searches").Insert(search)
}

// RetrieveSearches fetches the last 20 searches from the database
func RetrieveSearches(db *mgo.Session) (*SearchQueryResponse, error) {
	var results SearchQueryResponse

	if err := db.DB("").C("searches").Find(nil).Limit(20).Sort("-when").All(&results); err != nil {
		return nil, err
	}

	return &results, nil
}
