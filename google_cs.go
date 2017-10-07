package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// https://www.googleapis.com/customsearch/v1?q=lol+cats&ct=013327425812576012613:3t1bx0dps9c&searchType=image&num=2&start=2
const googleCSBaseURL = `https://www.googleapis.com/customsearch/v1`

// CustomSearchAPI provides an abstraction to the google custom search service
type CustomSearchAPI struct {
	values url.Values
}

// CustomSearchQuery provides an abstraction to the google custom api and is responsible for creating a new search request
type CustomSearchQuery struct {
	values url.Values
}

// NewCustomSearchAPI creates a new custom search api instance
func NewCustomSearchAPI() *CustomSearchAPI {
	v := url.Values{}
	v.Add("key", "AIzaSyDKA09pPyjSBYfPVb7Pc5Zq-MKSXXwKcg0")
	v.Add("cx", "013327425812576012613:3t1bx0dps9c")
	v.Add("searchType", "image")
	v.Add("num", "10")
	v.Add("start", "1")
	v.Add("safe", "high")

	return &CustomSearchAPI{
		values: v,
	}
}

// New creates a new query
func (csa *CustomSearchAPI) New() *CustomSearchQuery {
	val := url.Values{}

	for k, v := range csa.values {
		val[k] = v
	}

	q := CustomSearchQuery{values: val}

	return &q
}

// Search sets the query value
func (q *CustomSearchQuery) Search(query string) *CustomSearchQuery {
	q.values.Add("q", url.QueryEscape(query))

	return q
}

// Offset sets the starting index
func (q *CustomSearchQuery) Offset(offset int) *CustomSearchQuery {
	l := strconv.Itoa(offset)
	q.values.Set("start", l)

	return q
}

// Get sends the the request and returns the decoded response
func (q *CustomSearchQuery) Get() (*GoogleAPIResponse, error) {
	query := googleCSBaseURL + "?" + q.values.Encode()

	resp, err := http.Get(query)
	if err != nil {
		return nil, errors.Wrap(err, "Error querying api")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read response body")
	}

	apiResponse := GoogleAPIResponse{}
	err = json.Unmarshal(body, &apiResponse)

	return &apiResponse, err
}
