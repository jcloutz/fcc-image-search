package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"io"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

// Error Response defines the structure of an api error
type ErrorResponse struct {
	Error string `json:"error"`
}

// Handlers is responsible for processing requests to the api
type Handlers struct {
	masterDB *mgo.Session
	api      *CustomSearchAPI
}

var (
	ErrInvalidOffsetArgument          = errors.New("Invalid offset argument, must be an integer >= 1")
	ErrUnableToQueryAPI               = errors.New("Unable to query search API")
	ErrUnableToSaveQuery              = errors.New("Unable to save document to database")
	ErrUnableToLocatePreviousSearches = errors.New("Unable to locate previous searches")
)

// ImageSearch is the handler for querying google's api
func (h *Handlers) ImageSearch(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset := 1
	var err error

	if r.URL.Query().Get("offset") != "" {
		offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil || offset < 1 {
			h.RespondError(w, ErrInvalidOffsetArgument, http.StatusBadRequest)

			return
		}
	}

	resp, err := h.api.New().Offset(offset).Search(params["query"]).Get()
	if err != nil {
		h.RespondError(w, ErrUnableToQueryAPI, http.StatusInternalServerError)

		return
	}

	response := make(SearchResultResponse, len(resp.Items))
	for k, v := range resp.Items {
		response[k] = SearchResult{
			URL:       v.Link,
			Snippet:   v.Snippet,
			Context:   v.Image.ContextLink,
			Thumbnail: v.Image.ThumbnailLink,
		}
	}

	sess := h.masterDB.Copy()
	defer sess.Close()

	if err := ArchiveSearch(sess, params["query"]); err != nil {
		h.RespondError(w, ErrUnableToSaveQuery, 500)
	}

	h.Respond(w, response, 200)
	return
}

// Latest fetches the most recent searches from the database
func (h *Handlers) Latest(w http.ResponseWriter, r *http.Request, params map[string]string) {
	db := h.masterDB.Copy()
	defer db.Close()

	results, err := RetrieveSearches(db)
	if err != nil {
		h.RespondError(w, ErrUnableToLocatePreviousSearches, 404)

		return
	}

	h.Respond(w, results, http.StatusOK)
	return
}

// RespondError handle all error responses to the client
func (h *Handlers) RespondError(w http.ResponseWriter, apiError error, statusCode int) {

	e := ErrorResponse{
		Error: apiError.Error(),
	}

	h.Respond(w, e, statusCode)
}

// Respond handles all responses to the client
func (h *Handlers) Respond(w http.ResponseWriter, value interface{}, statusCode int) {

	js, err := json.Marshal(value)
	if err != nil {
		fmt.Println("error", err)
		js = []byte("{}")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	io.WriteString(w, string(js))
}
